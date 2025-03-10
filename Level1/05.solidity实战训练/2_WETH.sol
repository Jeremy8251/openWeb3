// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
WETH 是包装 ETH 主币，作为 ERC20 的合约。 标准的 ERC20 合约包括如下几个

3 个查询
    balanceOf: 查询指定地址的 Token 数量
    allowance: 查询指定地址对另外一个地址的剩余授权额度
    totalSupply: 查询当前合约的 Token 总量
2 个交易
    transfer: 从当前调用者地址发送指定数量的 Token 到指定地址。
        这是一个写入方法，所以还会抛出一个 Transfer 事件。
    transferFrom: 当向另外一个合约地址存款时，对方合约必须调用 transferFrom 才可以把 Token 拿到它自己的合约中。
2 个事件
    Transfer
    Approval
1 个授权
    approve: 授权指定地址可以操作调用者的最大 Token 数量。
 */
contract WETH {
    string public name = "Wrapped Ether";
    string public symbol = "WETH";
    uint8 public decimals = 18;
    event Approval(
        address indexed src,
        address indexed delegateAds,
        uint256 amount
    );
    event Transfer(address indexed src, address indexed toAds, uint256 amount);
    event Deposit(address indexed toAds, uint256 amount);
    event WithDraw(address indexed src, uint256 amount);
    mapping(address => uint256) public balanceOf;
    //授权指定地址可以操作调用者的最大 Token 数量
    mapping(address => mapping(address => uint256)) public allowance;

    // 存钱
    function deposit() public payable {
        // 存进钱
        balanceOf[msg.sender] += msg.value;

        emit Deposit(msg.sender, msg.value);
    }

    // 取钱
    function withDraw(uint256 amount) public {
        require(balanceOf[msg.sender] >= amount, "balance is not enough");
        balanceOf[msg.sender] -= amount;
        // 转入msg.sender款amount
        payable(msg.sender).transfer(amount);
        emit WithDraw(msg.sender, amount);
    }

    // 查看合约余额
    function totalSupply() public view returns (uint256) {
        return address(this).balance;
    }

    // 授权地址
    //授权指定地址（delegateAds）从调用者账户（msg.sender）中转移一定数量（amount）的代币
    function approval(
        address delegateAds,
        uint256 amount
    ) public returns (bool) {
        allowance[msg.sender][delegateAds] += amount;
        emit Approval(msg.sender, delegateAds, amount);
        return true;
    }

    // 转账
    function transfer(address toAds, uint256 amount) public returns (bool) {
        return transferFrom(msg.sender, toAds, amount);
    }

    function transferFrom(
        address src,
        address toAds,
        uint256 amount
    ) public returns (bool) {
        // 检查转出地址的余额是否足够
        require(balanceOf[src] >= amount, "src balance is not enough");
        // 如果调用者不是转出地址本人，需验证授权额度
        if (src != msg.sender) {
            // 检查朋友有没有超过你给的额度（100 块）
            require(allowance[src][msg.sender] >= amount, "allownce");
            // 朋友每刷一笔，额度就减少（比如刷了 30，剩余额度变 70）
            allowance[src][msg.sender] -= amount; // 减少授权额度
        }
        // 更新转出和接收地址的余额
        balanceOf[src] -= amount;
        balanceOf[toAds] += amount;
        // 触发 Transfer 事件（ERC20 标准要求）
        emit Transfer(src, toAds, amount);
        return true;
    }

    fallback() external payable {
        deposit();
    }

    receive() external payable {
        deposit();
    }
}

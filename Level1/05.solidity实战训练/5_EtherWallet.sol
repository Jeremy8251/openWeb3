// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
这一个实战主要是加深大家对 3 个取钱方法的使用。

任何人都可以发送金额到合约
只有 owner 可以取款
3 种取钱方式
 */
contract EtherWallet {
    address payable public immutable owner;
    event Log(string funName, address from, uint256 value, bytes data);

    constructor() {
        owner = payable(msg.sender);
    }

    fallback() external payable {
        emit Log("fallback", msg.sender, msg.value, "");
    }

    receive() external payable {
        emit Log("receive", msg.sender, msg.value, "");
    }

    function withDraw1() external {
        require(msg.sender == owner, "msg.sender is not owner");
        // owner.transfer 相比 msg.sender 更消耗Gas
        // owner.transfer(address(this).balance);
        payable(msg.sender).transfer(10 * 10 ** 18);
    }

    function withDraw2() external {
        require(msg.sender == owner, "msg.sender is not owner");
        bool success = payable(msg.sender).send(10 * 10 ** 18);
        require(success, "Send failed");
    }

    function withDraw3() external {
        require(msg.sender == owner, "msg.sender is not owner");
        (bool success, ) = msg.sender.call{value: address(this).balance}("");
        require(success, "call failed");
    }

    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}

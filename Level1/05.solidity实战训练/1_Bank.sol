// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
存钱罐合约
所有人都可以存钱
    ETH
只有合约 owner 才可以取钱
只要取钱，合约就销毁掉 selfdestruct
扩展：支持主币以外的资产
    ERC20
    ERC721
*/
contract Bank{

    address public immutable owner;
    // 存款
    event Deposit(address _ads, uint256 amount);
    // 提款
    event WithDraw(uint256 amount);

    receive() external payable{
        emit Deposit(msg.sender, msg.value);
    }
    constructor() payable{
        owner = msg.sender;
    }

    function withDraw() external {
        require(msg.sender == owner, "Not Owner");
        emit WithDraw(address(this).balance);
        // selfdestruct(payable(msg.sender));

    }

    function getBalance() external view returns (uint256){
        return address(this).balance;
    }

}

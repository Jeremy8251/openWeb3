// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Immutable {
    //它的值在合约部署时可以由构造函数设置，但一旦设置后就无法更改
    address public immutable MY_ADDRESS;
    uint256 public immutable MY_UINT;

    constructor(uint256 _myUint) {
        MY_ADDRESS = msg.sender;
        MY_UINT = _myUint;
    }
}

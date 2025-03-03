// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EtherUnits {
    uint256 public oneWei = 1 wei; // 1
    // 1 wei is equal to 1
    bool public isOneWei = (oneWei == 1); //true

    uint256 public oneGwei = 1 gwei; // 1000000000
    // 1 gwei is equal to 10^9 gwei
    bool public isOnewGwei = (oneGwei == 1e9); //true

    uint256 public oneEther = 1 ether; //1000000000 000000000
    // 1 ether is equal to 10^18 wei
    bool public isOneEther = (oneEther == 1e18); //true
}

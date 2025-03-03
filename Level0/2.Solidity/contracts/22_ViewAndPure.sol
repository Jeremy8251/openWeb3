// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
// View function 声明不会更改任何状态。
// Pure function 声明不会更改或读取任何 state 变量。
contract ViewAndPure {
    uint256 public x = 1;

    function addToX(uint256 y) public view returns (uint256) {
        return x + y;
    }

    function add(uint256 i, uint256 j) public pure returns (uint256) {
        return i + j;
    }
}

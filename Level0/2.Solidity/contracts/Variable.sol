// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract variable {
    uint public storageData; // 状态变量
    uint256 blockTime;
    constructor() {
        storageData = 10;
    }

    function getBlockTime() public view returns (uint256) {
        return block.timestamp;
    }

    function getStorageData() public view returns (uint256) {
        return storageData;
    }
}

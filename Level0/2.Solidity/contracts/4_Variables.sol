// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * Solidity 中有 3 种类型的变量
当地
  在函数中声明
  不存储在区块链上
state
  在函数外部声明
  存储在区块链上
global （提供有关区块链的信息）
 */
contract Variables {
    string public text = "Hello";
    uint256 public num = 123;

    function dosomethings() public view returns (uint256, uint256, address) {
        uint256 i = 456;

        uint256 timestamp = block.timestamp;
        address sender = msg.sender;

        return (i, timestamp, sender); //**0:uint256: 456 1:uint256: 1740819114 2:address: 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4 */
    }
}

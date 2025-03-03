// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

//这是一个简单的 Contract，你可以获取、递增和递减此 Contract 中的 count 存储。

contract Counter {
    uint public count;

    function get() public view returns (uint256) {
        return count;
    }

    function inc() public {
        count += 1;
    }

    function dec() public {
        count -= 1;
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Array {
    uint256[] public arr;
    uint256[] public arr2 = [1, 2, 3];
    uint256[10] public myFixedSizeArry;

    function get(uint256 i) public view returns (uint256) {
        return arr[i];
    }

    function getArr() public view returns (uint256[] memory) {
        return arr;
    }

    function push(uint256 i) public {
        arr.push(i);
    }

    function pop() public {
        arr.pop();
    }

    function getLength() public view returns (uint256) {
        return arr.length;
    }

    function remove(uint256 index) public {
        delete arr[index];
    }

    function examples() external pure returns (uint256[] memory) {
        //内存（memory）中创建一个固定大小的数组
        uint256[] memory a = new uint256[](5);
        return a;
    }
}

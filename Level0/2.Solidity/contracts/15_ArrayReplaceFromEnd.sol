// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ArrayReplaceFromEnd {
    uint256[] public arr;

    // Deleting an element creates a gap in the array.
    // One trick to keep the array compact is to
    // move the last element into the place to delete.
    function remove(uint256 index) public {
        // Move the last element into the place to delete
        arr[index] = arr[arr.length - 1];
        // Remove the last element
        // 改造了这里
        // arr.pop();
    }

    function getArr() public view returns (uint256[] memory) {
        return arr;
    }

    function test() public {
        arr = [1, 2, 3, 4];

        remove(1);
        // [1, 4, 3, 4]
        assert(arr.length == 4);
        assert(arr[0] == 1);
        assert(arr[1] == 4);
        assert(arr[2] == 3);
        assert(arr[3] == 4);

        remove(2);
        // [1, 4, 4，4]
        assert(arr.length == 4);
        assert(arr[0] == 1);
        assert(arr[1] == 4);
    }
}

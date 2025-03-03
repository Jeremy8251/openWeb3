// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
变量声明为 ， 或 以显式指定数据的位置。storagememorycalldata
storagevariable 是一个状态变量（存储在区块链上）
memory变量在内存中，并且在调用函数时存在
calldata包含函数参数的特殊数据位置
 */
contract DataLocations {
    uint256[] public arr;
    mapping(uint256 => address) map;
    struct MyStruct {
        uint256 foo;
    }

    mapping(uint256 => MyStruct) myStructs;

    function _f(
        uint256[] storage _arr,
        mapping(uint256 => address) storage _map,
        MyStruct storage _myStruct
    ) internal {
        // do something with storage variables
    }

    function g(uint256[] memory _arr) public returns (uint256[] memory) {
        // do something with memory array
    }

    function h(uint256[] calldata _arr) internal {
        // do something with calldata array
    }

    function f() public returns (MyStruct memory mys, MyStruct memory myt) {
        arr = [1, 2, 3, 4];

        myStructs[0] = MyStruct(123);

        MyStruct memory structItem;
        structItem.foo = 456;
        myStructs[1] = structItem;
        // call _f with state variables
        _f(arr, map, myStructs[1]);

        // get a struct from a mapping
        MyStruct storage myStruct = myStructs[1];
        // create a struct in memory
        MyStruct memory myMemStruct = MyStruct(0);
        return (myStruct, myMemStruct);
    }
}

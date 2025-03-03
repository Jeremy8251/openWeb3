// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Primitives {
    bool public boo = true;
    /*
    uint stands for unsigned integer, meaning non negative integers
    different sizes are available
        uint8   ranges from 0 to 2 ** 8 - 1
        uint16  ranges from 0 to 2 ** 16 - 1
        ...
        uint256 ranges from 0 to 2 ** 256 - 1
    */
    uint8 public u8 = 1;
    uint256 public u256 = 2 ** 256 - 1;
    uint256 public u = 123;

    int8 public i8 = 1;
    int256 public i256 = 456;

    int256 public i = -123;
    int256 public min = -2 ** 255 + 1;

    int256 public minInt = type(int256).min;//-57896044618658097711785492504343953926634992332820282019728792003956564819968
    int256 public maxInt = type(int256).max;//57896044618658097711785492504343953926634992332820282019728792003956564819967

    address public addr = 0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c;

    bytes1 a = 0xf5;
    bytes1 b = 0xff;

    bool public defautBoo; // false
    uint256 public defaultUint; //0
    int256 public defaultInt; //0
    address public defaultAddr; //0x0000000000000000000000000000000000000000
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// solcjs --bin Store.sol --abi Store.sol
// abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=Store.go
contract Store {
    event ItemSet(bytes32 key, bytes32 value);

    string public version;
    mapping(bytes32 => bytes32) public items;

    constructor(string memory _version) {
        version = _version;
    }

    function setItem(bytes32 key, bytes32 value) external {
        items[key] = value;
        emit ItemSet(key, value);
    }
}

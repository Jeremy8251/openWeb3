// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract ConstructorX {
    string public name;

    constructor(string memory _name) {
        name = _name;
    }
}

contract ConstructorY {
    string public text;

    constructor(string memory _text) {
        text = _text;
    }
}

contract B is ConstructorX("Input to X"), ConstructorY("Input to Y") {}

contract C is ConstructorX, ConstructorY {
    constructor(
        string memory _name,
        string memory _text
    ) ConstructorX(_name) ConstructorY(_text) {}
}

contract D is ConstructorX, ConstructorY {
    constructor() ConstructorX("X was called") ConstructorY("Y was called") {}
}

contract E is ConstructorX, ConstructorY {
    constructor() ConstructorY("Y was called") ConstructorX("X was called") {}
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Selector {
    bytes4 storedSelector;
    function storeSelector(bytes4 selector) public {
        storedSelector = selector;
    }
    function executeStoredFunction(uint x) public returns (uint) {
        return execute(storedSelector, x);
    }

    function square(uint x) public pure returns (uint) {
        return x * x;
    }

    function double(uint x) public pure returns (uint) {
        return x * 2;
    }

    function getSelector(string memory func) public pure returns (bytes4) {
        return bytes4(keccak256(abi.encodePacked(func)));
    }

    function getSquareSelector() public pure returns (bytes4) {
        return this.square.selector;
    }

    function getdoubleSelector() public pure returns (bytes4) {
        return bytes4(keccak256("double(uint256)"));
    }
    function execute(bytes4 _selector, uint _x) public returns (uint) {
        (bool success, bytes memory data) = address(this).call(
            abi.encodeWithSelector(_selector, _x)
        );
        require(success, "Execution failed");
        return abi.decode(data, (uint));
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Function {
    function returnMany() public pure returns (uint256, bool, uint256) {
        return (1, false, 2);
    }

    function named() public pure returns (uint256 x, bool b, uint256 y) {
        return (3, true, 4);
    }

    function assigned() public pure returns (uint256 x, bool b, uint256 y) {
        x = 5;
        b = false;
        y = 6;
    }

    // function that returns multiple values.
    function destructuringAssignments()
        public
        pure
        returns (uint256, bool, uint256, uint256, uint256)
    {
        (uint256 i, bool b, uint256 j) = returnMany();
        (uint256 x, , uint256 y) = (4, 5, 6);
        return (i, b, j, x, y);
    }

    function arrayInput(uint256[] memory _arr) public {}

    uint256[] public arr;

    function arrayout() public view returns (uint256[] memory) {
        return arr;
    }
}

contract XYZ {
    function someFuncWithManyInputs(
        uint256 x,
        uint256 y,
        uint256 z,
        address a,
        bool b,
        string memory c
    ) public pure returns (uint256) {}

    function callFunc() external pure returns (uint256) {
        return someFuncWithManyInputs(1, 2, 3, address(0), true, "c");
    }

    function callFuncWithKeyValue() external pure returns (uint256) {
        return
            someFuncWithManyInputs({
                a: address(0),
                b: true,
                c: "c",
                x: 1,
                y: 2,
                z: 3
            });
    }
}

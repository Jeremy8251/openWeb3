// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Loop {
    uint256 public j = 0;

    function loop() public {
        for (uint i = 0; i < 10; i++) {
            if (i == 5) continue;
            if (i == 8) break;
        }

        while (j < 10) {
            j++; // 10
        }
    }
}

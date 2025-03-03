// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
错误将撤消在事务期间对状态所做的所有更改。

您可以通过调用 、 或 来引发错误。requirerevertassert

require用于在执行之前验证输入和条件。
revert类似于 。有关详细信息，请参阅下面的代码。require
assert用于检查不应为 false 的代码。断言失败可能意味着存在 bug。
使用自定义错误来节省 gas。
 */

contract Error {
    function testRequire(uint256 _i) public pure {
        require(_i > 10, "Input must be greater than 10");
    }

    function testRevert(uint256 _i) public pure {
        // Revert is useful when the condition to check is complex.
        // This code does the exact same thing as the example above
        if (_i <= 10) {
            revert("Input must be greater than 10");
        }
    }

    uint256 public num;

    function testAssert() public view {
        assert(num == 0);
    }

    error InsufficientBalance(uint256 balance, uint256 withdrawAmount);

    function testCustomError(uint256 _withdrawAmount) public view {
        uint256 bal = address(this).balance;
        if (bal < _withdrawAmount) {
            revert InsufficientBalance({
                balance: bal,
                withdrawAmount: _withdrawAmount
            });
        }
    }
}

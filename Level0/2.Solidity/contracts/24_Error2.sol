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

contract Account {
    uint256 public balance;
    uint256 public constant MAX_UINT = 2 * 256 - 1;

    function deposit(uint256 _amount) public {
        uint256 oldBalance = balance;
        uint256 newBalance = balance + _amount;

        require(newBalance >= oldBalance, "OverFlow");

        balance = newBalance;
        assert(balance >= oldBalance);
    }

    function withdraw(uint256 _amount) public {
        uint256 oldBalance = balance;

        require(balance >= _amount, "require underflow");
        if (balance < _amount) {
            revert("revert Underflow");
        }
        balance -= _amount;

        assert(balance <= oldBalance);
    }
}

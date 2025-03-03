// SPDX-License-Identifier: MIT
// 修改编译目标为 cancun
pragma solidity ^0.8.24;

//存储在 transient storage 中的数据在事务处理后被清除。

interface ITest {
    function val() external view returns (uint256);

    function test() external;
}

contract Callback {
    uint256 public val;

    fallback() external {
        val = ITest(msg.sender).val();
    }

    function test(address target) external {
        ITest(target).test();
    }
}

contract TestStorage {
    uint256 public val;

    function test() public {
        val = 123;
        bytes memory b = "";
        (bool success, ) = msg.sender.call(b); // 捕获调用是否成功
        require(success, "Low-level call failed"); // 可选：失败时回滚交易‌
    }
}

contract TestTransientStorage {
    bytes32 constant SLOT = 0;

    function test() public {
        assembly {
            tstore(SLOT, 321)
        }
        bytes memory b = "";
        (bool success, ) = msg.sender.call(b);
        require(success, "Low-level call failed");

        assembly {
            tstore(SLOT, 0) // 添加清除操作
        }
    }

    function val() public view returns (uint256 v) {
        assembly {
            // 读取值
            v := tload(SLOT)
        }
        return v;
    }
}

contract ReentrancyGuard {
    bool private locked;

    modifier lock() {
        require(!locked);
        locked = true;
        _;
        locked = false;
    }

    // 35313 gas
    function test() public lock {
        // Ignore call error
        bytes memory b = "";
        (bool success, ) = msg.sender.call(b); // 捕获调用是否成功
        require(success, "Low-level call failed"); // 可选：失败时回滚交易‌
    }
}

contract ReentrancyGuardTransient {
    bytes32 constant SLOT = 0;

    modifier lock() {
        assembly {
            if tload(SLOT) {
                revert(0, 0)
            }
            tstore(SLOT, 1)
        }
        _;
        // 无论函数体是否成功都清除锁
        assembly {
            tstore(SLOT, 0)
        }
    }

    // 21887 gas
    function test() external lock {
        // Ignore call error
        bytes memory b = "";
        (bool success, ) = msg.sender.call(b); // 捕获调用是否成功
        require(success, "Low-level call failed"); // 可选：失败时回滚交易‌
    }
}

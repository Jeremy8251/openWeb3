// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Enum {
    enum Status {
        Pending, //0
        Shipped, //1
        Accepted, //2
        Rejected, //3
        Canceled //4
    }

    Status public status;

    function get() public view returns (Status) {
        return status;
    }

    function set(Status _status) public {
        status = _status;
    }

    function cancel() public {
        status = Status.Canceled;
    }

    function reset() public {
        delete status;
    }
}

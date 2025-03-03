// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

contract TodoList {
    uint public taskCount = 0;
    struct Task {
        uint i;
        string taskname;
        bool status;
    }

    mapping(uint => Task) public tasks;
    event TaskCreated(uint id, string taskname, bool status);

    event TaskStatus(uint id, bool status);

    constructor() {
        createTask("Todo List Tutorial");
    }

    function createTask(string memory _taskname) public {
        taskCount++;
        tasks[taskCount] = Task(taskCount, _taskname, false);
        emit TaskCreated(taskCount, _taskname, false);
    }

    function toggleStatus(uint _id) public {
        Task memory _task = tasks[_id];

        _task.status = !_task.status;
        tasks[_id] = _task;
        emit TaskStatus(_id, _task.status);
    }
}

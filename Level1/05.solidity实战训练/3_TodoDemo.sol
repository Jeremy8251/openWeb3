// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
TodoList: 是类似便签一样功能的东西，记录我们需要做的事情，以及完成状态。 
1.需要完成的功能
创建任务
修改任务名称
    任务名写错的时候
修改完成状态：
    手动指定完成或者未完成
    自动切换
        如果未完成状态下，改为完成
        如果完成状态，改为未完成
获取任务 
2.思考代码内状态变量怎么安排？ 
思考 1：思考任务 ID 的来源？ 我们在传统业务里，这里的任务都会有一个任务 ID，在区块链里怎么实现？？ 答：传统业务里，ID 可以是数据库自动生成的，也可以用算法来计算出来的，比如使用雪花算法计算出 ID 等。在区块链里我们使用数组的 index 索引作为任务的 ID，也可以使用自增的整型数据来表示。 
思考 2: 我们使用什么数据类型比较好？ 答：因为需要任务 ID，如果使用数组 index 作为任务 ID。则数据的元素内需要记录任务名称，任务完成状态，所以元素使用 struct 比较好。 如果使用自增的整型作为任务 ID，则整型 ID 对应任务，使用 mapping 类型比较符合。 
3.演示代码
 
*/
contract Demo {
    struct Todo {
        string name;
        bool isCompleted;
    }
    Todo[] public list;

    // 创建任务,
    function create(string memory _name) external {
        list.push(Todo({name: _name, isCompleted: false}));
    }

    // 修改任务名称 ,35327 gas
    function modifyName(uint256 _index, string memory _name) external {
        // 方法1: 直接修改，修改一个属性时候比较省 gas
        list[_index].name = _name;
    }

    // 修改任务名称2,	35426 gas
    function modifyName2(uint256 _index, string memory _name) external {
        // 方法2: 先获取储存到 storage，在修改，在修改多个属性的时候比较省 gas
        Todo storage temp = list[_index];
        temp.name = _name;
    }

    // 修改完成状态1:手动指定完成或者未完成,53236 gas
    function modifyStatus(uint256 _index, bool _status) external {
        list[_index].isCompleted = _status;
    }

    // 修改完成状态2:自动切换 toggle,39134 gas
    function modifyStatus2(uint256 _index) external {
        list[_index].isCompleted = !list[_index].isCompleted;
    }

    // 获取任务1: memory : 2次拷贝
    // 8278 gas
    function get(
        uint256 _index
    ) external view returns (string memory _name, bool _status) {
        Todo memory todo = list[_index];
        return (todo.name, todo.isCompleted);
    }

    // 获取任务2: storage : 1次拷贝
    // 预期：get2 的 gas 费用比较低（相对 get1）
    //	8152 gas
    function get2(
        uint256 _index
    ) external view returns (string memory _name, bool _status) {
        Todo storage todo = list[_index];
        return (todo.name, todo.isCompleted);
    }
}

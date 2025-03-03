const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");
module.exports = buildModule("TodoListModule", (m) => {
  const todoList = m.contract("TodoList", []);

  return { todoList };
});

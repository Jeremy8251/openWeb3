// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("RccModule", (m) => {

  const token = m.contract("RccToken");
  console.log("token deployed to:", token.address);
  return { token };
});

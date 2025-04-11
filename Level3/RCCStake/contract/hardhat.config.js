require("@nomicfoundation/hardhat-verify");
require("@nomicfoundation/hardhat-toolbox");
// require("hardhat-deploy");
// require("hardhat-deploy-ethers");
require("@nomicfoundation/hardhat-ethers");
require("hardhat-gas-reporter");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  defaultNetwork: "hardhat", // 将默认网络设置为 hardhat
  gasReporter: {
    enabled: true,
    currency: "USD",
    gasPrice: 21,
  },
  networks: {
    hardhat: {},
    // localhost: {
    //   url: "http://127.0.0.1:8545", // 本地 Hardhat 网络的默认 RPC 地址
    //   accounts: ["0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78XXX"], // 使用本地账户的私钥
    //   deployerAddress: "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603XX", // 定义本地部署者地址
    // },
    holesky: {
      url: "https://api.zan.top/node/v1/eth/holesky/XXXXX", // 替换为实际的网络 RPC URL
      accounts: ["XXXX"], // 使用私钥部署合约
      deployerAddress: "XXXX", // 定义部署者地址
    },
  },
};

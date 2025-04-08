require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-ethers");
require("ethers");
require("web3");
require("dotenv").config();
/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  networks: {
    bscTestnet: {
      url: "https://data-seed-prebsc-1-s1.binance.org:8545",
      accounts: [process.env.PRIVATE_KEY],
    },
    sepolia: {
      url: "https://sepolia.infura.io/v3/d8ed0bd1de8242d998a1405b6932ab33",
      accounts: [process.env.PRIVATE_KEY],
    },
    holesky: {
      url: "https://api.zan.top/node/v1/eth/mainnet/XXXXXXXXXXXXXXXXXXXX",
      accounts: [process.env.PRIVATE_KEY],
    },
  },
  solidity: {
    optimizer: {
      enabled: false,
      runs: 50,
    },
    compilers: [
      {
        version: "0.4.18",
        settings: {
          evmVersion: "berlin",
        },
      },
      {
        version: "0.5.16",
        settings: {
          evmVersion: "berlin",
        },
      },
      {
        version: "0.6.6",
        settings: {
          evmVersion: "berlin",
          optimizer: {
            enabled: true,
            runs: 1,
          },
        },
      },
      {
        version: "0.6.12",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1,
          },
          // evmVersion: "shanghai"
        },
      },
      {
        version: "0.8.20",
        settings: {
          optimizer: {
            enabled: true,
            runs: 1,
          },
          // evmVersion: "shanghai"
        },
      },
    ],
  },
};

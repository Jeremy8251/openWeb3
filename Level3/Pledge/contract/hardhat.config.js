require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-ethers");
require("ethers");
require("web3");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
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

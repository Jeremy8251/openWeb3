require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

const HOLESHY_URL = process.env.HOLESHY_URL;
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const PRIVATE_KEY_1 = process.env.PRIVATE_KEY_1;
const ETHERSCAN_KEY = process.env.ETHERSCAN_KEY;
/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",
  networks: {
    hardhat: {},
    Holesky: {
      url: HOLESHY_URL,
      accounts: [PRIVATE_KEY, PRIVATE_KEY_1],
      chainId: 17000,
    },
  },
};

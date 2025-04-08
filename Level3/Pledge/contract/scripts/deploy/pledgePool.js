// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.

let oracleAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
let swapRouter = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";
let feeAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
let multiSignatureAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";

const { ethers } = require("hardhat");

async function main() {
  // const [deployerMax,,,,deployerMin] = await ethers.getSigners();
  const [deployerMin, , , , deployerMax] = await ethers.getSigners();

  console.log("Deploying contracts with the account:", deployerMin.address);
  const accountBalance = await ethers.provider.getBalance(deployerMax.address);
  console.log("Account balance:", accountBalance);

  const pledgePoolToken = await ethers.getContractFactory("PledgePool");
  const pledgeAddress = await pledgePoolToken
    .connect(deployerMin)
    .deploy(oracleAddress, swapRouter, feeAddress, multiSignatureAddress);

  console.log("pledgeAddress address:", await pledgeAddress.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

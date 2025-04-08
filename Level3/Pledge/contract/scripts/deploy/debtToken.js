// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.

const { ethers } = require("hardhat");

let tokenName = "spBTC_1";
let tokenSymbol = "spBTC_1";
let multiSignatureAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";

async function main() {
  // const [deployerMax,,,,deployerMin] = await ethers.getSigners();
  const [deployerMin, , , , deployerMax] = await ethers.getSigners();

  console.log("Deploying contracts with the account:", deployerMin.address);

  const accountBalance = await ethers.provider.getBalance(deployerMax.address);
  console.log("Account balance:", accountBalance);

  const debtToken = await ethers.getContractFactory("DebtToken");
  const DebtToken = await debtToken
    .connect(deployerMin)
    .deploy(tokenName, tokenSymbol, multiSignatureAddress);

  console.log("DebtToken address:", await DebtToken.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

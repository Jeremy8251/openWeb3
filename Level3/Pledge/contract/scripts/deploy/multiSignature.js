const { ethers } = require("hardhat");

let multiSignatureAddress = [
  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
  "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
  "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
];
let threshold = 2;

async function main() {
  const [deployerMax, , , , deployerMin] = await ethers.getSigners();

  console.log("Deploying contracts with the account:", deployerMax.address);

  const accountBalance = await ethers.provider.getBalance(deployerMax.address);
  console.log("Account balance:", accountBalance);
  console.log("Account balance (ETH):", ethers.formatEther(accountBalance));
  const multiSignatureToken = await ethers.getContractFactory("multiSignature");
  const multiSignature = await multiSignatureToken
    .connect(deployerMax)
    .deploy(multiSignatureAddress, threshold, {
      gasLimit: 3000000, // 设置 gas 限制
      gasPrice: ethers.parseUnits("10", "gwei"), // 设置 gas 价格
    });
  console.log("multiSignature address:", await multiSignature.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

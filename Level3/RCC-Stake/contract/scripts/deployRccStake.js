const { ethers } = require("hardhat");
async function main() {
  const rccStake = await ethers.getContractFactory("RccStake");
  const stake = await rccStake.deploy();
  console.log("RccStake deployed to:", await stake.getAddress());//0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

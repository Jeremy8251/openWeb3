const { ethers } = require("hardhat");
async function main() {
  const rccToken = await ethers.getContractFactory("RccToken");
  const token = await rccToken.deploy();
  console.log("RccToken deployed to:", await token.getAddress()); //0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

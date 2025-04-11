const { ethers } = require("hardhat");

async function main() {
  // 部署 RccToken 合约
  //   console.log("Deploying RccToken...");
  //   const RccToken = await ethers.getContractFactory("RccToken");
  //   const rccToken = await RccToken.deploy();
  //   rccTokenAddress = await rccToken.getAddress();

  // 部署 RccStake 合约
  //   console.log("Deploying RccStake...");
  //   const RccStake = await ethers.getContractFactory("RccStake");
  //   const rccStake = await RccStake.deploy({ gasLimit: 5000000 }); // 设置 Gas 限额
  //   const rccStakeAddress = await rccStake.getAddress();

  rccStakeAddress = "0x7bb5395608E0029b242aB24D125cb9f4C715d8Bd"; // 替换为实际地址
  rccTokenAddress = "0x8C35E99751A7c7C2c3EF32d2eAe89cC9CC93eE02"; // 替换为实际地址
  rccStake = await ethers.getContractAt("RccStake", rccStakeAddress);

  console.log("RccToken deployed to:", rccTokenAddress);
  console.log("RccStake deployed to:", rccStakeAddress);

  // 初始化 RccStake 合约
  console.log("Initializing RccStake...");
  const startBlock = 0; // 设置起始区块
  const endBlock = 100000; // 设置结束区块
  const rccPerBlock = ethers.parseEther("1"); // 每区块奖励 1 RCC
  await rccStake.initialize(
    rccTokenAddress,
    startBlock,
    endBlock,
    rccPerBlock,
    {
      gasLimit: 5000000, // 设置 Gas 限额
    }
  );
  console.log("RccStake initialized with RccToken address:", rccTokenAddress);

  console.log("Deployment completed successfully!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});

// // const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");
const { assert, expect } = require("chai");
const { ethers, deployments } = require("hardhat");
describe("RccToken Contract", async () => {
  let rccToken;
  // 每次测试前进行部署
  beforeEach(async () => {
    console.log("测试前进行部署");
    // 所有操作先进行这一部署合约,使用deploy文件夹下的js文件带有导出all标识的文件
    // await deployments.fixture(["all"]);
    let RccToken = await ethers.getContractFactory("RccToken");
    rccToken = await RccToken.deploy();
    console.log("RccToken deployed to:", await rccToken.getAddress());
  });
  it("test", async () => {
    const name = await rccToken.name();
    console.log("name = ", name);
    assert.equal(name, "RccToken");
  });
  //   // 获取第一个账户
  //   // firstAccount = (await getNamedAccounts()).deployerAccount;
  //   // secondAccount = (await getNamedAccounts()).secondAccount;
  //   console.log("测试前进行部署", deployments);
  // 合约部署信息
  // const RccToken = await deployments.get("RccToken");
  // console.log("进行部署", RccToken);
  // mockV3Aggregator = await deployments.get("MockV3Aggregator");
  // 获取合约
  // rccToken = await ethers.getContractAt("RccToken", RccToken.address);
  //   console.log("RccToken合约地址:", rcc.address);
  //   // fundMeSecondAccount = await ethers.getContract("FundMe", secondAccount); //连接新的合约地址
});
//   // 测试合约
//   // it("test", async () => {
//   //   await rccToken.waitForDeployment();
//   //   const name = await rccToken.name();
//   //   console.log("name = ", name);
//   //   // console.log("firstAccount = ", firstAccount);

//   //   assert.equal(name, "RccToken");
//   // });

// });

//   //   it("Should fail if the unlockTime is not in the future", async function () {
//   //     // We don't use the fixture here because we want a different deployment
//   //     const latestTime = await time.latest();
//   //     const Lock = await ethers.getContractFactory("Lock");
//   //     await expect(Lock.deploy(latestTime, { value: 1 })).to.be.revertedWith(
//   //       "Unlock time should be in the future"
//   //     );
//   //   });
//   // });

//   // describe("Withdrawals", function () {
//   //   describe("Validations", function () {
//   //     it("Should revert with the right error if called too soon", async function () {
//   //       const { lock } = await loadFixture(deployOneYearLockFixture);

//   //       await expect(lock.withdraw()).to.be.revertedWith(
//   //         "You can't withdraw yet"
//   //       );
//   //     });

//   //     it("Should revert with the right error if called from another account", async function () {
//   //       const { lock, unlockTime, otherAccount } = await loadFixture(
//   //         deployOneYearLockFixture
//   //       );

//   //       // We can increase the time in Hardhat Network
//   //       await time.increaseTo(unlockTime);

//   //       // We use lock.connect() to send a transaction from another account
//   //       await expect(lock.connect(otherAccount).withdraw()).to.be.revertedWith(
//   //         "You aren't the owner"
//   //       );
//   //     });

//   //     it("Shouldn't fail if the unlockTime has arrived and the owner calls it", async function () {
//   //       const { lock, unlockTime } = await loadFixture(deployOneYearLockFixture);

//   //       // Transactions are sent using the first signer by default
//   //       await time.increaseTo(unlockTime);

//   //       await expect(lock.withdraw()).not.to.be.reverted;
//   //     });
//   //   });

//   // describe("Events", function () {
//   //   it("Should emit an event on withdrawals", async function () {
//   //     const { lock, unlockTime, lockedAmount } = await loadFixture(
//   //       deployOneYearLockFixture
//   //     );

//   //     await time.increaseTo(unlockTime);

//   //     await expect(lock.withdraw())
//   //       .to.emit(lock, "Withdrawal")
//   //       .withArgs(lockedAmount, anyValue); // We accept any value as `when` arg
//   //   });
//   // });
// // });

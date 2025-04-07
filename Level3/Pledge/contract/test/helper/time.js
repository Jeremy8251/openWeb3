const web3 = require("web3");
const BN = web3.utils.BN;
const { ethers } = require("hardhat");
// async function latestBlock() {
//   const block = await web3.eth.getBlock("latest");
//   return new BN(block.number);
// }
async function latestBlock() {
  const block = await ethers.provider.getBlock("latest");
  console.log("block.number", ethers.toBigInt(block.number));
  return ethers.toBigInt(block.number);
}

async function latestBlockNum() {
  const block = await ethers.provider.getBlock("latest");
  return block.number; // 直接返回区块号
}
// async function latestBlockNum() {
//   const block = await web3.eth.getBlock("latest");
//   return new BN(block.number).toNumber();
// }

async function showBlock(msg = "") {
  const block = await ethers.provider.getBlock("latest");
  console.log(`${msg} at block number: ${block.number}`);
}
// async function showBlock() {
//   console.log("showBlock");
//   const block = await web3.eth.getBlock("latest");
//   console.log("Block number: " + new BN(block.number).toString());
// }

// async function showBlock(msg) {
//   const block = await web3.eth.getBlock("latest");
//   console.log(msg + " at block number: " + new BN(block.number).toString());
// }

async function stopAutoMine() {
  //stop auto mine or it will mess the block number
  network.provider.send("evm_setIntervalMining", [600000]);
  // await network.provider.send("evm_setAutomine", [false])
}

function advanceBlock() {
  // return network.provider.send("evm_mine", [new Date().getTime()])
  return network.provider.send("evm_mine", []);
}

// Advance the block to the passed height
async function advanceBlockTo(target) {
  // stop interval mint,set to 600s
  await stopAutoMine();
  // 确保 target 是 BigInt 类型
  if (typeof target !== "bigint") {
    target = ethers.toBigInt(target);
  }
  // if (!BN.isBN(target)) {
  //   target = new BN(target);
  // }

  const currentBlock = await latestBlock();
  const start = Date.now();
  let notified;
  if (target < currentBlock)
    throw Error(
      `Target block #(${target}) is lower than current block #(${currentBlock})`
    );
  while ((await latestBlock()) < target) {
    if (!notified && Date.now() - start >= 5000) {
      notified = true;
      console.log(
        `\
${colors.white.bgBlack("@openzeppelin/test-helpers")} ${colors.black.bgYellow(
          "WARN"
        )} advanceBlockTo: Advancing too ` +
          "many blocks is causing this test to be slow."
      );
    }

    await advanceBlock();
  }
  console.log("advanceBlockTo", target.toString());

  await showBlock("arrive");
}

async function latest() {
  const block = await ethers.provider.getBlock("latest");
  // console.log("block.timestamp", ethers.toBigInt(block.timestamp));
  return ethers.toBigInt(block.timestamp);
}

// Returns the time of the last mined block in seconds
// async function latest() {
//   console.log("latestBlock", web3);
//   const block = await web3.eth.getBlock("latest");
//   console.log("block", block);
//   return new BN(block.timestamp);
// }

async function increase(seconds) {
  await network.provider.send("evm_increaseTime", [seconds]);
  await advanceBlock();
}

module.exports = {
  advanceBlockTo,
  advanceBlock,
  latestBlock,
  latestBlockNum,
  showBlock,
  stopAutoMine,
  latest,
  increase,
};

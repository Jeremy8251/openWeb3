const { ethers } = require("hardhat");
const web3 = require("web3");
const BN = web3.utils.BN;

async function mockUniswap(minter, weth) {
  const UniswapV2Factory = await ethers.getContractFactory("UniswapV2Factory");
  let uniswapFactory = await UniswapV2Factory.deploy(minter.address);

  const UniswapV2Router02 = await ethers.getContractFactory(
    "UniswapV2Router02"
  );
  let uniswapRouter = await UniswapV2Router02.deploy(
    uniswapFactory.address,
    weth.address
  );
  return [uniswapRouter, uniswapFactory];
}
// targetToken: {address: string, price: int}
async function mockAddLiquidity(
  router,
  token0,
  token1,
  minter,
  deadline,
  amount0,
  amount1
) {
  // approve
  await token0
    .connect(minter)
    .approve(await router.getAddress(), BigInt(amount0));
  await token1
    .connect(minter)
    .approve(await router.getAddress(), BigInt(amount1));
  // add
  await router
    .connect(minter)
    .addLiquidity(
      await token0.getAddress(),
      await token1.getAddress(),
      BigInt(amount0),
      BigInt(amount1),
      BigInt(0),
      BigInt(0),
      minter.address,
      deadline
    );
}

async function mockSwap(
  router,
  token0,
  swapAmount,
  minAmount,
  path,
  minter,
  deadline
) {
  // approve
  await token0.connect(minter).approve(router.address, BigInt(swapAmount));
  // swap
  await router
    .connect(minter)
    .swapExactTokensForTokens(
      BigInt(swapAmount),
      BigInt(minAmount),
      path,
      minter.address,
      deadline
    );
}

module.exports = {
  mockUniswap,
  mockSwap,
  mockAddLiquidity,
};

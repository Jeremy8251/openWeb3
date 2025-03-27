// // const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");
const { assert, expect } = require("chai");
const { ethers } = require("hardhat");
const { BigNumber } = require("ethers");
describe("RccStake Contract", async () => {
  let rccStake, rccToken, owner, admin, user;
  // 每次测试前进行部署
  beforeEach(async () => {
    console.log("测试前进行部署");
    // 所有操作先进行这一部署合约,使用deploy文件夹下的js文件带有导出all标识的文件
    // await deployments.fixture(["all"]);
    [owner, admin, user] = await ethers.getSigners();
    //打印在一行
    console.log(
      "owner:",
      owner.address,
      "admin:",
      admin.address,
      "user:",
      user.address
    );
    // 部署 RCC Token 合约
    let RccToken = await ethers.getContractFactory("RccToken");
    rccToken = await RccToken.deploy();
    console.log("RccToken deployed to:", await rccToken.getAddress());

    // 部署 RCC Stake 合约
    let RccStake = await ethers.getContractFactory("RccStake");
    rccStake = await RccStake.deploy();
    console.log("RccStake deployed to:", await rccStake.getAddress());
    // 初始化 RCCStake 合约
    await rccStake.initialize(
      rccToken, // RCC 代币
      0, // startBlock
      100, // endBlock
      ethers.parseEther("1") // rccPerBlock
    );
    console.log("initialize");
    // 授予 admin 账户 ADMIN_ROLE
    const ADMIN_ROLE = await rccStake.ADMIN_ROLE();
    await rccStake.connect(owner).grantRole(ADMIN_ROLE, admin.address);
  });
  it("should allow an admin to set the RCC token address", async function () {
    console.log("调用 setRCC 函数");
    const rccTokenAddress = await rccToken.getAddress();
    console.log("rccToken:", rccTokenAddress);

    await expect(rccStake.connect(admin).setRCC(rccTokenAddress))
      .to.emit(rccStake, "SetRCC") // 验证是否触发了 SetRCC 事件
      .withArgs(rccTokenAddress); // 验证事件参数是否正确

    // 验证 RCC 地址是否正确设置
    const rccAddress = await rccStake.RCC();
    console.log("rccAddress:", rccAddress);
    expect(rccAddress).to.equal(rccTokenAddress);
  });

  it("should revert if a non-admin tries to set the RCC token address", async function () {
    // 非管理员调用 setRCC 应失败
    console.log("非管理员调用 setRCC 应失败");
    console.log("admin:", admin.address);
    console.log("ADMIN_ROLE:", await rccStake.ADMIN_ROLE());
    console.log(
      "error:",
      `AccessControlUnauthorizedAccount("${
        user.address
      }", "${await rccStake.ADMIN_ROLE()}")`
    );
    // 验证错误
    await expect(rccStake.connect(user).setRCC(await rccToken.getAddress()))
      .to.be.revertedWithCustomError(
        rccStake,
        "AccessControlUnauthorizedAccount"
      )
      .withArgs(user.address, await rccStake.ADMIN_ROLE());
  });

  // 验证初始化数据正确性
  it("should initialize the contract with the correct data", async function () {
    // 验证 RCC 地址是否正确设置
    const rccAddress = await rccStake.RCC();
    expect(rccAddress).to.equal(await rccToken.getAddress());

    // 验证开始区块是否正确设置
    const startBlock = await rccStake.startBlock();
    expect(startBlock).to.equal(0);

    // 验证结束区块是否正确设置
    const endBlock = await rccStake.endBlock();
    expect(endBlock).to.equal(100);

    // 验证每区块 RCC 奖励是否正确设置
    const rccPerBlock = await rccStake.rccPerBlock();
    expect(rccPerBlock).to.equal(ethers.parseEther("1"));
  });

  it("should allow an admin to set the start block", async function () {
    // 设置开始区块
    const startBlock = 100;
    await expect(rccStake.connect(admin).setStartBlock(startBlock))
      .to.emit(rccStake, "SetStartBlock")
      .withArgs(startBlock);
    // 验证开始区块是否正确设置
    const _startBlock = await rccStake.startBlock();
    expect(_startBlock).to.equal(startBlock);
  });

  //验证设置结束区块
  it("should allow an admin to set the end block", async function () {
    // 调试日志
    console.log("startBlock:", await rccStake.startBlock());
    console.log("endBlock:", await rccStake.endBlock());

    await rccStake.connect(admin).setStartBlock(100);

    // 验证结束区块是否大于开始区块，require(endBlock > startBlock)
    await expect(rccStake.connect(admin).setEndBlock(50)).to.be.revertedWith(
      "end block must be larger than start block"
    );

    // 设置结束区块
    const endBlock = 200;
    await expect(rccStake.connect(admin).setEndBlock(endBlock))
      .to.emit(rccStake, "SetEndBlock")
      .withArgs(endBlock);
    // 验证结束区块是否正确设置
    const _endBlock = await rccStake.endBlock();
    expect(_endBlock).to.equal(endBlock);
  });

  // 验证暂停解押和恢复解押
  it("should allow an admin to pause and unpause the contract", async function () {
    // 暂停解押
    await expect(rccStake.connect(admin).pauseWithdraw())
      .to.emit(rccStake, "PauseWithdraw")
      .withArgs();
    // 验证是否暂停
    // await rccStake.pauseWithdraw();
    const paused = await rccStake.withdrawPaused();
    expect(paused).to.equal(true);

    // 恢复解押
    await expect(rccStake.connect(admin).unpauseWithdraw())
      .to.emit(rccStake, "UnpauseWithdraw")
      .withArgs();
    // 验证是否恢复
    const _paused = await rccStake.withdrawPaused();
    expect(_paused).to.equal(false);
  });

  // 验证暂停领取奖励和恢复领取奖励
  it("should allow an admin to pause and unpause reward claiming", async function () {
    // 暂停领取奖励
    await expect(rccStake.connect(admin).pauseClaim())
      .to.emit(rccStake, "PauseClaim")
      .withArgs();
    // 验证是否暂停
    const paused = await rccStake.claimPaused();
    expect(paused).to.equal(true);

    // 恢复领取奖励
    await expect(rccStake.connect(admin).unpauseClaim())
      .to.emit(rccStake, "UnpauseClaim")
      .withArgs();
    // 验证是否恢复
    const _paused = await rccStake.claimPaused();
    expect(_paused).to.equal(false);
  });

  // 验证设置每区块分发的 RCC 奖励数量
  it("should allow an admin to set the RCC per block", async function () {
    // 设置每区块 RCC 奖励数量
    const rccPerBlock = ethers.parseEther("2");
    await expect(rccStake.connect(admin).setRCCPerBlock(rccPerBlock))
      .to.emit(rccStake, "SetRCCPerBlock")
      .withArgs(rccPerBlock);
    // 验证每区块 RCC 奖励数量是否正确设置
    const _rccPerBlock = await rccStake.rccPerBlock();
    expect(_rccPerBlock).to.equal(rccPerBlock);
  });

  // 添加第一个 ETH 池
  it("should allow admin to add the first ETH pool", async function () {
    // 获取当前区块号
    let currentBlock = await ethers.provider.getBlockNumber();
    console.log("currentBlock before transaction:", currentBlock);
    console.log("startBlock:", await rccStake.startBlock());

    await expect(
      rccStake.connect(admin).addPool(
        ethers.ZeroAddress, // ETH 池地址
        10, // poolWeight
        ethers.parseEther("0.1"), // minDepositAmount
        10, // unstakeLockedBlocks
        false // withUpdate
      )
    )
      .to.emit(rccStake, "AddPool")
      .withArgs(
        ethers.ZeroAddress,
        10,
        currentBlock + 1, // lastRewardBlock，在交易执行后，区块号确实增加了 1
        ethers.parseEther("0.1"),
        10
      );
    currentBlock = await ethers.provider.getBlockNumber();
    console.log("block.number during transaction:", currentBlock);
    // 验证池的数量
    const poolLength = await rccStake.poolLength();
    expect(poolLength).to.equal(1);

    // 验证池的参数
    const pool = await rccStake.pool(0);
    expect(pool.stTokenAddress).to.equal(ethers.ZeroAddress);
    expect(pool.poolWeight).to.equal(10);
    expect(pool.minDepositAmount).to.equal(ethers.parseEther("0.1"));
    expect(pool.unstakeLockedBlocks).to.equal(10);
  });

  // 验证更新指定质押池updatePool
  it("should update the specified pool", async function () {
    // 添加第一个 ETH 池
    await rccStake.connect(admin).addPool(
      ethers.ZeroAddress, // ETH 池地址
      10, // poolWeight
      ethers.parseEther("0.1"), // minDepositAmount
      10, // unstakeLockedBlocks
      false // withUpdate
    );

    // 更新指定池
    await expect(rccStake.connect(admin).updatePool(0, 1, 5))
      .to.emit(rccStake, "UpdatePoolInfo")
      .withArgs(0, 1, 5);
  });

  // 验证设置Pool权重
  it("should allow an admin to set the pool weight", async function () {
    // 添加第一个 ETH 池
    await rccStake.connect(admin).addPool(
      ethers.ZeroAddress, // ETH 池地址
      10, // poolWeight
      ethers.parseEther("0.1"), // minDepositAmount
      10, // unstakeLockedBlocks
      false // withUpdate
    );
    // console.log("BigNumber:", ethers.toBigInt(0));
    // 计算总权重
    let totalWeight = ethers.toBigInt(0);
    const poolLength = await rccStake.poolLength();
    for (let i = 0; i < poolLength; i++) {
      const poolWeight = (await rccStake.pool(i)).poolWeight;
      // console.log("poolWeight:", poolWeight);
      totalWeight = totalWeight + poolWeight;
    }
    // console.log("totalWeight before update:", totalWeight.toString());

    // 计算新的总权重
    const oldWeight = (await rccStake.pool(0)).poolWeight;
    const newWeight = ethers.toBigInt(20);
    totalWeight = totalWeight - oldWeight + newWeight;
    // console.log("totalWeight after update:", totalWeight.toString());

    // 设置池权重
    await expect(rccStake.connect(admin).setPoolWeight(0, 20, true))
      .to.emit(rccStake, "SetPoolWeight")
      .withArgs(0, 20, totalWeight);
  });

  //验证从 _from 到 _to 区块之间的奖励倍数
  it("should return the reward multiplier from _from to _to blocks", async function () {
    // 设置开始区块
    await rccStake.connect(admin).setStartBlock(100);
    // 设置结束区块
    await rccStake.connect(admin).setEndBlock(200);

    // 设置每区块 RCC 奖励数量
    await rccStake.connect(admin).setRCCPerBlock(ethers.parseEther("1"));

    // 获取奖励倍数,因为from大于开始区块，将开始区块赋值给from，所以从100-180区块之间的奖励倍数
    const multiplier = await rccStake.getMultiplier(150, 180);
    console.log("multiplier:", multiplier.toString());
    // 计算区块间的奖励倍数
    const blocks = BigInt(180 - 100);
    const rccPerBlock = BigInt(await rccStake.rccPerBlock());
    console.log("rccPerBlock:", rccPerBlock);
    const reward = blocks * rccPerBlock;
    console.log("reward:", reward.toString());
    expect(reward).to.equal(BigInt(multiplier));
  });

  // 验证用户质押的代币数量
  it("should return the amount of tokens staked by a user", async function () {
    let currentBlock = await ethers.provider.getBlockNumber();
    console.log("Current block:", currentBlock);
    // 添加第一个质押池
    await rccStake.connect(admin).addPool(
      ethers.ZeroAddress, // ETH 池地址
      10, // poolWeight
      ethers.parseEther("0.1"), // minDepositAmount
      10, // unstakeLockedBlocks
      false // withUpdate
    );

    // 检查质押池数量
    const poolLength = await rccStake.poolLength();
    console.log("Pool length:", poolLength.toString());
    expect(poolLength).to.equal(1);

    await rccStake.connect(user).depositETH({ value: ethers.parseEther("1") });

    // 用户发起解押请求
    await rccStake.connect(user).unstake(0, ethers.parseEther("0.5"));

    // 模拟区块前进，确保解锁区块号已达到
    for (let i = 0; i < 10; i++) {
      await ethers.provider.send("evm_mine", []);
    }
    currentBlock = await ethers.provider.getBlockNumber();
    console.log("Current block:", currentBlock);

    // 获取用户质押的代币数量
    const [requestAmount, pendingWithdrawAmount] =
      await rccStake.withdrawAmount(0, user.address);
    // 打印返回值
    console.log("Request Amount:", requestAmount.toString());
    console.log("Pending Withdraw Amount:", pendingWithdrawAmount.toString());

    // 验证返回值
    expect(requestAmount).to.equal(ethers.parseEther("0.5")); // 解押请求总量
    expect(pendingWithdrawAmount).to.equal(ethers.parseEther("0.5")); // 待解押总量
  });
});

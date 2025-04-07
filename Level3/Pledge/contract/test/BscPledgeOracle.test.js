const { expect } = require("chai");
const { ethers } = require("hardhat");
const { show } = require("./helper/meta.js");
const BN = require("bn.js");

describe("BscPledgeOracle ContracbusdAddress", function () {
  let BscPledgeOracle,
    bscPledgeOracle,
    busdAddress,
    btcAddress,
    multiSignature,
    minter,
    addr1,
    addr2;
  const threshold = 2; // 签名阈值

  beforeEach(async () => {
    [minter, addr1, addr2] = await ethers.getSigners();

    // 部署 multiSignature 合约
    const MultiSignature = await ethers.getContractFactory("multiSignature");
    multiSignature = await MultiSignature.deploy(
      [minter.address, addr1.address, addr2.address],
      threshold
    );
    const multiSignatureAddress = await multiSignature.getAddress();
    console.log("MultiSignature Address:", multiSignatureAddress);

    // 部署 BscPledgeOracle 合约
    const BscPledgeOracle = await ethers.getContractFactory("BscPledgeOracle");
    bscPledgeOracle = await BscPledgeOracle.deploy(multiSignatureAddress);
    console.log("bscPledgeOracle Address:", await bscPledgeOracle.getAddress());

    const DebtToken = await ethers.getContractFactory("DebtToken");
    busdToken = await DebtToken.deploy("BEP20Token", "BEP20", minter.address);
    busdAddress = await busdToken.getAddress();

    const BtcToken = await ethers.getContractFactory("DebtToken");
    btcToken = await BtcToken.deploy("BtcToken", "Btc", minter.address);
    btcAddress = await btcToken.getAddress();
  });

  async function createMultiSignature() {
    const msgHash = await multiSignature.getApplicationHash(
      minter.address,
      await bscPledgeOracle.getAddress()
    );
    console.log("msgHash:", msgHash);
    // 创建多签申请
    await multiSignature
      .connect(minter)
      .createApplication(await bscPledgeOracle.getAddress());
    await multiSignature.connect(minter).signApplication(msgHash);
    await multiSignature.connect(addr1).signApplication(msgHash); // 达到阈值
  }

  it("should set and get the price of an asset", async function () {
    const price = ethers.parseUnits("100", 18);

    await createMultiSignature();

    // 设置价格
    await bscPledgeOracle.connect(minter).setPrice(busdAddress, price);
    // 获取价格
    const fetchedPrice = await bscPledgeOracle.getPrice(busdAddress);
    expect(fetchedPrice).to.equal(price);
  });

  it("can not set price without authorization", async function () {
    await expect(
      bscPledgeOracle.connect(addr1).setPrice(busdAddress, 100000)
    ).to.be.revertedWith("multiSignatureClient : This tx is not aprroved");
  });

  it("Admin set price operation", async function () {
    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(0).toString()
    );

    await createMultiSignature();

    await bscPledgeOracle.connect(minter).setPrice(busdAddress, 100000000);
    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(100000000).toString()
    );
  });

  it("Administrators set prices in batches", async function () {
    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(0).toString()
    );
    expect(await bscPledgeOracle.getPrice(btcAddress)).to.equal(
      BigInt(0).toString()
    );
    let busdIndex = new BN(busdAddress.substring(2), 16).toString(10);
    let btcIndex = new BN(btcAddress.substring(2), 16).toString(10);

    await createMultiSignature();

    await bscPledgeOracle
      .connect(minter)
      .setPrices([busdIndex, btcIndex], [100, 100]);
    expect(await bscPledgeOracle.getUnderlyingPrice(0)).to.equal(
      BigInt(100).toString()
    );
    expect(await bscPledgeOracle.getUnderlyingPrice(1)).to.equal(
      BigInt(100).toString()
    );
  });

  it("Get price according to INDEX", async function () {
    await createMultiSignature();

    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(0).toString()
    );
    let underIndex = new BN(busdAddress.substring(2), 16).toString(10);
    await bscPledgeOracle
      .connect(minter)
      .setUnderlyingPrice(underIndex, 100000000);
    expect(await bscPledgeOracle.getUnderlyingPrice(underIndex)).to.equal(
      BigInt(100000000).toString()
    );
  });

  it("Set price according to INDEX", async function () {
    await createMultiSignature();
    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(0).toString()
    );
    let underIndex = new BN(busdAddress.substring(2), 16).toString(10);
    await bscPledgeOracle
      .connect(minter)
      .setUnderlyingPrice(underIndex, 100000000);
    expect(await bscPledgeOracle.getPrice(busdAddress)).to.equal(
      BigInt(100000000).toString()
    );
  });

  it("Set AssetsAggregator", async function () {
    await createMultiSignature();
    let arrData = await bscPledgeOracle.getAssetsAggregator(busdAddress);
    show(arrData[0]);
    expect(arrData[0]).to.equal("0x0000000000000000000000000000000000000000");
    await bscPledgeOracle
      .connect(minter)
      .setAssetsAggregator(busdAddress, btcAddress, 18);
    let data = await bscPledgeOracle.getAssetsAggregator(busdAddress);
    expect(data[0]).to.equal(btcAddress);
    expect(data[1]).to.equal(18);
  });
});

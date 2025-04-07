const { expect } = require("chai");
const { show } = require("./helper/meta.js");

describe("DebtToken", function () {
  let debtToken;
  const threshold = 2; // 签名阈值
  beforeEach(async () => {
    [minter, addr1, addr2, addr3, _] = await ethers.getSigners();
    // 部署 multiSignature 合约
    const MultiSignature = await ethers.getContractFactory("multiSignature");
    multiSignature = await MultiSignature.deploy(
      [minter.address, addr1.address, addr2.address],
      threshold
    );
    const multiSignatureAddress = await multiSignature.getAddress();
    console.log("MultiSignature Address:", multiSignatureAddress);

    const DebtToken = await ethers.getContractFactory("DebtToken");
    debtToken = await DebtToken.deploy(
      "spBUSD_1",
      "spBUSD_1",
      multiSignatureAddress
    );
  });

  async function createMultiSignature() {
    const msgHash = await multiSignature.getApplicationHash(
      minter.address,
      await debtToken.getAddress()
    );
    console.log("msgHash:", msgHash);
    // 创建多签申请
    await multiSignature
      .connect(minter)
      .createApplication(await debtToken.getAddress());
    await multiSignature.connect(minter).signApplication(msgHash);
    await multiSignature.connect(addr2).signApplication(msgHash); // 达到阈值
    await multiSignature.connect(addr1).signApplication(msgHash); // 达到阈值
  }

  it("check if mint right", async function () {
    let amount = await debtToken.totalSupply();
    show({ amount });
    expect(await debtToken.totalSupply()).to.equal(BigInt(0).toString());
    // mint
    await createMultiSignature(); //签名
    await debtToken.addMinter(minter.address);
    await debtToken
      .connect(minter)
      .mint(minter.address, BigInt(100000000 * 1e18));
    expect(await debtToken.balanceOf(minter.address)).to.equal(
      BigInt(100000000 * 1e18).toString()
    );
  });

  it("can not mint without authorization", async function () {
    await expect(
      debtToken.connect(addr1).mint(addr2.address, 100000)
    ).to.be.revertedWith("Token: caller is not the minter");
  });

  it("can not add minter by others", async function () {
    await createMultiSignature(); //签名
    // 非授权用户尝试调用 addMinter
    await expect(
      debtToken.connect(addr1).addMinter(addr1.address)
    ).to.be.revertedWith("multiSignatureClient : This tx is not aprroved");
    // ).to.be.revertedWith("Ownable: caller is not the owner");
  });

  it("after addMinter by owner, mint by minter should succeed", async function () {
    await createMultiSignature(); //签名
    await debtToken.connect(minter).addMinter(addr1.address);
    expect(await debtToken.balanceOf(addr2.address)).to.equal(0);

    await debtToken.connect(addr1).mint(addr2.address, 100000);
    expect(await debtToken.balanceOf(addr2.address)).to.equal(100000);

    await debtToken.connect(addr1).mint(addr2.address, 10000);
    expect(await debtToken.balanceOf(addr2.address)).to.equal(110000);
  });

  it("after delMinter, pre minter can not mint", async function () {
    await createMultiSignature(); //签名
    await debtToken.connect(minter).addMinter(addr1.address);
    await debtToken.connect(addr1).mint(addr2.address, 100000);
    expect(await debtToken.balanceOf(addr2.address)).to.equal(100000);
    // delete
    await debtToken.connect(minter).delMinter(addr1.address);
    await expect(
      debtToken.connect(addr1).mint(addr2.address, 100000)
    ).to.be.revertedWith("Token: caller is not the minter");
  });

  it("isMinter and getMinterLength should work well", async function () {
    expect(await debtToken.connect(minter).getMinterLength()).to.equal(0);

    await createMultiSignature(); //签名

    await debtToken.connect(minter).addMinter(addr1.address);
    await debtToken.connect(minter).addMinter(addr2.address);
    expect(await debtToken.connect(minter).getMinterLength()).to.equal(2);

    expect(await debtToken.isMinter(addr1.address)).to.equal(true);
    expect(await debtToken.isMinter(addr2.address)).to.equal(true);
    expect(await debtToken.isMinter(minter.address)).to.equal(false);

    await debtToken.connect(minter).delMinter(addr2.address);
    expect(await debtToken.connect(minter).getMinterLength()).to.equal(1);
    expect(await debtToken.isMinter(addr2.address)).to.equal(false);
  });

  it("getMinter should work well", async function () {
    await createMultiSignature(); //签名

    await debtToken.connect(minter).addMinter(addr1.address);
    await debtToken.connect(minter).addMinter(addr2.address);
    expect(await debtToken.getMinter(0)).to.equal(addr1.address);
    expect(await debtToken.getMinter(1)).to.equal(addr2.address);
    await expect(debtToken.getMinter(2)).to.be.revertedWith(
      "Token: index out of bounds"
    );

    await debtToken.connect(minter).delMinter(addr1.address);
    expect(await debtToken.getMinter(0)).to.equal(addr2.address);
    await expect(debtToken.getMinter(1)).to.be.revertedWith(
      "Token: index out of bounds"
    );
  });

  it("totalSupply should work well", async function () {
    await createMultiSignature(); //签名
    await debtToken.connect(minter).addMinter(addr1.address);
    await debtToken.connect(addr1).mint(addr2.address, 200);
    expect(await debtToken.totalSupply()).to.equal("200");
  });
});

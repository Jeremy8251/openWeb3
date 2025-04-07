const { expect } = require("chai");

describe("multiSignature", function () {
  let multiSignature, owner1, owner2, owner3, nonOwner;
  const threshold = 2; // 签名阈值
  beforeEach(async () => {
    [owner1, owner2, owner3, nonOwner] = await ethers.getSigners();

    // 部署 multiSignature 合约
    const MultiSignature = await ethers.getContractFactory("multiSignature");
    multiSignature = await MultiSignature.deploy(
      [owner1.address, owner2.address, owner3.address], // 签名者列表
      threshold // 签名阈值
    );
  });
  it("should deploy the contract with correct owners and threshold", async function () {
    // 验证签名者列表
    expect(await multiSignature.signatureOwners(0)).to.equal(owner1.address);
    expect(await multiSignature.signatureOwners(1)).to.equal(owner2.address);
    expect(await multiSignature.signatureOwners(2)).to.equal(owner3.address);

    // 验证阈值
    expect(await multiSignature.threshold()).to.equal(threshold);
  });

  it("should allow owners to create an application", async function () {
    const tx = await multiSignature
      .connect(owner1)
      .createApplication(owner2.address);
    await tx.wait();

    const msghash = await multiSignature.getApplicationHash(
      owner1.address,
      owner2.address
    );
    const applicationCount = await multiSignature.getApplicationCount(msghash);

    expect(applicationCount).to.equal(1);
  });

  it("should allow owners to sign an application", async function () {
    const msghash = await multiSignature.getApplicationHash(
      owner1.address,
      owner2.address
    );

    // 创建申请
    await multiSignature.connect(owner1).createApplication(owner2.address);

    // 签名申请
    await multiSignature.connect(owner1).signApplication(msghash);

    const [applicant, signatures] = await multiSignature.getApplicationInfo(
      msghash,
      0
    );
    expect(applicant).to.equal(owner1.address);
    expect(signatures).to.include(owner1.address);
  });

  it("should not allow non-owners to sign an application", async function () {
    const msghash = await multiSignature.getApplicationHash(
      owner1.address,
      owner2.address
    );

    // 创建申请
    await multiSignature.connect(owner1).createApplication(owner2.address);

    // 非签名者尝试签名
    await expect(
      multiSignature.connect(nonOwner).signApplication(msghash)
    ).to.be.revertedWith(
      "Multiple Signature : caller is not in the ownerList!"
    );
  });

  it("should validate signatures when threshold is met", async function () {
    const msghash = await multiSignature.getApplicationHash(
      owner1.address,
      owner2.address
    );

    // 创建申请
    await multiSignature.connect(owner1).createApplication(owner2.address);

    // 签名申请
    await multiSignature.connect(owner1).signApplication(msghash);
    await multiSignature.connect(owner2).signApplication(msghash);

    // 验证签名是否有效
    const validIndex = await multiSignature.getValidSignature(msghash, 0);
    expect(validIndex).to.equal(1); // 第一个申请满足阈值
  });

  it("should not validate signatures if threshold is not met", async function () {
    const msghash = await multiSignature.getApplicationHash(
      owner1.address,
      owner2.address
    );

    // 创建申请
    await multiSignature.connect(owner1).createApplication(owner2.address);

    // 仅一个签名
    await multiSignature.connect(owner1).signApplication(msghash);

    // 验证签名是否有效
    const validIndex = await multiSignature.getValidSignature(msghash, 0);
    expect(validIndex).to.equal(0); // 未满足阈值
  });
});

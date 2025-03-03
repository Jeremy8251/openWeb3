const { expect } = require("chai");
const hre = require("hardhat");
describe("Shipping", function () {
  let shippingContract;
  before(async () => {
    // 生成合约实例并且复用
    shippingContract = await hre.ethers.deployContract("Shipping", []);
  });

  it("should return the status Pending", async () => {
    expect(await shippingContract.Status()).to.equal("Pending");
  });

  it("should return the status Shipped", async () => {
    // Calling the Shipped() function
    await shippingContract.Shipped();
    // Checking if the status is Shipped
    expect(await shippingContract.Status()).to.equal("Shipped");
  });

  it("should return correct event description", async () => {
    await expect(shippingContract.Delivered()) // 验证事件是否被触发
      .to.emit(shippingContract, "LogNewAlert") // 验证事件的参数是否符合预期
      .withArgs("Your package has arrived");
  });
});

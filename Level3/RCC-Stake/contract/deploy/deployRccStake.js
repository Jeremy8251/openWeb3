/**
 * 部署函数
 * @param {Object} hre - Hardhat 运行时环境
 */
module.exports = async (hre) => {
  // 打印日志，表示部署函数正在执行
  console.log("Deploying RccStake contract");
  // 获取 Hardhat 运行时环境中的 deployments 对象
  const deployment = hre.deployments;

  // npx hardhat deploy --network Holesky 会生成deployments文件夹，再次运行会复用原有合约地址
  // 但如果修改了合约，需要重新部署
  // 可以通过npx hardhat deploy --network Holesk --reset 重置合约地址
  // 部署 RccToken 合约
  await deployment.deploy("RccStake", {
    from: hre.network.config.deployerAddress, // 使用配置中的部署者地址
    args: [], // 合约构造函数参数
    log: true, // 打印部署日志
  });

  console.log("RccStake deployed successfully!");
};

module.exports.tags = ["all", "RccToken"];

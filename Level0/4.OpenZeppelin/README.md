# Sample Hardhat Project

This project demonstrates a basic Hardhat use case. It comes with a sample contract, a test for that contract, and a Hardhat Ignition module that deploys that contract.

Try running some of the following tasks:

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/Lock.js
```

```项目运行过程
C:\My Documents\Project\Learn\OpenWeb3\3.hardhat-project>npm init
npm install --save-dev hardhat
npx hardhat init
npx hardhat node
npx hardhat compile
 npx hardhat ignition deploy ./ignition/modules/Lock.js --network localhost
 npx hardhat ignition deploy ignition/modules/Shipping.js --network localhost

 Deployed Addresses
ShippingModule#Shipping - 0x5FbDB2315678afecb367f032d93F642f64180aa3


 npx hardhat test test/Shipping.js
```

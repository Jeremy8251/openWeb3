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



C:\My Documents\Project\openWeb3\Level0\7.todolist>npx hardhat ignition deploy ./ignition/modules/TodoList.js --network Holesky
√ Confirm deploy to network Holesky (17000)? ... yes
Hardhat Ignition 🚀

Resuming existing deployment from .\ignition\deployments\chain-17000

Deploying [ TodoListModule ]

Batch #1
  Executed TodoListModule#TodoList

[ TodoListModule ] successfully deployed 🚀

Deployed Addresses

TodoListModule#TodoList - 0xe63226B6c87b347fD7F902581Be3b0BE7c2Ee3Dd

https://holesky.etherscan.io/address/0xe63226B6c87b347fD7F902581Be3b0BE7c2Ee3Dd#code
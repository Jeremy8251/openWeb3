### **PledgePool 合约场景与作用**

`PledgePool` 是一个基于区块链的借贷协议合约，旨在为用户提供去中心化的借贷服务。以下是对整个项目的功能、使用场景和作用的详细分析：

---

### **1. 项目核心功能**

#### **1.1 借贷池的创建与管理**
- **功能**：
  - 管理多个借贷池，每个池都有独立的参数（如利率、最大供应量、抵押率等）。
  - 支持创建新的借贷池，设置借贷相关的参数。
- **作用**：
  - 为用户提供灵活的借贷市场，允许用户根据需求选择不同的借贷池。

#### **1.2 用户存款与借款**
- **功能**：
  - 用户可以存入资金（`depositLend`），成为出借人，赚取利息。
  - 用户可以质押资产（`depositBorrow`），成为借款人，获取流动性。
- **作用**：
  - 为用户提供资金存取和借贷的功能，满足流动性需求。

#### **1.3 结算与清算**
- **功能**：
  - 支持借贷池的结算（`settle`），计算借贷金额和利息。
  - 支持清算（`liquidate`），当借款人的抵押品价值不足时，触发清算。
- **作用**：
  - 确保借贷池的资金安全，防止借款人违约。

#### **1.4 代币化的借贷权益**
- **功能**：
  - 出借人存款后会收到 `spToken`（代表存款权益）。
  - 借款人质押后会收到 `jpToken`（代表借款权益）。
- **作用**：
  - 通过代币化的方式，记录用户的借贷权益，方便后续的流转和管理。

#### **1.5 手续费与收益管理**
- **功能**：
  - 支持设置借贷手续费（`lendFee` 和 `borrowFee`）。
  - 手续费会转入指定的 `feeAddress`。
- **作用**：
  - 为协议提供收益来源，支持协议的可持续发展。

---

### **2. 使用场景**

#### **2.1 去中心化借贷市场**
- **场景**：
  - 用户希望通过质押资产（如 BTC）借入稳定币（如 BUSD）。
  - 用户希望通过存入稳定币赚取利息。
- **实现**：
  - 借款人调用 `depositBorrow` 质押资产。
  - 出借人调用 `depositLend` 存入资金。

#### **2.2 流动性管理**
- **场景**：
  - 用户希望在资金池中存入资金，并在需要时提取本金和利息。
- **实现**：
  - 用户调用 `withdrawLend` 提取本金和利息。

#### **2.3 风险管理与清算**
- **场景**：
  - 当借款人的抵押品价值不足时，触发清算，保护出借人的资金安全。
- **实现**：
  - 系统调用 `liquidate` 清算借款人的抵押品。

#### **2.4 收益分配**
- **场景**：
  - 协议希望通过手续费获取收益，用于维护和发展。
- **实现**：
  - 手续费会自动转入 `feeAddress`。

---

### **3. 合约模块与作用**

#### **3.1 数据结构**
- **`PoolBaseInfo`**：
  - 存储每个借贷池的基本信息（如利率、最大供应量、抵押率等）。
- **`PoolDataInfo`**：
  - 存储每个借贷池的动态数据（如结算金额、清算金额等）。
- **`LendInfo` 和 `BorrowInfo`**：
  - 跟踪每个用户在借贷池中的存款和借款状态。

#### **3.2 核心函数**
- **借贷相关**：
  - `depositLend`：用户存入资金，成为出借人。
  - `depositBorrow`：用户质押资产，成为借款人。
  - `withdrawLend`：出借人提取本金和利息。
  - `withdrawBorrow`：借款人提取剩余的保证金。
- **结算与清算**：
  - `settle`：结算借贷池，计算借贷金额和利息。
  - `liquidate`：清算借贷池，处理违约情况。
- **管理相关**：
  - `createPoolInfo`：创建新的借贷池。
  - `setFee`：设置借贷手续费。
  - `setPause`：暂停或恢复合约。

#### **3.3 事件**
- **核心事件**：
  - `DepositLend` 和 `DepositBorrow`：记录用户的存款和借款操作。
  - `RefundLend` 和 `RefundBorrow`：记录用户的退款操作。
  - `ClaimLend` 和 `ClaimBorrow`：记录用户的权益领取操作。
  - `WithdrawLend` 和 `WithdrawBorrow`：记录用户的提取操作。
  - `StateChange`：记录借贷池的状态变化。

---

### **4. 项目作用**

#### **4.1 为用户提供去中心化的借贷服务**
- 用户可以通过存款赚取利息，或通过质押资产获取流动性。
- 借贷过程完全去中心化，无需依赖第三方机构。

#### **4.2 提高资金利用率**
- 通过借贷池的设计，聚合用户的资金，提高资金利用率。

#### **4.3 提供风险管理机制**
- 通过清算机制，保护出借人的资金安全。
- 通过手续费机制，为协议提供收益来源。

#### **4.4 支持代币化的权益管理**
- 通过 `spToken` 和 `jpToken`，记录用户的借贷权益，方便后续的流转和管理。

---

### **5. 项目优势与改进建议**

#### **优势**
1. **模块化设计**：
   - 借贷池的设计灵活，支持多种资产和参数配置。
2. **去中心化**：
   - 所有操作均通过智能合约执行，无需信任第三方。
3. **风险管理**：
   - 支持清算机制，保护出借人的资金安全。

#### **改进建议**
1. **优化清算机制**：
   - 增加清算奖励，激励用户参与清算。
2. **支持更多资产**：
   - 增加对更多代币的支持，扩大用户群体。
3. **用户体验优化**：
   - 提供前端界面，方便用户操作。
4. **安全性审计**：
   - 对合约进行第三方安全审计，确保资金安全。

---

### **6. 当前项目的设计支持多种资产**
从 `PledgePool` 和 `AddressPrivileges` 合约的设计来看，项目已经支持多种资产的借贷。以下是关键点：

#### **6.1 借贷池的灵活性**
- **`createPoolInfo` 方法**：
  - 借贷池的创建允许指定不同的借款代币（`borrowToken`）和出借代币（`lendToken`）。
  - 这意味着可以创建一个以 BTC 作为抵押品、借出 USDT 的借贷池。

#### **6.2 权限管理**
- **`AddressPrivileges` 合约**：
  - 通过 `addMinter` 和 `delMinter` 方法，可以动态添加或移除代币的铸币权限。
  - 这意味着可以为不同的代币（如 BTC 或 USDT）设置铸币权限，支持多种资产的借贷。

#### **6.3 代币化的权益管理**
- **`spToken` 和 `jpToken`**：
  - 出借人和借款人的权益通过代币化（`spToken` 和 `jpToken`）进行管理。
  - 这些代币的灵活性允许支持多种资产的权益记录。

---

### **7. 使用 BTC 抵押借款的实现**

#### **7.1 创建 BTC 抵押池**
管理员可以调用 `createPoolInfo` 方法，创建一个以 BTC 作为抵押品、借出 USDT 的借贷池：

```javascript
await pledgePool.createPoolInfo(
  1680000000, // settleTime
  1685000000, // endTime
  5000000,    // interestRate (5%)
  ethers.utils.parseEther("1000"), // maxSupply (USDT)
  150000000,  // mortgageRate (150%)
  btcToken.address, // 借款代币 (BTC)
  usdtToken.address, // 出借代币 (USDT)
  spToken.address,   // 存款权益代币
  jpToken.address,   // 借款权益代币
  110000000          // autoLiquidateThreshold (110%)
);
```

---

#### **7.2 借款人质押 BTC**
借款人调用 `depositBorrow` 方法，质押 BTC 并借出 USDT：

```javascript
await btcToken.approve(pledgePool.address, ethers.utils.parseEther("1")); // 质押 1 BTC
await pledgePool.depositBorrow(0, ethers.utils.parseEther("1")); // 借款池 ID 为 0
```

---

#### **7.3 出借人存入 USDT**
出借人调用 `depositLend` 方法，存入 USDT 并获得 `spToken`：

```javascript
await usdtToken.approve(pledgePool.address, ethers.utils.parseEther("500")); // 存入 500 USDT
await pledgePool.depositLend(0, ethers.utils.parseEther("500")); // 借款池 ID 为 0
```

---

#### **7.4 借贷池结算**
管理员调用 `settle` 方法，结算借贷池，计算利息和本金：

```javascript
await pledgePool.settle(0); // 借款池 ID 为 0
```

---

#### **7.5 清算（如果需要）**
如果借款人的抵押品（BTC）价值不足，触发清算：

```javascript
await pledgePool.liquidate(0, borrower.address); // 借款池 ID 为 0
```

---

### **8. 扩展支持其他代币的步骤**

#### **8.1 添加代币支持**
如果需要支持新的代币（如 BTC 或 USDT），需要确保以下几点：
1. **代币合约已部署**：
   - BTC 和 USDT 需要是符合 ERC20 标准的代币合约。
2. **铸币权限**：
   - 如果代币需要铸币功能，可以通过 `AddressPrivileges` 合约添加铸币权限：
     ```javascript
     await addressPrivileges.addMinter(btcToken.address);
     await addressPrivileges.addMinter(usdtToken.address);
     ```

#### **8.2 创建新的借贷池**
调用 `createPoolInfo` 方法，为新的代币创建借贷池。例如：
- BTC 作为抵押品，借出 USDT。
- USDT 作为抵押品，借出 BTC。

#### **8.3 调整清算逻辑**
确保清算逻辑支持新的代币。例如：
- 在清算时，根据 BTC 的市场价格计算抵押品价值。
- 使用预言机（Oracle）获取 BTC 和 USDT 的价格。

---

### **9. 使用场景示例**

#### **场景 1：BTC 抵押借 USDT**
1. 借款人质押 1 BTC，借出 500 USDT。
2. 出借人存入 500 USDT，赚取利息。
3. 如果 BTC 价格下跌，触发清算，出借人获得 BTC 抵押品。

#### **场景 2：USDT 抵押借 BTC**
1. 借款人质押 1000 USDT，借出 0.5 BTC。
2. 出借人存入 0.5 BTC，赚取利息。
3. 如果 USDT 价格下跌，触发清算，出借人获得 USDT 抵押品。

---

### **10. 总结**

`PledgePool` 是一个功能完善的去中心化借贷协议，适用于多种场景，包括资金存取、借贷、清算等。通过模块化的设计和代币化的权益管理，通过 `PledgePool` 和 `AddressPrivileges` 合约的灵活设计，项目可以支持多种代币的借贷，包括 BTC 和 USDT。管理员可以通过创建新的借贷池，配置不同的借款代币和出借代币，满足用户的多样化需求。同时，通过清算机制和代币化的权益管理，确保借贷池的资金安全和高效运作。



### 项目执行

```shell
npm install
npx hardhat compile
npx hardhat test
npx hardhat ignition deploy ./ignition/modules/PledgePool.js

npx hardhat node
npx hardhat test test\BscPledgeOracle.test.js
npx hardhat test test\DebtToken.test.js
npx hardhat test test\PledgePool.test.js
```

### 项目部署

#### 部署 multiSignature 到 sepolia

````
敲黑板
!!! 需要先部署 multiSignature 到 sepolia ，然后 注意 那3管理员钱包地址，写你自己的，可以交互控制！！！
```js
// scripts/deploy/multiSignature.js
// 第五行代码开始！！！
let multiSignatureAddress = ["0x3D7155586d33a31851e28bd4Ead18A413Bc8F599",
                            "0xc3C6Ef79897Df94ddd86189A86BD9c5c7bB93Cf6",
                            "0x3B720fBacd602bccd65F82c20F8ECD5Bbb295c0a"];
let threshold = 2;
````

```shell
npx hardhat run scripts/deploy/multiSignature.js --network sepolia
```

这里，我们得到了一个多签名地址，然后 在 scripts/deploy/debtToken.js 中 使用这个地址
就叫 multiSignatureAddress

#### 部署 debtToken 到 sepolia

敲黑板
！！！这里 multiSignatureAddress 取上面部署得到的地址！！！

```js
// scripts/deploy/debtToken.js
// 第10行代码开始！！！
let multiSignatureAddress = "0xa5D1E71aC4cE6336a70E8a0cb1B6DFa87BccEf4c";
```

```shell
npx hardhat run scripts/deploy/debtToken.js --network sepolia
```

#### 部署 swapRouter 到 sepolia

```shell
npx hardhat run scripts/deploy/swapRouter.js --network sepolia
```

#### 部署 pledgePool 到 sepolia

```
npx hardhat run scripts/deploy/pledgePool.js --network sepolia


WBNB:0xd0772b878adb5c739b878e2afa060cea4a3fbc14
https://sepolia.etherscan.io/address/0xd0772b878adb5c739b878e2afa060cea4a3fbc14#code

PancakeFactory
0x5e1B1049AB259cB09e341B4f0d9426896b89fA9f

PANCAKEROUTER
0x3b75bC4e6dBAcd54023aFCB8dF0Bcd040086EabF
https://sepolia.etherscan.io/address/0x3b75bc4e6dbacd54023afcb8df0bcd040086eabf#code

multiSignature
0x1257F1804B73b8125f399A2c440763DF86FF6B50
https://sepolia.etherscan.io/address/0x1257f1804b73b8125f399a2c440763df86ff6b50#code

BscPledgeOracle
0xB574D61E7121320D708C6eC988c9CDEEc0cDDAEa
https://sepolia.etherscan.io/address/0xb574d61e7121320d708c6ec988c9cdeec0cddaea#code

PledgePool
0xbEd2F048532b859EA0272E87C07489ad7A1772DE
https://sepolia.etherscan.io/address/0xbed2f048532b859ea0272e87c07489ad7a1772de#code

DebtToken
Jpbtc
0x3b80F1c05e331eb742Db0696038F349EEFEdae5d
https://sepolia.etherscan.io/address/0x3b80f1c05e331eb742db0696038f349eefedae5d#code

Jpbusd
0x363a91fB59bEC3399a9f656A76304CDa9B34E66d
https://sepolia.etherscan.io/address/0x363a91fb59bec3399a9f656a76304cda9b34e66d



#### BSC TEST NETWORK CONTRACTS

- BUSD TOKEN : 0xE676Dcd74f44023b95E0E2C6436C97991A7497DA
- BTC TOKEN : 0xB5514a4FA9dDBb48C3DE215Bc9e52d9fCe2D8658
- DAI TOKEN : 0x490BC3FCc845d37C1686044Cd2d6589585DE9B8B
- BNB TOKEN : 0x0000000000000000000000000000000000000000

- ORACLE ADDRESS: 0x272aCa56637FDaBb2064f19d64BC3dE64A85A1b2
- SWAP ADDRESS: 0xbe9c40a0eab26a4223309ea650dea0dd4612767e
- FEE ADDRESS： 0x0ff66Eb23C511ABd86fC676CE025Ca12caB2d5d4
- PLEDGE ADDRESS: 0x216f718A983FCCb462b338FA9c60f2A89199490c
- MULTISIGNATURE: 0xcdC5A05A0A68401d5FCF7d136960CBa5aEa990Dd

```

### 本地测试

```shell
npx hardhat node

Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10000 ETH)
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10000 ETH)
Private Key: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

Account #2: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC (10000 ETH)
Private Key: 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

Account #3: 0x90F79bf6EB2c4f870365E785982E1f101E93b906 (10000 ETH)
Private Key: 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6


npx hardhat run scripts/deploy/multiSignature.js

Deploying contracts with the account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Account balance: 10000000000000000000000n
Account balance (ETH): 10000.0
multiSignature address: 0x5FbDB2315678afecb367f032d93F642f64180aa3


npx hardhat run scripts/deploy/debtToken.js

Deploying contracts with the account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Account balance: 10000000000000000000000n
DebtToken address: 0x5FbDB2315678afecb367f032d93F642f64180aa3


npx hardhat run scripts/deploy/swapRouter.js

router 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
factory 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512

npx hardhat run scripts/deploy/oracle.js

Deploying contracts with the account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Account balance: 10000000000000000000000n
Oracle address: 0x5FbDB2315678afecb367f032d93F642f64180aa3

npx hardhat run scripts/deploy/pledgePool.js

Deploying contracts with the account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Account balance: 10000000000000000000000n
pledgeAddress address: 0x5FbDB2315678afecb367f032d93F642f64180aa3
```


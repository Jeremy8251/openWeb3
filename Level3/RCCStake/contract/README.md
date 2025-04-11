```shell
# 部署指令
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/Lock.js
```

```shell
npm install @openzeppelin/contracts
npm install @openzeppelin/contracts-upgradeable
npx hardhat ignition deploy ./ignition/modules/Rcc.ts --network localhost
npx hardhat run scripts/deployRccToken.js
npx hardhat run scripts/deployRccStake.js
npx hardhat test ./test/RccStake.js

npx hardhat run scripts/deploy.js --network holesky

Deploying RccToken...
RccToken deployed to: 0x8C35E99751A7c7C2c3EF32d2eAe89cC9CC93eE02
Deploying RccStake...
RccStake deployed to: 0x7bb5395608E0029b242aB24D125cb9f4C715d8Bd
Initializing RccStake...
RccStake initialized with RccToken address: 0x8C35E99751A7c7C2c3EF32d2eAe89cC9CC93eE02
Deployment completed successfully!

# RccToken合约部署
https://eth-holesky.blockscout.com/address/0x8C35E99751A7c7C2c3EF32d2eAe89cC9CC93eE02

# RccStake
https://eth-holesky.blockscout.com/address/0x7bb5395608E0029b242aB24D125cb9f4C715d8Bd

# 初始化失败，out of gas
https://eth-holesky.blockscout.com/tx/0xf0f23b4032e10e3f71c0e3815f20ff63dbaa83e551ff8300af9ccc9f90b4d828
# 初始化成功,调整gas限制
https://eth-holesky.blockscout.com/tx/0xb4eb87cd71824015fee20808007f5d774b933f7927f37e26d0355481bc396b83
```

## `RCCStake` 智能合约

用于实现基于区块链的质押和奖励分发系统。它允许用户通过质押代币（或以太币）来赚取奖励代币 RCC，同时支持多种质押池的管理和灵活的奖励分配机制。以下是代码的详细说明：

### 合约继承与库

`RCCStake` 合约继承了多个 OpenZeppelin 提供的模块，包括 `Initializable`（支持可升级合约的初始化）、`UUPSUpgradeable`（支持 UUPS 升级模式）、`PausableUpgradeable`（支持暂停功能）和 `AccessControlUpgradeable`（提供基于角色的访问控制）。此外，合约还使用了 `SafeERC20`、`Address` 和 `Math` 库，分别用于安全的 ERC20 操作、地址类型的工具函数以及数学运算。

### 数据结构

合约定义了三个主要的数据结构：

1. **Pool**：表示质押池的信息，包括质押代币地址、池权重、奖励分发的最后区块号、累计奖励、质押代币总量、最小质押金额和解锁区块数等。
2. **UnstakeRequest**：表示用户的解押请求，包括解押金额和解锁区块号。
3. **User**：表示用户的质押信息，包括质押金额、已分发的奖励、待领取的奖励以及解押请求列表。

### 状态变量

合约包含多个状态变量：

- `startBlock` 和 `endBlock`：质押的起始和结束区块号。
- `rccPerBlock`：每个区块分发的 RCC 奖励数量。
- `withdrawPaused` 和 `claimPaused`：分别表示是否暂停解押和奖励领取功能。
- `RCC`：RCC 奖励代币的合约地址。
- `totalPoolWeight`：所有质押池的总权重。
- `pool`：质押池的数组。
- `user`：用户质押信息的映射。

### 事件与修饰符

合约定义了多个事件，用于记录重要操作（如新增质押池、更新奖励参数、用户存款等）。

修饰符如 `checkPid`、`whenNotClaimPaused` 和 `whenNotWithdrawPaused` 用于验证输入或限制功能的执行条件。

### 核心功能

1. **初始化与升级**：

   - `initialize` 函数用于初始化合约，包括设置 RCC 代币地址、起始和结束区块号、每区块奖励等。
   - `_authorizeUpgrade` 函数限制只有具有 `UPGRADE_ROLE` 的账户可以升级合约。

2. **管理员功能**：

   - 管理员可以通过 `setRCC` 设置 RCC 代币地址，通过 `pauseWithdraw` 和 `pauseClaim` 暂停解押或奖励领取功能。
   - 管理员还可以通过 `addPool` 添加新的质押池，或通过 `updatePool` 和 `setPoolWeight` 更新质押池的参数。

3. **用户功能**：

   - 用户可以通过 `depositETH` 或 `deposit` 存入以太币或代币进行质押。
   - 用户可以通过 `unstake` 发起解押请求，解押的代币会在指定的解锁区块后通过 `withdraw` 提取。
   - 用户可以通过 `claim` 领取质押奖励。

4. **奖励计算**：

   - `getMultiplier` 函数计算指定区块范围内的奖励倍数。
   - `pendingRCC` 和 `pendingRCCByBlockNumber` 函数计算用户在某个质押池中的待领取奖励。
   - `updatePool` 和 `massUpdatePools` 函数更新质押池的奖励分配状态。

5. **内部工具函数**：
   - `_deposit` 处理用户的质押逻辑，包括更新奖励状态和用户信息。
   - `_safeRCCTransfer` 和 `_safeETHTransfer` 确保安全地转移 RCC 或以太币，避免因余额不足导致的错误。

### 核心计算逻辑

在 `RCCStake.sol` 合约中，程序运行的核心计算逻辑主要围绕 **奖励分配** 和 **用户质押状态更新** 展开。以下是程序运行中的关键计算逻辑说明：

---

### **1. 奖励分配逻辑**

奖励分配的核心逻辑是根据用户的质押比例和质押池的权重，计算用户在特定时间段内应得的奖励。

#### **奖励分配公式**

```text
用户奖励 = 用户质押金额 × 池的累计奖励 - 用户已领取奖励 + 用户待领取奖励
```

#### **相关变量**

- **`accRCCPerST`**：质押池的累计奖励，每个质押代币对应的奖励数量。
- **`finishedRCC`**：用户已领取的奖励。
- **`pendingRCC`**：用户待领取的奖励。
- **`stAmount`**：用户质押的代币数量。

#### **计算逻辑**

1. **计算奖励倍数**：

   - 奖励倍数表示从上次奖励分发到当前区块之间的奖励范围。
   - 公式：
     ```solidity
     multiplier = (_to - _from) × rccPerBlock
     ```
   - `_from` 是池的 `lastRewardBlock`，`_to` 是当前区块号。

2. **计算池的奖励总量**：

   - 奖励总量根据池的权重占比计算：
     ```solidity
     totalRCC = multiplier × poolWeight / totalPoolWeight
     ```

3. **更新池的累计奖励**：

   - 如果池中有质押代币，则更新 `accRCCPerST`：
     ```solidity
     accRCCPerST += totalRCC × 1 ether / stTokenAmount
     ```
   - 这里乘以 `1 ether` 是为了提高精度，避免小数精度丢失。

4. **计算用户的奖励**：
   - 用户奖励根据用户的质押比例计算：
     ```solidity
     用户奖励 = 用户质押金额 × accRCCPerST / 1 ether - finishedRCC + pendingRCC
     ```

---

### **2. 用户质押逻辑**

用户质押时，合约需要更新用户和质押池的状态。

#### **质押逻辑**

1. **更新池的奖励状态**：

   - 调用 `updatePool` 更新池的累计奖励和最后奖励区块号。

2. **计算用户的待领取奖励**：

   - 如果用户已有质押金额，则计算其待领取奖励：
     ```solidity
     pendingRCC += 用户质押金额 × accRCCPerST / 1 ether - finishedRCC
     ```

3. **更新用户的质押金额**：

   - 将用户存入的金额累加到其总质押金额：
     ```solidity
     stAmount += _amount
     ```

4. **更新池的质押总量**：

   - 将用户存入的金额累加到池的总质押金额：
     ```solidity
     stTokenAmount += _amount
     ```

5. **更新用户的已分发奖励**：
   - 根据最新的累计奖励更新用户的 `finishedRCC`：
     ```solidity
     finishedRCC = 用户质押金额 × accRCCPerST / 1 ether
     ```

---

### **3. 解押逻辑**

用户解押时，合约会记录解押请求，并在解锁区块后允许用户提取。

#### **解押逻辑**

1. **更新池的奖励状态**：

   - 调用 `updatePool` 更新池的累计奖励和最后奖励区块号。

2. **计算用户的待领取奖励**：

   - 在解押前，计算用户的待领取奖励并更新：
     ```solidity
     pendingRCC += 用户质押金额 × accRCCPerST / 1 ether - finishedRCC
     ```

3. **记录解押请求**：

   - 将解押金额和解锁区块号记录到用户的解押请求列表：
     ```solidity
     requests.push(UnstakeRequest({
         amount: _amount,
         unlockBlocks: 当前区块号 + unstakeLockedBlocks
     }));
     ```

4. **更新用户和池的质押金额**：
   - 减少用户的质押金额和池的总质押金额：
     ```solidity
     stAmount -= _amount;
     stTokenAmount -= _amount;
     ```

---

### **4. 提取解锁的质押代币**

用户在解锁区块后可以提取解押的代币。

#### **提取逻辑**

1. **遍历解押请求**：

   - 遍历用户的解押请求列表，计算已解锁的代币总量：
     ```solidity
     if (unlockBlocks <= 当前区块号) {
         pendingWithdraw += amount;
     }
     ```

2. **移除已处理的请求**：

   - 将未解锁的请求向前移动，并移除已解锁的请求。

3. **转移解锁的代币**：
   - 如果是 ETH 池，调用 `_safeETHTransfer` 转移 ETH。
   - 如果是代币池，调用 `safeTransfer` 转移代币。

---

### **5. 奖励领取逻辑**

用户可以随时领取其待分发的奖励。

#### **领取逻辑**

1. **更新池的奖励状态**：

   - 调用 `updatePool` 更新池的累计奖励和最后奖励区块号。

2. **计算用户的待领取奖励**：

   - 根据用户的质押金额和池的累计奖励计算待领取奖励：
     ```solidity
     pendingRCC += 用户质押金额 × accRCCPerST / 1 ether - finishedRCC
     ```

3. **转移奖励**：

   - 调用 `_safeRCCTransfer` 将奖励转移给用户。

4. **更新用户的已分发奖励**：
   - 更新用户的 `finishedRCC`：
     ```solidity
     finishedRCC = 用户质押金额 × accRCCPerST / 1 ether
     ```

---

### **6. 计算示例**

假设：

- 每区块奖励 `rccPerBlock = 10 RCC`。
- 总权重 `totalPoolWeight = 100`。
- 某池权重 `poolWeight = 50`。
- 用户 A 在该池质押了 100 个代币，池的总质押量为 200 个代币。

#### **奖励分配**

1. **计算奖励倍数**：

   - 假设从区块 100 到区块 110：multiplier = (\_to - \_from) × rccPerBlock
     ```solidity
     multiplier = (110 - 100) × 10 = 100 RCC
     ```

2. **计算池的奖励总量**：totalRCC = multiplier × poolWeight / totalPoolWeight

   ```solidity
   totalRCC = 100 × 50 / 100 = 50 RCC
   ```

3. **更新池的累计奖励**：accRCCPerST += totalRCC × 1 ether / stTokenAmount

   ```solidity
   accRCCPerST += 50 × 1 ether / 200 = 0.25 ether
   ```

4. **计算用户奖励**：用户质押金额 × accRCCPerST / 1 ether

   ```solidity
   用户奖励 = 100 × 0.25 ether / 1 ether = 25 RCC
   ```

---

### **总结**

`RCCStake.sol` 的计算逻辑围绕奖励分配展开，核心是通过池的权重和用户的质押比例，动态计算奖励并更新用户和池的状态。合约通过精确的数学计算和安全的操作，确保奖励分配的公平性和准确性，同时支持灵活的质押和解押操作。

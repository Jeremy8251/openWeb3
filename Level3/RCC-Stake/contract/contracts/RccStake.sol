// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol"; //ERC20 标准接口，用于与代币合约交互
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol"; //提供安全的 ERC20 操作，避免低级调用失败
import "@openzeppelin/contracts/utils/Address.sol"; //提供与地址类型相关的工具函数
import "@openzeppelin/contracts/utils/math/Math.sol"; //提供数学运算工具

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol"; //支持合约初始化（用于可升级合约）
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol"; //支持 UUPS 升级模式
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol"; //提供基于角色的访问控制
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol"; //支持合约的暂停功能

// Initializable：支持初始化函数。
// UUPSUpgradeable：支持 UUPS 升级模式。
// PausableUpgradeable：支持暂停功能。
// AccessControlUpgradeable：支持基于角色的访问控制
contract RccStake is
    Initializable,
    UUPSUpgradeable,
    PausableUpgradeable,
    AccessControlUpgradeable
{
    //引入库函数，简化操作
    using SafeERC20 for IERC20;
    using Address for address;
    using Math for uint256;
    // 管理员角色哈希值
    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    // 升级权限角色哈希值
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    //ETH 质押池的 ID，固定为 0
    uint256 public constant ETH_PID = 0;
    //质押池的结构体，存储池的基本信息
    struct Pool {
        address stTokenAddress; // 质押代币地址
        uint256 poolWeight; // 质押池权重
        uint256 lastRewardBlock; // 上次发放奖励的区块高度
        uint256 accRCCPerST; // 每个质押代币的累计奖励
        uint256 stTokenAmount; // 质押代币总量
        uint256 minDepositAmount; // 最小质押数量
        uint256 unstakeLockedBlocks; // 解锁区块数量
    }
    //用户的解押请求
    struct UnstakeRequest {
        uint256 amount; // 解押请求代币数量
        uint256 unlockBlocks; // 解锁区块号
    }

    // 用户的质押信息
    struct User {
        uint256 stAmount; // 质押代币数量
        uint256 finishedRCC; // 已领取的奖励
        uint256 pendingRCC; // 待领取的奖励
        UnstakeRequest[] requests; // 解锁请求列表
    }
    // 状态变量
    uint256 public startBlock; // 质押开始的区块号
    uint256 public endBlock; // 质押结束的区块号
    uint256 public rccPerBlock; // 每区块分发的 RCC 奖励数量

    bool public withdrawPaused; // 是否暂停解押功能
    bool public claimPaused; // 是否暂停领取奖励功能

    IERC20 public RCC; // RCC代币合约地址

    uint256 public totalPoolWeight; // 所有质押池的总权重
    Pool[] public pool; // 质押池列表

    mapping(uint256 => mapping(address => User)) public user; // 用户质押信息

    // 事件,用于记录合约操作的关键信息
    event SetRCC(IERC20 indexed RCC);
    event PauseWithdraw();
    event UnpauseWithdraw();
    event PauseClaim();
    event UnpauseClaim();
    event SetStartBlock(uint256 indexed startBlock);
    event SetEndBlock(uint256 indexed endBlock);
    event SetRCCPerBlock(uint256 indexed rccPerBlock);
    event AddPool(
        address indexed stTokenAddress,
        uint256 indexed poolWeight,
        uint256 indexed lastRewardBlock,
        uint256 minDepositAmount,
        uint256 unstakeLockedBlocks
    );
    event UpdatePoolInfo(
        uint256 indexed poolId,
        uint256 indexed minDepositAmount,
        uint256 indexed unstakeLockedBlocks
    );
    event SetPoolWeight(
        uint256 indexed poolId,
        uint256 indexed poolWeight,
        uint256 totalPoolWeight
    );
    event UpdatePool(
        uint256 indexed poolId,
        uint256 indexed lastRewardBlock,
        uint256 totalRCC
    );
    event Deposit(address indexed user, uint256 indexed poolId, uint256 amount);
    event RequestUnstake(
        address indexed user,
        uint256 indexed poolId,
        uint256 amount
    );
    event Withdraw(
        address indexed user,
        uint256 indexed poolId,
        uint256 amount,
        uint256 indexed blockNumber
    );
    event Claim(
        address indexed user,
        uint256 indexed poolId,
        uint256 rccReward
    );

    // modifier
    // 检查pid合法
    modifier checkPid(uint256 _pid) {
        require(_pid < pool.length, "invalid pool id");
        _;
    }

    // 领取没有暂停
    modifier whenNotClaimPaused() {
        require(!claimPaused, "claim paused");
        _;
    }

    // 解押功能没有暂停
    modifier whenNotWithdrawPaused() {
        require(!withdrawPaused, "withdraw is paused");
        _;
    }

    function initialize(
        IERC20 _RCC,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _rccPerBlock
    ) public initializer {
        require(
            _startBlock <= _endBlock && _rccPerBlock > 0,
            "invalid parameters"
        );
        // 初始化访问控制模块
        __AccessControl_init();
        // 初始化UUPS可升级模块
        __UUPSUpgradeable_init();
        // 将默认管理员角色授予部署合约的账户
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        // 将管理员角色授予部署合约的账户
        _grantRole(ADMIN_ROLE, msg.sender);
        // 将升级权限授予部署合约的账户
        _grantRole(UPGRADE_ROLE, msg.sender);
        // 设置RCC代币的合约地址
        setRCC(_RCC);
        // 设置质押开始的区块号
        startBlock = _startBlock;
        // 设置质押结束的区块号
        endBlock = _endBlock;
        // 设置每个区块分发的RCC奖励数量
        rccPerBlock = _rccPerBlock;
    }

    // 升级授权,限制只有具有 `UPGRADE_ROLE` 的账户可以升级合约
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyRole(UPGRADE_ROLE) {}

    // 管理员功能
    // 设置RCC代币地址
    function setRCC(IERC20 _RCC) public onlyRole(ADMIN_ROLE) {
        RCC = _RCC;

        emit SetRCC(RCC);
    }

    // 暂停解押功能
    function pauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(!withdrawPaused, "withdraw paused");
        withdrawPaused = true;

        emit PauseWithdraw();
    }

    // 恢复解押功能
    function unpauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(withdrawPaused, "withdraw not paused");
        withdrawPaused = false;

        emit UnpauseWithdraw();
    }

    // 暂停领取奖励功能
    function pauseClaim() public onlyRole(ADMIN_ROLE) {
        require(!claimPaused, "claim paused");
        claimPaused = true;

        emit PauseClaim();
    }

    // 恢复领取奖励功能
    function unpauseClaim() public onlyRole(ADMIN_ROLE) {
        require(claimPaused, "claim not paused");
        claimPaused = false;

        emit UnpauseClaim();
    }

    // 设置质押开始的区块号
    function setStartBlock(uint256 _startBlock) public onlyRole(ADMIN_ROLE) {
        require(
            _startBlock <= endBlock,
            "start block must be smaller than end block"
        );
        startBlock = _startBlock;

        emit SetStartBlock(startBlock);
    }

    // 设置质押结束的区块号
    function setEndBlock(uint256 _endBlock) public onlyRole(ADMIN_ROLE) {
        require(
            _endBlock >= startBlock,
            "end block must be larger than start block"
        );

        endBlock = _endBlock;

        emit SetEndBlock(endBlock);
    }

    // 设置每区块分发的 RCC 奖励数量
    function setRCCPerBlock(uint256 _rccPerBlock) public onlyRole(ADMIN_ROLE) {
        require(_rccPerBlock > 0, "invalid parameter");
        rccPerBlock = _rccPerBlock;

        emit SetRCCPerBlock(rccPerBlock);
    }

    //添加新的质押池
    function addPool(
        address _stTokenAddress,
        uint256 _poolWeight,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks,
        bool _withUpdate
    ) public onlyRole(ADMIN_ROLE) {
        if (pool.length > 0) {
            require(
                _stTokenAddress != address(0x0),
                "invalid staking token address"
            );
        } else {
            //第一个池必须是 ETH 池（stTokenAddress 为 0x0）
            require(
                _stTokenAddress == address(0x0),
                "invalid staking token address"
            );
        }

        //检查解押锁定区块数是否大于 0
        require(_unstakeLockedBlocks > 0, "invalid withdraw locked blocks");
        require(block.number < endBlock, "Already ended");

        //如果 _withUpdate 为 true，更新所有池的奖励状态
        if (_withUpdate) {
            massUpdatePools();
        }

        uint256 lastRewardBlock = block.number > startBlock
            ? block.number
            : startBlock;
        totalPoolWeight = totalPoolWeight + _poolWeight;
        // 设置池的初始参数并添加到 pool 数组
        pool.push(
            Pool({
                stTokenAddress: _stTokenAddress,
                poolWeight: _poolWeight,
                lastRewardBlock: lastRewardBlock,
                accRCCPerST: 0,
                stTokenAmount: 0,
                minDepositAmount: _minDepositAmount,
                unstakeLockedBlocks: _unstakeLockedBlocks
            })
        );

        emit AddPool(
            _stTokenAddress,
            _poolWeight,
            lastRewardBlock,
            _minDepositAmount,
            _unstakeLockedBlocks
        );
    }

    //更新指定质押池的最小质押金额和解押锁定区块数
    function updatePool(
        uint256 _pid,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks
    ) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        pool[_pid].minDepositAmount = _minDepositAmount;
        pool[_pid].unstakeLockedBlocks = _unstakeLockedBlocks;
        emit UpdatePoolInfo(_pid, _minDepositAmount, _unstakeLockedBlocks);
    }

    //设置质押池权重
    function setPoolWeight(
        uint256 _pid,
        uint256 _poolWeight,
        bool _withUpdate
    ) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        // 检查权重是否大于 0
        require(_poolWeight > 0, "invalid pool weight");
        // 如果 _withUpdate 为 true，更新所有池的奖励状态
        if (_withUpdate) {
            // 更新所有池的奖励状态
            massUpdatePools();
        }
        // 重新计算总权重
        totalPoolWeight = totalPoolWeight - pool[_pid].poolWeight + _poolWeight;
        // 更新池的权重
        pool[_pid].poolWeight = _poolWeight;
        // 触发 SetPoolWeight 事件
        emit SetPoolWeight(_pid, _poolWeight, totalPoolWeight);
    }

    //获取质押池数量
    function poolLength() external view returns (uint256) {
        return pool.length;
    }

    //计算从 _from 到 _to 区块之间的奖励倍数
    function getMultiplier(
        uint256 _from,
        uint256 _to
    ) public view returns (uint256 multiplier) {
        //确保 _from 和 _to 在质押的有效区块范围内
        require(_from <= _to, "invalid block");
        if (_from >= startBlock) {
            _from = startBlock;
        }
        if (_to > endBlock) {
            _to = endBlock;
        }
        require(_from <= _to, "end block must be greater than start block");
        bool success;
        //计算奖励倍数,奖励倍数 = (_to - _from) * rccPerBlock
        (success, multiplier) = (_to - _from).tryMul(rccPerBlock);
        require(success, "multiplier overflow");
    }

    // 获取用户在指定池中的待领取奖励数量
    function pendingRCC(
        uint256 _pid,
        address _user
    ) external view checkPid(_pid) returns (uint256) {
        return pendingRCCByBlockNumber(_pid, _user, block.number);
    }

    // 按区块号计算用户待领取奖励
    function pendingRCCByBlockNumber(
        uint256 _pid,
        address _user,
        uint256 _blockNumber
    ) public view checkPid(_pid) returns (uint256) {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][_user];
        uint256 accRCCPerST = pool_.accRCCPerST;
        uint256 stSupply = pool_.stTokenAmount;
        // 如果当前区块号大于池的最后奖励区块号，则更新累计奖励
        if (_blockNumber > pool_.lastRewardBlock && stSupply != 0) {
            uint256 multiplier = getMultiplier(
                pool_.lastRewardBlock,
                _blockNumber
            );
            // 计算当前池的奖励
            uint256 rccForPool = (multiplier * pool_.poolWeight) /
                totalPoolWeight;
            // 更新累计奖励
            accRCCPerST = accRCCPerST + (rccForPool * (1 ether)) / stSupply;
        }
        // 返回用户的待领取奖励数量
        return
            (user_.stAmount * accRCCPerST) /
            (1 ether) -
            user_.finishedRCC +
            user_.pendingRCC;
    }

    // 获取用户质押的代币数量
    function stakingBalance(
        uint256 _pid,
        address _user
    ) external view checkPid(_pid) returns (uint256) {
        return user[_pid][_user].stAmount;
    }

    // 获取用户质押的代币数量，返回解押请求总量，待解押总量
    function withdrawAmount(
        uint256 _pid,
        address _user
    )
        public
        view
        checkPid(_pid)
        returns (uint256 requestAmount, uint256 pendingWithdrawAmount)
    {
        // 获取用户的质押信息
        User storage user_ = user[_pid][_user];
        // 计算用户的解押请求总量
        for (uint256 i = 0; i < user_.requests.length; i++) {
            // 如果解锁区块号小于等于当前区块号，则解押请求生效
            if (user_.requests[i].unlockBlocks <= block.number) {
                // 累加用户的待解押总量
                pendingWithdrawAmount =
                    pendingWithdrawAmount +
                    user_.requests[i].amount;
            }
            // 累加用户的解押请求总量
            requestAmount = requestAmount + user_.requests[i].amount;
        }
    }

    // 更新指定质押池的奖励数量，与最新区块同步
    function updatePool(uint256 _pid) public checkPid(_pid) {
        Pool storage pool_ = pool[_pid];
        // 如果当前区块号小于等于池的最后奖励区块号，则无需更新
        if (block.number <= pool_.lastRewardBlock) {
            return;
        }
        // 计算从 lastRewardBlock 到当前区块之间的奖励倍数,
        // 并乘以池的权重
        // 得到当前池的奖励总量
        (bool success1, uint256 totalRCC) = getMultiplier(
            pool_.lastRewardBlock,
            block.number
        ).tryMul(pool_.poolWeight);

        require(success1, "overflow");
        // 根据池的权重占比计算实际分配给池的奖励
        (success1, totalRCC) = totalRCC.tryDiv(totalPoolWeight);
        require(success1, "overflow");
        // 检查池的质押总量
        uint256 stSupply = pool_.stTokenAmount;
        if (stSupply > 0) {
            //将奖励单位转换为更高精度（避免小数精度丢失）
            (bool success2, uint256 totalRCC_) = totalRCC.tryMul(1 ether);
            require(success2, "overflow");
            // 将结果除以池的质押代币总量 stSupply，计算每个质押代币的奖励
            (success2, totalRCC_) = totalRCC_.tryDiv(stSupply);
            require(success2, "overflow");
            // 将计算出的奖励累加到池的 accRCCPerST 中
            (bool success3, uint256 accRCCPerST) = pool_.accRCCPerST.tryAdd(
                totalRCC_
            );
            require(success3, "overflow");
            // 更新池的 accRCCPerST
            pool_.accRCCPerST = accRCCPerST;
        }
        //更新池的 lastRewardBlock
        pool_.lastRewardBlock = block.number;
        //触发 UpdatePool 事件，记录池的奖励更新信息
        emit UpdatePool(_pid, pool_.lastRewardBlock, totalRCC);
    }

    // 更新所有质押池的奖励
    function massUpdatePools() public {
        uint256 length = pool.length;
        for (uint256 pid = 0; pid < length; pid++) {
            updatePool(pid);
        }
    }

    // 允许用户存入 ETH 进行质押
    function depositETH() public payable whenNotPaused {
        // 检查池是否为 ETH 池
        Pool storage pool_ = pool[ETH_PID];
        require(
            pool_.stTokenAddress == address(0x0),
            "invalid staking token address"
        );

        //检查存入金额是否大于等于池的最小质押金额
        uint256 _amount = msg.value;
        require(
            _amount >= pool_.minDepositAmount,
            "deposit amount is too small"
        );
        //调用内部函数 _deposit 处理质押逻辑
        _deposit(ETH_PID, _amount);
    }

    // 存入代币进行质押
    function deposit(
        uint256 _pid,
        uint256 _amount
    ) public whenNotPaused checkPid(_pid) {
        // 检查池是否支持代币质押（_pid 不能为 0）
        require(_pid != 0, "deposit not support ETH staking");
        Pool storage pool_ = pool[_pid];
        //检查存入金额是否大于池的最小质押金额
        require(
            _amount > pool_.minDepositAmount,
            "deposit amount is too small"
        );

        if (_amount > 0) {
            //调用 safeTransferFrom 将用户的代币转移到合约
            IERC20(pool_.stTokenAddress).safeTransferFrom(
                msg.sender,
                address(this),
                _amount
            );
        }
        //处理质押逻辑
        _deposit(_pid, _amount);
    }

    //用户发起解押请求，将质押的代币标记为待解押状态
    function unstake(
        uint256 _pid,
        uint256 _amount
    ) public whenNotPaused checkPid(_pid) whenNotWithdrawPaused {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];
        // 检查用户的质押代币数量是否足够
        require(user_.stAmount >= _amount, "Not enough staking token balance");
        // 更新池的奖励
        updatePool(_pid);
        // 计算用户的解押请求总量
        uint256 pendingRCC_ = (user_.stAmount * pool_.accRCCPerST) /
            (1 ether) -
            user_.finishedRCC;
        if (pendingRCC_ > 0) {
            // 计算用户的待领取奖励并更新
            user_.pendingRCC = user_.pendingRCC + pendingRCC_;
        }

        if (_amount > 0) {
            // 将用户的质押代币数量减去解押数量
            user_.stAmount = user_.stAmount - _amount;
            // 将解押请求添加到用户的请求列表，设置解锁区块号
            user_.requests.push(
                UnstakeRequest({
                    amount: _amount,
                    // 解锁区块号 = 当前区块号 + 解锁区块数量
                    unlockBlocks: block.number + pool_.unstakeLockedBlocks
                })
            );
        }
        // 将解押数量从池的质押代币总量中减去
        pool_.stTokenAmount = pool_.stTokenAmount - _amount;
        //  计算用户的已领取奖励
        user_.finishedRCC = (user_.stAmount * pool_.accRCCPerST) / (1 ether);
        // 触发 RequestUnstake 事件，记录用户的解押请求信息
        emit RequestUnstake(msg.sender, _pid, _amount);
    }

    // 提取解锁的质押代币
    function withdraw(
        uint256 _pid
    ) public whenNotPaused checkPid(_pid) whenNotWithdrawPaused {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];
        uint256 pendingWithdraw_; // 记录用户当前可提取的解锁代币总量
        uint256 popNum_; //记录需要从解押请求列表中移除的请求数量
        // 遍历用户的解押请求
        for (uint256 i = 0; i < user_.requests.length; i++) {
            // 如果当前请求的 unlockBlocks 大于当前区块号（block.number），说明该请求尚未解锁，停止遍历
            if (user_.requests[i].unlockBlocks > block.number) {
                break;
            }
            //如果请求已解锁，将其金额累加到 pendingWithdraw_
            pendingWithdraw_ = pendingWithdraw_ + user_.requests[i].amount;
            // 增加 popNum_，记录需要移除的解锁请求数量
            popNum_++;
        }

        // 移除已处理的解押请求
        for (uint256 i = 0; i < user_.requests.length - popNum_; i++) {
            // 将未解锁的请求向前移动，覆盖已解锁的请求
            user_.requests[i] = user_.requests[i + popNum_];
        }

        for (uint256 i = 0; i < popNum_; i++) {
            // 调用 pop() 方法，从数组末尾移除多余的请求
            user_.requests.pop();
        }

        // 将解锁的代币转移给用户（支持 ETH 和代币）
        if (pendingWithdraw_ > 0) {
            if (pool_.stTokenAddress == address(0x0)) {
                // 如果池的质押代币地址为 0x0，说明是 ETH 池，调用 _safeETHTransfer 转移 ETH。
                _safeETHTransfer(msg.sender, pendingWithdraw_);
            } else {
                // 如果是代币池，调用 safeTransfer 转移代币
                IERC20(pool_.stTokenAddress).safeTransfer(
                    msg.sender,
                    pendingWithdraw_
                );
            }
        }

        emit Withdraw(msg.sender, _pid, pendingWithdraw_, block.number);
    }

    // 领取质押奖励
    function claim(
        uint256 _pid
    ) public whenNotPaused checkPid(_pid) whenNotClaimPaused {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];
        // 更新池的奖励状态
        updatePool(_pid);
        // 计算用户的待领取奖励 = 质押数量 * 每个质押代币的累计奖励 - 已领取奖励 + 待领取奖励
        uint256 pendingRCC_ = (user_.stAmount * pool_.accRCCPerST) /
            (1 ether) -
            user_.finishedRCC +
            user_.pendingRCC;

        if (pendingRCC_ > 0) {
            // 将用户的待领取奖励清零
            user_.pendingRCC = 0;
            // 转移给用户
            _safeRCCTransfer(msg.sender, pendingRCC_);
        }
        // 更新用户的已领取奖励
        user_.finishedRCC = (user_.stAmount * pool_.accRCCPerST) / (1 ether);

        emit Claim(msg.sender, _pid, pendingRCC_);
    }

    //处理用户的质押逻辑，更新用户和池的状态
    function _deposit(uint256 _pid, uint256 _amount) internal {
        // 获取质押池和用户信息
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];
        // 更新池的奖励状态
        updatePool(_pid);

        if (user_.stAmount > 0) {
            // uint256 accST = user_.stAmount.mulDiv(pool_.accRCCPerST, 1 ether);
            // 计算用户的待领取奖励，质押金额乘以池的累计奖励
            (bool success1, uint256 accST) = user_.stAmount.tryMul(
                pool_.accRCCPerST
            );
            require(success1, "user stAmount mul accRCCPerST overflow");
            // 除以 1 ether，调整奖励的精度
            (success1, accST) = accST.tryDiv(1 ether);
            require(success1, "accST div 1 ether overflow");

            //计算用户的待领取奖励
            // 从累计奖励中减去用户已领取的奖励 finishedRCC
            (bool success2, uint256 pendingRCC_) = accST.trySub(
                user_.finishedRCC
            );
            require(success2, "accST sub finishedRCC overflow");

            if (pendingRCC_ > 0) {
                // 累加到用户的 pendingRCC
                (bool success3, uint256 _pendingRCC) = user_.pendingRCC.tryAdd(
                    pendingRCC_
                );
                require(success3, "user pendingRCC overflow");
                // 更新用户的待领取奖励
                user_.pendingRCC = _pendingRCC;
            }
        }

        if (_amount > 0) {
            // 将用户的质押代币数量累加到用户的 stAmount
            (bool success4, uint256 stAmount) = user_.stAmount.tryAdd(_amount);
            require(success4, "user stAmount overflow");
            user_.stAmount = stAmount;
        }
        // 将用户的质押代币数量累加到池的 stTokenAmount
        (bool success5, uint256 stTokenAmount) = pool_.stTokenAmount.tryAdd(
            _amount
        );
        require(success5, "pool stTokenAmount overflow");
        pool_.stTokenAmount = stTokenAmount;
        // 更新用户的已分发奖励 finishedRCC
        // user_.finishedRCC = user_.stAmount.mulDiv(pool_.accRCCPerST, 1 ether);
        // 将用户的质押金额 stAmount 乘以池的累计奖励 accRCCPerST
        (bool success6, uint256 finishedRCC) = user_.stAmount.tryMul(
            pool_.accRCCPerST
        );
        require(success6, "user stAmount mul accRCCPerST overflow");

        // 除以 1 ether，调整奖励的精度
        (success6, finishedRCC) = finishedRCC.tryDiv(1 ether);
        require(success6, "finishedRCC div 1 ether overflow");
        // 更新用户的 finishedRCC
        user_.finishedRCC = finishedRCC;
        // 触发 Deposit 事件，记录用户的存款操作
        emit Deposit(msg.sender, _pid, _amount);
    }

    // 安全地将 RCC 代币从合约转移到指定地址，避免因余额不足导致的错误
    //_to：接收 RCC 代币的目标地址。
    // _amount：需要转移的 RCC 代币数量。
    function _safeRCCTransfer(address _to, uint256 _amount) internal {
        // 查询当前合约地址中持有的 RCC 代币余额
        uint256 RCCBal = RCC.balanceOf(address(this));

        // 如果 _amount 大于合约中的 RCC 余额 RCCBal：
        if (_amount > RCCBal) {
            // 只能转移合约中现有的全部 RCC 余额 RCCBal
            RCC.transfer(_to, RCCBal);
        } else {
            // 如果 _amount 小于或等于 RCCBal,按照请求的数量 _amount 转移 RCC
            RCC.transfer(_to, _amount);
        }
    }

    // 安全地将 ETH 从合约转移到指定地址，确保转账操作成功
    function _safeETHTransfer(address _to, uint256 _amount) internal {
        // 向目标地址 _to 发送 _amount 数量的 ETH
        (bool success, bytes memory data) = address(_to).call{value: _amount}(
            ""
        );

        require(success, "ETH transfer call failed");
        if (data.length > 0) {
            // 如果目标地址返回了数据，则对数据进行解码并验证其有效性
            // 使用 abi.decode(data, (bool)) 将返回的字节数据解码为布尔值
            // 如果解码后的布尔值为 false，则抛出异常并回滚交易，错误信息
            // 如果目标地址是普通账户（非合约），通常不会返回数据，因此无需解码
            require(
                abi.decode(data, (bool)),
                "ETH transfer operation did not succeed"
            );
        }
    }
}

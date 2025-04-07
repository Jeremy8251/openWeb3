// SPDX-License-Identifier: MIT
pragma solidity 0.6.12;

import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../library/SafeTransfer.sol";
import "../interface/IDebtToken.sol";
import "../interface/IBscPledgeOracle.sol";
import "../interface/IUniswapV2Router02.sol";
import "../multiSignature/multiSignatureClient.sol";

contract PledgePool is ReentrancyGuard, SafeTransfer, multiSignatureClient {
    using SafeMath for uint256; //SafeMath 库，并将其应用于 uint256 类型
    using SafeERC20 for IERC20; //SafeERC20 库，并将其应用于 IERC20 接口
    // default decimal 最小单位
    uint256 internal constant calDecimal = 1e18;
    // Based on the decimal of the commission and interest
    uint256 internal constant baseDecimal = 1e8; //与佣金或利息相关的精度
    uint256 public minAmount = 100e18; //最小金额为 100 个以太币
    // one years 一年
    uint256 constant baseYear = 365 days;
    // 池的状态
    enum PoolState {
        MATCH, //匹配状态
        EXECUTION, //执行状态
        FINISH, //完成状态
        LIQUIDATION, //清算状态
        UNDONE //未完成状态
    }
    // 默认状态
    PoolState constant defaultChoice = PoolState.MATCH;
    // 全局暂停状态
    bool public globalPaused = false;
    // pancake swap router
    address public swapRouter; //PancakeSwap 路由器的地址
    // receiving fee address
    address payable public feeAddress; //接收费用的地址
    // oracle address
    IBscPledgeOracle public oracle; //预言机（Oracle）的地址
    // fee
    uint256 public lendFee; //借出操作的费用，费用从借出金额中扣除，影响的是出借人实际借出的资金
    uint256 public borrowFee; //借入(款)操作的费用，费用从借款金额中扣除，影响的是借款人实际收到的资金

    // 每个池的基本信息
    struct PoolBaseInfo {
        uint256 settleTime; // 结算时间
        uint256 endTime; // 结束时间
        uint256 interestRate; // 池的固定利率，单位是1e8 (1e8)
        uint256 maxSupply; // 池的最大限额
        uint256 lendSupply; // 当前实际存款的借款
        uint256 borrowSupply; // 当前实际存款的借款
        uint256 martgageRate; // 池的抵押率，单位是1e8 (1e8)
        address lendToken; // 借款方代币地址 (比如 BUSD..)
        address borrowToken; // 借款方代币地址 (比如 BTC..)
        PoolState state; // 状态 'MATCH, EXECUTION, FINISH, LIQUIDATION, UNDONE'
        IDebtToken spCoin; // sp_token的erc20地址 (比如 spBUSD_1..)
        IDebtToken jpCoin; // jp_token的erc20地址 (比如 jpBTC_1..)
        uint256 autoLiquidateThreshold; // 自动清算阈值 (触发清算阈值)
    }
    // total base pool.存储多个资金池的基本信息
    PoolBaseInfo[] public poolBaseInfo;

    // 每个池的数据信息
    struct PoolDataInfo {
        uint256 settleAmountLend; // 结算时的实际出借金额
        uint256 settleAmountBorrow; // 结算时的实际借款金额
        uint256 finishAmountLend; // 完成时的实际出借金额
        uint256 finishAmountBorrow; // 完成时的实际借款金额
        uint256 liquidationAmounLend; // 清算时的实际出借金额
        uint256 liquidationAmounBorrow; // 清算时的实际借款金额
    }
    // total data pool 存储多个资金池的动态数据信息
    PoolDataInfo[] public poolDataInfo;

    // 借款用户信息
    struct BorrowInfo {
        uint256 stakeAmount; // 当前借款的质押金额
        uint256 refundAmount; // 多余的退款金额
        bool hasNoRefund; // 默认为false，false = 未退款，true = 已退款
        bool hasNoClaim; // 默认为false，false = 未认领，true = 已认领
    }
    // Info of each user that stakes tokens.  {user.address : {pool.index : user.borrowInfo}}
    //跟踪每个用户在不同资金池中的借款状态
    mapping(address => mapping(uint256 => BorrowInfo)) public userBorrowInfo;

    // 借款用户信息
    struct LendInfo {
        uint256 stakeAmount; // 当前借款的质押金额
        uint256 refundAmount; // 超额退款金额
        bool hasNoRefund; // 默认为false，false = 无退款，true = 已退款
        bool hasNoClaim; // 默认为false，false = 无索赔，true = 已索赔
    }

    // Info of each user that stakes tokens.  {user.address : {pool.index : user.lendInfo}}
    // 存储每个用户在不同资金池中的出借信息
    mapping(address => mapping(uint256 => LendInfo)) public userLendInfo;

    // 事件
    // 存款借出事件，from是借出者地址，token是借出的代币地址，amount是借出的数量，mintAmount是生成的数量
    event DepositLend(
        address indexed from,
        address indexed token,
        uint256 amount,
        uint256 mintAmount
    );
    // 借出退款事件，from是退款者地址，token是退款的代币地址，refund是退款的数量
    event RefundLend(
        address indexed from,
        address indexed token,
        uint256 refund
    );
    // 借出索赔事件，from是索赔者地址，token是索赔的代币地址，amount是索赔的数量
    event ClaimLend(
        address indexed from,
        address indexed token,
        uint256 amount
    );
    // 提取借出事件，from是提取者地址，token是提取的代币地址，amount是提取的数量，burnAmount是销毁的数量
    event WithdrawLend(
        address indexed from,
        address indexed token,
        uint256 amount,
        uint256 burnAmount
    );
    // 存款借入事件，from是借入者地址，token是借入的代币地址，amount是借入的数量，mintAmount是生成的数量
    event DepositBorrow(
        address indexed from,
        address indexed token,
        uint256 amount,
        uint256 mintAmount
    );
    // 借入退款事件，from是退款者地址，token是退款的代币地址，refund是退款的数量
    event RefundBorrow(
        address indexed from,
        address indexed token,
        uint256 refund
    );
    // 借入索赔事件，from是索赔者地址，token是索赔的代币地址，amount是索赔的数量
    event ClaimBorrow(
        address indexed from,
        address indexed token,
        uint256 amount
    );
    // 提取借入事件，from是提取者地址，token是提取的代币地址，amount是提取的数量，burnAmount是销毁的数量
    event WithdrawBorrow(
        address indexed from,
        address indexed token,
        uint256 amount,
        uint256 burnAmount
    );
    // 交换事件，fromCoin是交换前的币种地址，toCoin是交换后的币种地址，fromValue是交换前的数量，toValue是交换后的数量
    event Swap(
        address indexed fromCoin,
        address indexed toCoin,
        uint256 fromValue,
        uint256 toValue
    );
    // 紧急借入提取事件，from是提取者地址，token是提取的代币地址，amount是提取的数量
    event EmergencyBorrowWithdrawal(
        address indexed from,
        address indexed token,
        uint256 amount
    );
    // 紧急借出提取事件，from是提取者地址，token是提取的代币地址，amount是提取的数量
    event EmergencyLendWithdrawal(
        address indexed from,
        address indexed token,
        uint256 amount
    );
    // 状态改变事件，pid是项目id，beforeState是改变前的状态，afterState是改变后的状态
    event StateChange(
        uint256 indexed pid,
        uint256 indexed beforeState,
        uint256 indexed afterState
    );
    // 设置费用事件，newLendFee是新的借出费用，newBorrowFee是新的借入费用
    event SetFee(uint256 indexed newLendFee, uint256 indexed newBorrowFee);
    // 设置交换路由器地址事件，oldSwapAddress是旧的交换地址，newSwapAddress是新的交换地址
    event SetSwapRouterAddress(
        address indexed oldSwapAddress,
        address indexed newSwapAddress
    );
    // 设置费用地址事件，oldFeeAddress是旧的费用地址，newFeeAddress是新的费用地址
    event SetFeeAddress(
        address indexed oldFeeAddress,
        address indexed newFeeAddress
    );
    // 设置最小数量事件，oldMinAmount是旧的最小数量，newMinAmount是新的最小数量
    event SetMinAmount(
        uint256 indexed oldMinAmount,
        uint256 indexed newMinAmount
    );

    constructor(
        address _oracle,
        address _swapRouter,
        address payable _feeAddress,
        address _multiSignature //多签合约的地址，用于权限管理。
    ) public multiSignatureClient(_multiSignature) {
        require(_oracle != address(0), "Is zero address");
        require(_swapRouter != address(0), "Is zero address");
        require(_feeAddress != address(0), "Is zero address");
        // 预言机（Oracle）的地址，用于获取链外数据（如价格）
        oracle = IBscPledgeOracle(_oracle);
        // PancakeSwap 路由器的地址，用于代币交换
        swapRouter = _swapRouter;
        // 接收费用的地址，标记为 payable，表示它可以接收原生代币（如 BNB）
        feeAddress = _feeAddress;
        lendFee = 0;
        borrowFee = 0;
    }

    /**
     * @dev Set the lend fee and borrow fee
     * @notice Only allow administrators to operate
     */
    //该函数允许管理员设置借贷相关的费用，包括出借费用（lendFee）和借款费用（borrowFee）。
    function setFee(uint256 _lendFee, uint256 _borrowFee) external validCall {
        lendFee = _lendFee;
        borrowFee = _borrowFee;
        emit SetFee(_lendFee, _borrowFee);
    }

    /**
     * @dev Set swap router address, example pancakeswap or babyswap..
     * @notice Only allow administrators to operate
     */
    //允许管理员更新合约中存储的交换路由器地址
    function setSwapRouterAddress(address _swapRouter) external validCall {
        require(_swapRouter != address(0), "Is zero address");
        emit SetSwapRouterAddress(swapRouter, _swapRouter);
        swapRouter = _swapRouter;
    }

    /**
     * @dev Set up the address to receive the handling fee
     * @notice Only allow administrators to operate
     */
    // 允许管理员更新合约中存储的手续费接收地址
    function setFeeAddress(address payable _feeAddress) external validCall {
        require(_feeAddress != address(0), "Is zero address");
        emit SetFeeAddress(feeAddress, _feeAddress);
        feeAddress = _feeAddress;
    }

    /**
     * @dev Set the min amount
     */
    //允许管理员动态更新合约中存储的最小金额限制
    function setMinAmount(uint256 _minAmount) external validCall {
        emit SetMinAmount(minAmount, _minAmount);
        minAmount = _minAmount;
    }

    /**
     * @dev Query pool length
     */
    // 当前资金池的总数量
    function poolLength() external view returns (uint256) {
        return poolBaseInfo.length;
    }

    /**
     * @dev 创建一个新的借贷池。函数接收一系列参数，包括结算时间、结束时间、利率、最大供应量、抵押率、借款代币、借出代币、SP代币、JP代币和自动清算阈值。
     *  Can only be called by the owner.
     */
    function createPoolInfo(
        uint256 _settleTime,
        uint256 _endTime,
        uint64 _interestRate,
        uint256 _maxSupply,
        uint256 _martgageRate,
        address _lendToken,
        address _borrowToken,
        address _spToken,
        address _jpToken,
        uint256 _autoLiquidateThreshold
    ) public validCall {
        // 检查是否已设置token ...
        // 需要结束时间大于结算时间
        require(
            _endTime > _settleTime,
            "createPool:end time grate than settle time"
        );
        // 需要_jpToken不是零地址
        require(_jpToken != address(0), "createPool:is zero address");
        // 需要_spToken不是零地址
        require(_spToken != address(0), "createPool:is zero address");

        // 推入基础池信息
        poolBaseInfo.push(
            PoolBaseInfo({
                settleTime: _settleTime,
                endTime: _endTime,
                interestRate: _interestRate,
                maxSupply: _maxSupply,
                lendSupply: 0,
                borrowSupply: 0,
                martgageRate: _martgageRate,
                lendToken: _lendToken,
                borrowToken: _borrowToken,
                state: defaultChoice,
                spCoin: IDebtToken(_spToken),
                jpCoin: IDebtToken(_jpToken),
                autoLiquidateThreshold: _autoLiquidateThreshold
            })
        );
        // 推入池数据信息
        poolDataInfo.push(
            PoolDataInfo({
                settleAmountLend: 0,
                settleAmountBorrow: 0,
                finishAmountLend: 0,
                finishAmountBorrow: 0,
                liquidationAmounLend: 0,
                liquidationAmounBorrow: 0
            })
        );
    }

    /**
     * @dev Get pool state
     * @notice returned is an int integer
     */
    // 返回指定资金池的当前状态
    function getPoolState(uint256 _pid) public view returns (uint256) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        return uint256(pool.state);
    }

    /**
     * @dev 存款人执行存款操作
     * @notice 池状态必须为MATCH
     * @param _pid 是池索引
     * @param _stakeAmount 是用户的质押金额
     */
    function depositLend(
        uint256 _pid,
        uint256 _stakeAmount
    ) external payable nonReentrant notPause timeBefore(_pid) stateMatch(_pid) {
        // 时间和状态的限制
        // 从 poolBaseInfo 数组中获取指定资金池的基本信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        // 从 userLendInfo 映射中获取当前用户在该资金池中的存款信息
        LendInfo storage lendInfo = userLendInfo[msg.sender][_pid];
        // 边界条件,检查存款金额是否超过资金池的最大限额（maxSupply）
        require(
            _stakeAmount <= (pool.maxSupply).sub(pool.lendSupply),
            "depositLend: lendSupply exceeds maxSupply"
        );
        //调用 getPayableAmount 函数计算实际存款金额
        uint256 amount = getPayableAmount(pool.lendToken, _stakeAmount);
        // 确保存款金额大于合约中设置的最小金额
        require(amount > minAmount, "depositLend: Amount less than minimum");
        // 保存借款用户信息,存款状态标记
        lendInfo.hasNoClaim = false; //未领取
        lendInfo.hasNoRefund = false; //未退款
        if (pool.lendToken == address(0)) {
            // 如果存款的代币是原生代币（如 ETH 或 BNB），则使用 msg.value 更新用户的存款金额和资金池的存款总额
            lendInfo.stakeAmount = lendInfo.stakeAmount.add(msg.value);
            pool.lendSupply = pool.lendSupply.add(msg.value);
        } else {
            //如果存款的代币是 ERC20 代币，则使用 _stakeAmount 更新用户的存款金额和资金池的存款总额
            lendInfo.stakeAmount = lendInfo.stakeAmount.add(_stakeAmount);
            pool.lendSupply = pool.lendSupply.add(_stakeAmount);
        }
        emit DepositLend(msg.sender, pool.lendToken, _stakeAmount, amount);
    }

    /**
     * @dev 退还过量存款给存款人
     * @notice 池状态不等于匹配和未完成
     * @param _pid 是池索引
     */
    function refundLend(
        uint256 _pid
    ) external nonReentrant notPause timeAfter(_pid) stateNotMatchUndone(_pid) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池的基本信息
        PoolDataInfo storage data = poolDataInfo[_pid]; // 获取池的数据信息
        LendInfo storage lendInfo = userLendInfo[msg.sender][_pid]; // 获取用户的出借信息
        // 限制金额
        require(lendInfo.stakeAmount > 0, "refundLend: not pledged"); // 需要用户已经质押了一定数量
        require(
            pool.lendSupply.sub(data.settleAmountLend) > 0,
            "refundLend: not refund"
        ); // 需要池中还有未退还的金额
        require(!lendInfo.hasNoRefund, "refundLend: repeat refund"); // 需要用户没有重复退款
        // 用户份额 = 当前质押金额 / 总金额
        uint256 userShare = lendInfo.stakeAmount.mul(calDecimal).div(
            pool.lendSupply
        );
        // refundAmount = 总退款金额 * 用户份额
        uint256 refundAmount = (pool.lendSupply.sub(data.settleAmountLend))
            .mul(userShare)
            .div(calDecimal);
        // 退款操作
        _redeem(payable(msg.sender), pool.lendToken, refundAmount);
        // 更新用户信息
        lendInfo.hasNoRefund = true;
        lendInfo.refundAmount = lendInfo.refundAmount.add(refundAmount);
        emit RefundLend(msg.sender, pool.lendToken, refundAmount); // 触发退款事件
    }

    /**
     * @dev 存款人接收 sp_toke,主要功能是让存款人领取 sp_token
     * @notice 池状态不等于匹配和未完成
     * @param _pid 是池索引
     */
    function claimLend(
        uint256 _pid
    ) external nonReentrant notPause timeAfter(_pid) stateNotMatchUndone(_pid) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池的基本信息
        PoolDataInfo storage data = poolDataInfo[_pid]; // 获取池的数据信息
        LendInfo storage lendInfo = userLendInfo[msg.sender][_pid]; // 获取用户的借款信息
        // 金额限制
        require(lendInfo.stakeAmount > 0, "claimLend: cannot claim sp_token"); // 需要用户的质押金额大于0
        require(!lendInfo.hasNoClaim, "claimLend: cannot claim again"); // 用户不能再次领取
        // 用户份额 = 当前质押金额 / 总金额
        uint256 userShare = lendInfo.stakeAmount.mul(calDecimal).div(
            pool.lendSupply
        );
        // totalSpAmount = settleAmountLend
        uint256 totalSpAmount = data.settleAmountLend; // 总的Sp金额等于借款结算金额
        // 用户 sp 金额 = totalSpAmount * 用户份额
        uint256 spAmount = totalSpAmount.mul(userShare).div(calDecimal);
        // 铸造 sp token
        pool.spCoin.mint(msg.sender, spAmount);
        // 更新领取标志
        lendInfo.hasNoClaim = true;
        emit ClaimLend(msg.sender, pool.borrowToken, spAmount); // 触发领取借款事件
    }

    /**
     * @dev 存款人取回本金和利息
     * @notice 池的状态可能是完成或清算
     * @param _pid 是池索引
     * @param _spAmount 是销毁的sp数量
     */
    function withdrawLend(
        uint256 _pid,
        uint256 _spAmount
    ) external nonReentrant notPause stateFinishLiquidation(_pid) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        PoolDataInfo storage data = poolDataInfo[_pid];
        require(_spAmount > 0, "withdrawLend: Withdraw amount is zero");
        // 销毁 sp_token
        pool.spCoin.burn(msg.sender, _spAmount);
        // 计算销毁份额
        uint256 totalSpAmount = data.settleAmountLend;
        // sp份额 = _spAmount/totalSpAmount
        uint256 spShare = _spAmount.mul(calDecimal).div(totalSpAmount);
        // 完成
        if (pool.state == PoolState.FINISH) {
            require(
                block.timestamp > pool.endTime,
                "withdrawLend: less than end time"
            );
            // 赎回金额 = finishAmountLend * sp份额
            uint256 redeemAmount = data.finishAmountLend.mul(spShare).div(
                calDecimal
            );
            // 退款动作
            _redeem(payable(msg.sender), pool.lendToken, redeemAmount);
            emit WithdrawLend(
                msg.sender,
                pool.lendToken,
                redeemAmount,
                _spAmount
            );
        }
        // 清算
        if (pool.state == PoolState.LIQUIDATION) {
            require(
                block.timestamp > pool.settleTime,
                "withdrawLend: less than match time"
            );
            // 赎回金额
            uint256 redeemAmount = data.liquidationAmounLend.mul(spShare).div(
                calDecimal
            );
            // 退款动作
            _redeem(payable(msg.sender), pool.lendToken, redeemAmount);
            emit WithdrawLend(
                msg.sender,
                pool.lendToken,
                redeemAmount,
                _spAmount
            );
        }
    }

    /**
     * @dev 紧急提取贷款
     * @notice 池状态必须是未完成
     * @param _pid 是池索引
     */
    function emergencyLendWithdrawal(
        uint256 _pid
    ) external nonReentrant notPause stateUndone(_pid) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池的基本信息
        require(pool.lendSupply > 0, "emergencLend: not withdrawal"); // 要求贷款供应大于0
        // 贷款紧急提款
        LendInfo storage lendInfo = userLendInfo[msg.sender][_pid]; // 获取用户的贷款信息
        // 限制金额
        require(lendInfo.stakeAmount > 0, "refundLend: not pledged"); // 要求质押金额大于0
        require(!lendInfo.hasNoRefund, "refundLend: again refund"); // 要求没有退款
        // 退款操作
        _redeem(payable(msg.sender), pool.lendToken, lendInfo.stakeAmount); // 执行赎回操作
        // 更新用户信息
        lendInfo.hasNoRefund = true; // 设置没有退款为真
        emit EmergencyLendWithdrawal(
            msg.sender,
            pool.lendToken,
            lendInfo.stakeAmount
        ); // 触发紧急贷款提款事件
    }

    /**
     * @dev 借款人质押操作
     * @param _pid 是池子索引
     * @param _stakeAmount 是用户质押的数量
     */
    function depositBorrow(
        uint256 _pid,
        uint256 _stakeAmount
    ) external payable nonReentrant notPause timeBefore(_pid) stateMatch(_pid) {
        // 基础信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池子基础信息
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][_pid]; // 获取用户借款信息
        // 动作
        uint256 amount = getPayableAmount(pool.borrowToken, _stakeAmount); // 获取应付金额
        require(amount > 0, "depositBorrow: deposit amount is zero"); // 要求质押金额大于0
        // 保存用户信息
        borrowInfo.hasNoClaim = false; // 设置用户未提取质押物
        borrowInfo.hasNoRefund = false; // 设置用户未退款
        // 更新信息
        if (pool.borrowToken == address(0)) {
            // 如果借款代币是0地址（即ETH）
            borrowInfo.stakeAmount = borrowInfo.stakeAmount.add(msg.value); // 更新用户质押金额
            pool.borrowSupply = pool.borrowSupply.add(msg.value); // 更新池子借款供应量
        } else {
            // 如果借款代币不是0地址（即其他ERC20代币）
            borrowInfo.stakeAmount = borrowInfo.stakeAmount.add(_stakeAmount); // 更新用户质押金额
            pool.borrowSupply = pool.borrowSupply.add(_stakeAmount); // 更新池子借款供应量
        }
        emit DepositBorrow(msg.sender, pool.borrowToken, _stakeAmount, amount); // 触发质押借款事件
    }

    /**
     * @dev 退还给借款人的过量存款，当借款人的质押量大于0，且借款供应量减去结算借款量大于0，且借款人没有退款时，计算退款金额并进行退款。
     * @notice 池状态不等于匹配和未完成
     * @param _pid 是池状态
     */
    function refundBorrow(
        uint256 _pid
    ) external nonReentrant notPause timeAfter(_pid) stateNotMatchUndone(_pid) {
        // 基础信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池的基础信息
        PoolDataInfo storage data = poolDataInfo[_pid]; // 获取池的数据信息
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][_pid]; // 获取借款人的信息
        // 条件
        require(
            pool.borrowSupply.sub(data.settleAmountBorrow) > 0,
            "refundBorrow: not refund"
        ); // 需要借款供应量减去结算借款量大于0
        require(borrowInfo.stakeAmount > 0, "refundBorrow: not pledged"); // 需要借款人的质押量大于0
        require(!borrowInfo.hasNoRefund, "refundBorrow: again refund"); // 需要借款人没有退款
        // 计算用户份额
        uint256 userShare = borrowInfo.stakeAmount.mul(calDecimal).div(
            pool.borrowSupply
        ); // 用户份额等于借款人的质押量乘以计算小数点后的位数，然后除以借款供应量
        uint256 refundAmount = (pool.borrowSupply.sub(data.settleAmountBorrow))
            .mul(userShare)
            .div(calDecimal); // 退款金额等于（借款供应量减去结算借款量）乘以用户份额，然后除以计算小数点后的位数
        // 动作
        _redeem(payable(msg.sender), pool.borrowToken, refundAmount); // 赎回
        // 更新用户信息
        borrowInfo.refundAmount = borrowInfo.refundAmount.add(refundAmount); // 更新借款人的退款金额
        borrowInfo.hasNoRefund = true; // 设置借款人已经退款
        emit RefundBorrow(msg.sender, pool.borrowToken, refundAmount); // 触发退款事件
    }

    /**
     * @dev 借款人接收 sp_token 和贷款资金
     * @notice 池状态不等于匹配和未完成
     * @param _pid 是池状态
     */
    function claimBorrow(
        uint256 _pid
    ) external nonReentrant notPause timeAfter(_pid) stateNotMatchUndone(_pid) {
        // 池基本信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        PoolDataInfo storage data = poolDataInfo[_pid];
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][_pid];
        // 限制
        require(
            borrowInfo.stakeAmount > 0,
            "claimBorrow: no jp_token to claim"
        );
        require(!borrowInfo.hasNoClaim, "claimBorrow: cannot claim again");
        // 总 jp 数量 = settleAmountLend * martgageRate
        uint256 totalJpAmount = data
            .settleAmountLend
            .mul(pool.martgageRate)
            .div(baseDecimal);
        uint256 userShare = borrowInfo.stakeAmount.mul(calDecimal).div(
            pool.borrowSupply
        );
        uint256 jpAmount = totalJpAmount.mul(userShare).div(calDecimal);
        // 铸造 jp token
        pool.jpCoin.mint(msg.sender, jpAmount);
        // 索取贷款资金
        uint256 borrowAmount = data.settleAmountLend.mul(userShare).div(
            calDecimal
        );
        _redeem(payable(msg.sender), pool.lendToken, borrowAmount);
        // 更新用户信息
        borrowInfo.hasNoClaim = true;
        emit ClaimBorrow(msg.sender, pool.borrowToken, jpAmount);
    }

    /**
     * @dev 借款人提取剩余的保证金，这个函数首先检查提取的金额是否大于0，然后销毁相应数量的JPtoken。接着，它计算JPtoken的份额，并根据池的状态（完成或清算）进行相应的操作。如果池的状态是完成，它会检查当前时间是否大于结束时间，然后计算赎回金额并进行赎回。如果池的状态是清算，它会检查当前时间是否大于匹配时间，然后计算赎回金额并进行赎回。
     * @param _pid 是池状态
     * @param _jpAmount 是用户销毁JPtoken的数量
     */
    function withdrawBorrow(
        uint256 _pid,
        uint256 _jpAmount
    ) external nonReentrant notPause stateFinishLiquidation(_pid) {
        // 声明池基础信息和数据信息的存储
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        PoolDataInfo storage data = poolDataInfo[_pid];
        // 要求提取的金额大于0
        require(_jpAmount > 0, "withdrawBorrow: withdraw amount is zero");
        // 销毁jp token
        pool.jpCoin.burn(msg.sender, _jpAmount);
        // jp份额
        uint256 totalJpAmount = data
            .settleAmountLend
            .mul(pool.martgageRate)
            .div(baseDecimal);
        uint256 jpShare = _jpAmount.mul(calDecimal).div(totalJpAmount);
        // 完成状态
        if (pool.state == PoolState.FINISH) {
            // 要求当前时间大于结束时间
            require(
                block.timestamp > pool.endTime,
                "withdrawBorrow: less than end time"
            );
            uint256 redeemAmount = jpShare.mul(data.finishAmountBorrow).div(
                calDecimal
            );
            _redeem(payable(msg.sender), pool.borrowToken, redeemAmount);
            emit WithdrawBorrow(
                msg.sender,
                pool.borrowToken,
                _jpAmount,
                redeemAmount
            );
        }
        // 清算状态
        if (pool.state == PoolState.LIQUIDATION) {
            // 要求当前时间大于匹配时间
            require(
                block.timestamp > pool.settleTime,
                "withdrawBorrow: less than match time"
            );
            uint256 redeemAmount = jpShare.mul(data.liquidationAmounBorrow).div(
                calDecimal
            );
            _redeem(payable(msg.sender), pool.borrowToken, redeemAmount);
            emit WithdrawBorrow(
                msg.sender,
                pool.borrowToken,
                _jpAmount,
                redeemAmount
            );
        }
    }

    /**
     * @dev 紧急借款提取
     * @notice 在极端情况下，总存款为0，或者总保证金为0，在某些极端情况下，如总存款为0或总保证金为0时，借款者可以进行紧急提取。首先，代码会获取池子的基本信息和借款者的借款信息，然后检查借款供应和借款者的质押金额是否大于0，以及借款者是否已经进行过退款。如果这些条件都满足，那么就会执行赎回操作，并标记借款者已经退款。最后，触发一个紧急借款提取的事件。
     * @param _pid 是池子的索引
     */
    function emergencyBorrowWithdrawal(
        uint256 _pid
    ) external nonReentrant notPause stateUndone(_pid) {
        // 获取池子的基本信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        // 确保借款供应大于0
        require(pool.borrowSupply > 0, "emergencyBorrow: not withdrawal");
        // 获取借款者的借款信息
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][_pid];
        // 确保借款者的质押金额大于0
        require(borrowInfo.stakeAmount > 0, "refundBorrow: not pledged");
        // 确保借款者没有进行过退款
        require(!borrowInfo.hasNoRefund, "refundBorrow: again refund");
        // 执行赎回操作
        _redeem(payable(msg.sender), pool.borrowToken, borrowInfo.stakeAmount);
        // 标记借款者已经退款
        borrowInfo.hasNoRefund = true;
        // 触发紧急借款提取事件
        emit EmergencyBorrowWithdrawal(
            msg.sender,
            pool.borrowToken,
            borrowInfo.stakeAmount
        );
    }

    /**
     * @dev Can it be settle
     * @param _pid is pool index
     */
    // 判断当前时间是否已经超过指定资金池的结算时间（settleTime）
    function checkoutSettle(uint256 _pid) public view returns (bool) {
        return block.timestamp > poolBaseInfo[_pid].settleTime;
    }

    /**
     * @dev  结算
     * @param _pid 是池子的索引
     */
    function settle(uint256 _pid) public validCall {
        // 获取基础池信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        // 获取数据池信息
        PoolDataInfo storage data = poolDataInfo[_pid];
        // 需要当前时间大于池子的结算时间
        require(
            block.timestamp > poolBaseInfo[_pid].settleTime,
            "settle: less than settle time"
        );
        // 池子的状态必须是匹配状态
        require(
            pool.state == PoolState.MATCH,
            "settle: pool state must be match"
        );
        if (pool.lendSupply > 0 && pool.borrowSupply > 0) {
            // 获取标的物价格
            uint256[2] memory prices = getUnderlyingPriceView(_pid);
            // 总保证金价值 = 保证金数量 * 保证金价格
            uint256 totalValue = pool
                .borrowSupply
                .mul(prices[1].mul(calDecimal).div(prices[0]))
                .div(calDecimal);
            // 转换为稳定币价值
            uint256 actualValue = totalValue.mul(baseDecimal).div(
                pool.martgageRate
            );
            if (pool.lendSupply > actualValue) {
                // 总借款大于总借出
                data.settleAmountLend = actualValue;
                data.settleAmountBorrow = pool.borrowSupply;
            } else {
                // 总借款小于总借出
                data.settleAmountLend = pool.lendSupply;
                data.settleAmountBorrow = pool
                    .lendSupply
                    .mul(pool.martgageRate)
                    .div(prices[1].mul(baseDecimal).div(prices[0]));
            }
            // 更新池子状态
            pool.state = PoolState.EXECUTION;
            // 触发事件
            emit StateChange(
                _pid,
                uint256(PoolState.MATCH),
                uint256(PoolState.EXECUTION)
            );
        } else {
            // 极端情况，借款或借出任一为0
            pool.state = PoolState.UNDONE;
            data.settleAmountLend = pool.lendSupply;
            data.settleAmountBorrow = pool.borrowSupply;
            // 触发事件
            emit StateChange(
                _pid,
                uint256(PoolState.MATCH),
                uint256(PoolState.UNDONE)
            );
        }
    }

    /**
     * @dev Can it be finish
     * @param _pid is pool index
     */
    //检查指定资金池是否已经达到结束时间
    function checkoutFinish(uint256 _pid) public view returns (bool) {
        return block.timestamp > poolBaseInfo[_pid].endTime;
    }

    /**
     * @dev 完成一个借贷池的操作，包括计算利息、执行交换操作、赎回费用和更新池子状态等步骤。
     * @param _pid 是池子的索引
     */
    function finish(uint256 _pid) public validCall {
        // 获取基础池子信息和数据信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        PoolDataInfo storage data = poolDataInfo[_pid];

        // 验证当前时间是否大于池子的结束时间
        require(
            block.timestamp > poolBaseInfo[_pid].endTime,
            "finish: less than end time"
        );
        // 验证池子的状态是否为执行状态
        require(
            pool.state == PoolState.EXECUTION,
            "finish: pool state must be execution"
        );

        // 获取借款和贷款的token
        (address token0, address token1) = (pool.borrowToken, pool.lendToken);

        // 计算时间比率 = ((结束时间 - 结算时间) * 基础小数)/365天
        uint256 timeRatio = (
            (pool.endTime.sub(pool.settleTime)).mul(baseDecimal)
        ).div(baseYear);

        // 计算利息 = 时间比率 * 利率 * 结算贷款金额
        uint256 interest = timeRatio
            .mul(pool.interestRate.mul(data.settleAmountLend))
            .div(1e16);

        // 计算贷款金额 = 结算贷款金额 + 利息
        uint256 lendAmount = data.settleAmountLend.add(interest);

        // 计算销售金额 = 贷款金额*(1+贷款费)
        uint256 sellAmount = lendAmount.mul(lendFee.add(baseDecimal)).div(
            baseDecimal
        );

        // 执行交换操作
        (uint256 amountSell, uint256 amountIn) = _sellExactAmount(
            swapRouter,
            token0,
            token1,
            sellAmount
        );

        // 验证交换后的金额是否大于等于贷款金额
        require(amountIn >= lendAmount, "finish: Slippage is too high");

        // 如果交换后的金额大于贷款金额，计算费用并赎回
        if (amountIn > lendAmount) {
            uint256 feeAmount = amountIn.sub(lendAmount);
            // 贷款费
            _redeem(feeAddress, pool.lendToken, feeAmount);
            data.finishAmountLend = amountIn.sub(feeAmount);
        } else {
            data.finishAmountLend = amountIn;
        }

        // 计算剩余的借款金额并赎回借款费
        uint256 remianNowAmount = data.settleAmountBorrow.sub(amountSell);
        uint256 remianBorrowAmount = redeemFees(
            borrowFee,
            pool.borrowToken,
            remianNowAmount
        );
        data.finishAmountBorrow = remianBorrowAmount;

        // 更新池子状态为完成
        pool.state = PoolState.FINISH;

        // 触发状态改变事件
        emit StateChange(
            _pid,
            uint256(PoolState.EXECUTION),
            uint256(PoolState.FINISH)
        );
    }

    /**
     * @dev 检查清算条件,它首先获取了池子的基础信息和数据信息，然后计算了保证金的当前价值和清算阈值，最后比较了这两个值，如果保证金的当前价值小于清算阈值，那么就满足清算条件，函数返回true，否则返回false。
     * @param _pid 是池子的索引
     */
    function checkoutLiquidate(uint256 _pid) external view returns (bool) {
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取基础池信息
        PoolDataInfo storage data = poolDataInfo[_pid]; // 获取池数据信息
        // 保证金价格
        uint256[2] memory prices = getUnderlyingPriceView(_pid); // 获取标的价格视图
        // 保证金当前价值 = 保证金数量 * 保证金价格
        uint256 borrowValueNow = data
            .settleAmountBorrow
            .mul(prices[1].mul(calDecimal).div(prices[0]))
            .div(calDecimal);
        // 清算阈值 = settleAmountLend*(1+autoLiquidateThreshold)
        uint256 valueThreshold = data
            .settleAmountLend
            .mul(baseDecimal.add(pool.autoLiquidateThreshold))
            .div(baseDecimal);
        return borrowValueNow < valueThreshold; // 如果保证金当前价值小于清算阈值，则返回true，否则返回false
    }

    /**
     * @dev 清算
     * @param _pid 是池子的索引
     */
    function liquidate(uint256 _pid) public validCall {
        PoolDataInfo storage data = poolDataInfo[_pid]; // 获取池子的数据信息
        PoolBaseInfo storage pool = poolBaseInfo[_pid]; // 获取池子的基本信息
        require(
            block.timestamp > pool.settleTime,
            "current time is less than match time"
        ); // 需要当前时间大于结算时间
        require(
            pool.state == PoolState.EXECUTION,
            "liquidate: pool state must be execution"
        ); // 需要池子的状态是执行状态
        // sellamount
        (address token0, address token1) = (pool.borrowToken, pool.lendToken); // 获取借款和贷款的token
        // 时间比率 = ((结束时间 - 结算时间) * 基础小数)/365天
        uint256 timeRatio = (
            (pool.endTime.sub(pool.settleTime)).mul(baseDecimal)
        ).div(baseYear);
        // 利息 = 时间比率 * 利率 * 结算贷款金额
        uint256 interest = timeRatio
            .mul(pool.interestRate.mul(data.settleAmountLend))
            .div(1e16);
        // 贷款金额 = 结算贷款金额 + 利息
        uint256 lendAmount = data.settleAmountLend.add(interest);
        // sellamount = lendAmount*(1+lendFee)
        // 添加贷款费用
        uint256 sellAmount = lendAmount.mul(lendFee.add(baseDecimal)).div(
            baseDecimal
        );
        (uint256 amountSell, uint256 amountIn) = _sellExactAmount(
            swapRouter,
            token0,
            token1,
            sellAmount
        ); // 卖出准确的金额
        // 可能会有滑点，amountIn - lendAmount < 0;
        if (amountIn > lendAmount) {
            uint256 feeAmount = amountIn.sub(lendAmount); // 费用金额
            // 贷款费用
            _redeem(feeAddress, pool.lendToken, feeAmount);
            data.liquidationAmounLend = amountIn.sub(feeAmount);
        } else {
            data.liquidationAmounLend = amountIn;
        }
        // liquidationAmounBorrow  借款费用
        uint256 remianNowAmount = data.settleAmountBorrow.sub(amountSell); // 剩余的现在的金额
        uint256 remianBorrowAmount = redeemFees(
            borrowFee,
            pool.borrowToken,
            remianNowAmount
        ); // 剩余的借款金额
        data.liquidationAmounBorrow = remianBorrowAmount;
        // 更新池子状态
        pool.state = PoolState.LIQUIDATION;
        // 事件
        emit StateChange(
            _pid,
            uint256(PoolState.EXECUTION),
            uint256(PoolState.LIQUIDATION)
        );
    }

    /**
     * @dev 费用计算,计算并赎回费用。首先，它计算费用，这是通过乘以费率并除以基数来完成的。如果计算出的费用大于0，它将从费用地址赎回相应的费用。最后，它返回的是原始金额减去费用。
     */
    function redeemFees(
        uint256 feeRatio,
        address token,
        uint256 amount
    ) internal returns (uint256) {
        // 计算费用，费用 = 金额 * 费率 / 基数
        uint256 fee = amount.mul(feeRatio) / baseDecimal;
        // 如果费用大于0
        if (fee > 0) {
            // 从费用地址赎回相应的费用
            _redeem(feeAddress, token, fee);
        }
        // 返回金额减去费用
        return amount.sub(fee);
    }

    /**
     * @dev Get the swap path
     */
    // 根据传入的两个代币地址（token0 和 token1）生成一个用于代币交换的路径数组（path）
    function _getSwapPath(
        address _swapRouter,
        address token0,
        address token1
    ) internal pure returns (address[] memory path) {
        IUniswapV2Router02 IUniswap = IUniswapV2Router02(_swapRouter);
        path = new address[](2);
        path[0] = token0 == address(0) ? IUniswap.WETH() : token0;
        path[1] = token1 == address(0) ? IUniswap.WETH() : token1;
    }

    /**
     * @dev Get input based on output
     */
    // 根据指定的输出代币数量（amountOut）和代币交换路径，计算所需的输入代币数量
    function _getAmountIn(
        address _swapRouter,
        address token0,
        address token1,
        uint256 amountOut
    ) internal view returns (uint256) {
        IUniswapV2Router02 IUniswap = IUniswapV2Router02(_swapRouter);
        address[] memory path = _getSwapPath(swapRouter, token0, token1);
        uint[] memory amounts = IUniswap.getAmountsIn(amountOut, path);
        return amounts[0];
    }

    /**
     * @dev sell Exact Amount
     */
    // 计算出售代币的数量
    function _sellExactAmount(
        address _swapRouter,
        address token0,
        address token1,
        uint256 amountout
    ) internal returns (uint256, uint256) {
        uint256 amountSell = amountout > 0
            ? _getAmountIn(swapRouter, token0, token1, amountout)
            : 0;
        return (amountSell, _swap(_swapRouter, token0, token1, amountSell));
    }

    /**
     * @dev Swap
     */
    // 在去中心化交易所（DEX）中执行代币交换操作
    function _swap(
        address _swapRouter,
        address token0,
        address token1,
        uint256 amount0
    ) internal returns (uint256) {
        if (token0 != address(0)) {
            _safeApprove(token0, address(_swapRouter), type(uint256).max);
        }
        if (token1 != address(0)) {
            _safeApprove(token1, address(_swapRouter), type(uint256).max);
        }
        IUniswapV2Router02 IUniswap = IUniswapV2Router02(_swapRouter);
        address[] memory path = _getSwapPath(_swapRouter, token0, token1);
        uint256[] memory amounts;
        uint currentTime = block.timestamp;
        if (token0 == address(0)) {
            amounts = IUniswap.swapExactETHForTokens{value: amount0}(
                0,
                path,
                address(this),
                currentTime + 30
            );
        } else if (token1 == address(0)) {
            amounts = IUniswap.swapExactTokensForETH(
                amount0,
                0,
                path,
                address(this),
                currentTime + 30
            );
        } else {
            amounts = IUniswap.swapExactTokensForTokens(
                amount0,
                0,
                path,
                address(this),
                currentTime + 30
            );
        }
        emit Swap(token0, token1, amounts[0], amounts[amounts.length - 1]);
        return amounts[amounts.length - 1];
    }

    /**
     * @dev Approve
     */
    // 通过低级调用（call）的方式，安全地执行 ERC-20 代币的 approve 方法
    function _safeApprove(address token, address to, uint256 value) internal {
        (bool success, bytes memory data) = token.call(
            abi.encodeWithSelector(0x095ea7b3, to, value)
        );
        require(
            success && (data.length == 0 || abi.decode(data, (bool))),
            "!safeApprove"
        );
    }

    /**
     * @dev 获取最新的预言机价格
     */
    function getUnderlyingPriceView(
        uint256 _pid
    ) public view returns (uint256[2] memory) {
        // 从基础池中获取指定的池
        PoolBaseInfo storage pool = poolBaseInfo[_pid];
        // 创建一个新的数组来存储资产
        uint256[] memory assets = new uint256[](2);
        // 将借款和贷款的token添加到资产数组中
        assets[0] = uint256(uint160(pool.lendToken));
        assets[1] = uint256(uint160(pool.borrowToken));
        // 从预言机获取资产的价格
        uint256[] memory prices = oracle.getPrices(assets);
        // 返回价格数组
        return [prices[0], prices[1]];
    }

    /**
     * @dev set Pause
     */
    // 切换合约的全局暂停状态
    function setPause() public validCall {
        globalPaused = !globalPaused;
    }

    //确保在调用被修饰的函数时，合约未处于暂停状态
    modifier notPause() {
        require(globalPaused == false, "Stake has been suspended");
        _;
    }

    //确保调用被修饰的函数时，当前时间早于指定资金池的结算时间（settleTime）
    modifier timeBefore(uint256 _pid) {
        require(
            block.timestamp < poolBaseInfo[_pid].settleTime,
            "Less than this time"
        );
        _;
    }
    //确保调用被修饰的函数时，当前时间晚于指定资金池的结算时间（settleTime）
    modifier timeAfter(uint256 _pid) {
        require(
            block.timestamp > poolBaseInfo[_pid].settleTime,
            "Greate than this time"
        );
        _;
    }
    // 确保调用被修饰的函数时，指定资金池的状态为匹配状态
    modifier stateMatch(uint256 _pid) {
        require(
            poolBaseInfo[_pid].state == PoolState.MATCH,
            "state: Pool status is not equal to match"
        );
        _;
    }
    // 确保调用被修饰的函数时，指定资金池的状态为未完成状态
    modifier stateNotMatchUndone(uint256 _pid) {
        require(
            poolBaseInfo[_pid].state == PoolState.EXECUTION ||
                poolBaseInfo[_pid].state == PoolState.FINISH ||
                poolBaseInfo[_pid].state == PoolState.LIQUIDATION,
            "state: not match and undone"
        );
        _;
    }
    // 确保调用被修饰的函数时，指定资金池的状态为未完成状态
    modifier stateFinishLiquidation(uint256 _pid) {
        require(
            poolBaseInfo[_pid].state == PoolState.FINISH ||
                poolBaseInfo[_pid].state == PoolState.LIQUIDATION,
            "state: finish liquidation"
        );
        _;
    }
    // 确保调用被修饰的函数时，指定资金池的状态为未完成状态
    modifier stateUndone(uint256 _pid) {
        require(
            poolBaseInfo[_pid].state == PoolState.UNDONE,
            "state: state must be undone"
        );
        _;
    }
}

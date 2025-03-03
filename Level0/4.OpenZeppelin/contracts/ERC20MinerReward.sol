// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract ERC20MineReward is ERC20 {
   event LogNewAlert(string description, address indexed _from, uint256 _n);

   constructor() ERC20("MinerReward", "MToken"){
//调用父合约构造函数，创建名为MinerReward、符号为MToken的ERC20代币‌
   }

   function _reward() public {
    //block.coinbase在权益证明(PoS)链中可能指向验证者地址，需确认目标链环境是否支持该特性‌
    _mint(block.coinbase, 20);//向当前区块矿工地址block.coinbase铸造20个代币‌
    emit LogNewAlert('_rewarded', block.coinbase, block.number);//触发LogNewAlert事件记录铸造信息‌
   }
}
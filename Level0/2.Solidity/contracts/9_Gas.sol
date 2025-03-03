// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Gas {
    uint256 public i = 0;

    function forever() public {
        while (true) {
            i += 1;
        }
    }
    /**
     * 
状态	0x0 交易已打包但执行失败
交易哈希	0xca12e3ffb87095435141c862134acb19af48d7e7c8ea6420b8b8b2a2ab18e8d6
区块哈希	0xe0db26a62c3c3edd9848f2c2c40a4cefc5229add508b0eba07fcbac58555c08a
区块号	28
from	0x5B38Da6a701c568545dCfcB03FcB875f56beddC4
to	Gas.forever() 0x4a9C121080f6D9250Fc0143f41B595fD172E31bf
gas	3000000 gas
交易成本	3000000 gas 
执行成本	2978936 gas 
输入	0x9ff...9a603
输出	0x
解码输入	{}
解码输出	{}
日志	[]
原始日志	[]
transact to Gas.forever errored: Error occurred: out of gas.

out of gas
	The transaction ran out of gas. Please increase the Gas Limit.

If the transaction failed for not having enough gas, try increasing the gas limit gently.
     */
}

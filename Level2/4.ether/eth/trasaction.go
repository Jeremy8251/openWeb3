package eth

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryTransaction(client *ethclient.Client) {
	ctx := context.Background()
	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainId = ", chainId) //chainId =  17000
	blockNumber := big.NewInt(3429131)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	var signer types.Signer
	for _, tx := range block.Transactions() {
		fmt.Println("====================")
		fmt.Println(tx.Hash().Hex())
		// fmt.Println(tx.Value().String())    // 100000000000000000
		// fmt.Println(tx.Gas())               // 21000
		// fmt.Println(tx.GasPrice().Uint64()) // 100000000000
		// fmt.Println(tx.Nonce())             // 245132
		// fmt.Println(tx.Data())              // []
		// fmt.Println(tx.To().Hex())          // 0x8F9aFd209339088Ced7Bc0f57Fe08566ADda3587

		// 出现 transaction type not supported 错误的原因是代码中使用的交易类型与当前签名器（EIP155Signer）不兼容
		// if sender, err := types.Sender(types.NewEIP155Signer(chainId), tx); err == nil {
		// 将 NewEIP155Signer 替换为 LatestSignerForChainID，自动适配交易类型：
		signer = types.LatestSignerForChainID(chainId)
		sender, err := types.Sender(signer, tx)
		if err == nil {
			fmt.Println("sender", sender.Hex())
			//0x6c7086ab06e4d0459826a073aa205e209d190b2b2bcef4439b5fa6c35d291596
			//sender 0xC6392aD8A14794eA57D237D12017E7295bea2363
		} else {
			fmt.Println("sender get error")
			log.Fatal(err)
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receipt = ", receipt.Status)
		fmt.Println("receipt = ", receipt.Logs)
	}
	blockHash := common.HexToHash("0x3bb204ff34ebe3e1d29139602001f8e8a5c893926a85ef69bf7a6f7979261b20")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count = ", count)
	fmt.Println("=========22222===========")
	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		// break
	}
	fmt.Println("=========33333===========")
	//https://holesky.beaconcha.in/tx/0x6032e8fd3adb799c40dc5edf0769624bff43f0db347e1eea4f27a3e7df332c6b
	txHash := common.HexToHash("0x6032e8fd3adb799c40dc5edf0769624bff43f0db347e1eea4f27a3e7df332c6b")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(isPending)
	fmt.Println("Hash:", tx.Hash().Hex())
	fmt.Println("Time:", tx.Time())
	fmt.Println("Gaslimit:", tx.Gas())
	fmt.Println("GasPrice:", tx.GasPrice())
	fmt.Println("Cost:", tx.Cost())
	fmt.Println("ChainId:", tx.ChainId())
	fmt.Println("Value:", tx.Value())
	fmt.Println("To:", tx.To())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		fmt.Println("sender error")
	}
	fmt.Println("sender:", sender)

	etherInWei, _ := new(big.Float).SetString(tx.Cost().String()) // 1e18 as a big integer
	etherValue := big.NewFloat(math.Pow10(18))                    // 设置足够精度以避免舍入误差                                  // 将wei转换为big.Float以便进行除法操作
	valueEther := new(big.Float).Quo(etherInWei, etherValue)      // 除以1e18得到Ether值
	fmt.Println("valueEther: ", valueEther, " eth")
}

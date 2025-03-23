package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryBlock(client *ethclient.Client) {
	// address := common.HexToAddress("0xd2B1a5cc4BED90A121580015DF7F403B2F4FD8B8")
	// balance, err := client.BalanceAt(context.Background(), address, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// etherInWei, _ := new(big.Int).SetString("1000000000000000000", 10) // 1e18 as a big integer
	// etherValue := new(big.Float).SetPrec(128)                          // 设置足够精度以避免舍入误差
	// etherValue.SetInt(balance)                                         // 将wei转换为big.Float以便进行除法操作
	// etherValue.Quo(etherValue, new(big.Float).SetInt(etherInWei))      // 除以1e18得到Ether值
	// etherValueAsFloat64, _ := etherValue.Float64()                     // 转换为float64以便打印或进一步处理
	// fmt.Println("余额 = ", etherValueAsFloat64)

	header, err := client.HeaderByNumber(context.Background(), nil) //参数 nil 表示‌获取最新的区块头‌。
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("header =", header.Number.String())
	fmt.Println(header.Time)
	fmt.Println(header.Difficulty.Uint64())
	fmt.Println(header.Hash().Hex())

	blockNumber := big.NewInt(3429131)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())     // 3429131
	fmt.Println(block.Time())                // 1712798400
	fmt.Println(block.Difficulty().Uint64()) // 0
	fmt.Println(block.Hash().Hex())          // 0x3bb204ff34ebe3e1d29139602001f8e8a5c893926a85ef69bf7a6f7979261b20
	fmt.Println(len(block.Transactions()))   // 16
}

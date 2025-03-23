package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func QueryReceipt(client *ethclient.Client) {
	blockNumber := big.NewInt(3429131)
	blockHash := common.HexToHash("0x3bb204ff34ebe3e1d29139602001f8e8a5c893926a85ef69bf7a6f7979261b20")
	receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("receiptByHash = ", receiptByHash[0])
	receiptsByNum, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("receiptsByNum = ", receiptsByNum[0])
	fmt.Println(receiptByHash[0] == receiptsByNum[0])              // false
	fmt.Println(IsSameReceipt(receiptByHash[0], receiptsByNum[0])) // true
}

func IsSameReceipt(r1, r2 *types.Receipt) bool {
	//比较交易哈希、状态码等不可变字段
	return r1.TxHash == r2.TxHash &&
		r1.Status == r2.Status &&
		r1.BlockHash == r2.BlockHash
}

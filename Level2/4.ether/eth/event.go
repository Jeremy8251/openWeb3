package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// https://eth-holesky.blockscout.com/tx/0x44609f104e457480fd72760974563d22dad046e87c2babd3bc78f258147304e1
// 32 位测试
// var key = "0x0000000000000000000000000000000000000000000000000000000000000055"
// var value = "0x0000000000000000000000000000000000000000000000000000000000005555"
var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

// 查询事件
func EventWithabiByWSS(wss string) {
	fmt.Println("wss :", wss)
	client, err := ethclient.Dial(wss)
	if err != nil {
		log.Fatal(err)
	}
	// 从 go-ethereum 包中导入 FilterQuery 结构体并用过滤选项初始化
	contractAddress := common.HexToAddress(contractAddr)
	// 创建筛选查询
	query := ethereum.FilterQuery{
		// BlockHash
		FromBlock: big.NewInt(3520562),
		ToBlock:   big.NewInt(3530562),
		Addresses: []common.Address{
			contractAddress,
		},
		// Topics: [][]common.Hash{
		//  {},
		//  {},
		// },
	}
	// 调用 ethclient 的 FilterLogs，它接收我们的查询并将返回所有的匹配事件日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	// 返回的所有日志将是 ABI 编码，因此它们本身不会非常易读。为了解码日志，我们需要导入我们智能合约的 ABI
	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}
	// 通过日志进行迭代并将它们解码为我么可以使用的类型
	for _, vLog := range logs {
		fmt.Println("BlockHash:", vLog.BlockHash.Hex())
		fmt.Println("BlockNumber:", vLog.BlockNumber)
		fmt.Println("TxHash:", vLog.TxHash.Hex())

		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Key:", string(event.Key[:]))
		fmt.Println("Value:", string(event.Value[:]))
		fmt.Println(common.Bytes2Hex(event.Key[:]))
		fmt.Println(common.Bytes2Hex(event.Value[:]))
		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}
		fmt.Println("topics[0]=", topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
		if len(topics) > 1 {
			fmt.Println("indexed topics:", topics[1:])
		}

	}

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("signature topics=", hash.Hex())

}

// 订阅事件
func EventSubscribeByWSS(wss string) {
	fmt.Println("wss :", wss)
	client, err := ethclient.Dial(wss)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(contractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(StoreABI)))
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog.BlockHash.Hex())
			fmt.Println(vLog.BlockNumber)
			fmt.Println(vLog.TxHash.Hex())
			event := struct {
				Key   [32]byte
				Value [32]byte
			}{}
			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(common.Bytes2Hex(event.Key[:]))
			fmt.Println(common.Bytes2Hex(event.Value[:]))
			var topics []string
			for i := range vLog.Topics {
				topics = append(topics, vLog.Topics[i].Hex())
			}
			fmt.Println("topics[0]=", topics[0])
			if len(topics) > 1 {
				fmt.Println("index topic:", topics[1:])
			}
		}
	}
}

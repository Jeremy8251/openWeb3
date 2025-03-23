package eth

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SubscribeNewBlockByWSS(wss string) {
	fmt.Println("wss :", wss)
	client, err := ethclient.Dial(wss)
	if err != nil {
		log.Fatal(err)
	}
	//创建一个新的通道，用于接收最新的区块头
	headers := make(chan *types.Header)
	// 接收我们刚创建的区块头通道，该方法将返回一个订阅对象
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("接收中....")
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("block.Hash:", block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println("block.Number:", block.Number().Uint64()) // 3477413
			fmt.Println("block.Time:", block.Time())              // 1529525947
			fmt.Println("block.Nonce:", block.Nonce())
		}
	}
}

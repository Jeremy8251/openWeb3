package main

import (
	"ether/eth"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// 查询交易
// https://eth-holesky.blockscout.com/tx/0x44609f104e457480fd72760974563d22dad046e87c2babd3bc78f258147304e1

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// 从环境变量中读取值
	url := os.Getenv("URL")
	privateKey1 := os.Getenv("PRIVATE_KEY1")
	privateKey2 := os.Getenv("PRIVATE_KEY2")
	accountAddress1 := os.Getenv("Account_Address1")
	wss := os.Getenv("WSS")
	fmt.Println(privateKey1, privateKey2, accountAddress1, url, wss)
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("we have a connection", client)

	// eth.QueryBlock(client)
	// eth.QueryTransaction(client)
	// eth.QueryReceipt(client)
	// eth.CreateWallet(client, privateKey2)
	// eth.Transfer(client, privateKey1)
	// eth.TokenTransfer(client, privateKey1)
	// eth.QueryTokenBalance(client, accountAddress1)
	// eth.SubscribeNewBlock(wss)
	// eth.DeployByabi(client, privateKey1)
	// eth.DeployBybin(client, privateKey1)
	// eth.LoadContract(client)
	// eth.QueryContract(client)
	// eth.ExecContractByGo(client, privateKey1)
	// eth.ExecContractByabi(client, privateKey1)
	// eth.ExecContractBynone(client, privateKey1)
	// eth.EventWithabiByWSS(wss)
	// eth.EventSubscribeByWSS(wss)
	// eth.CreateKeystore()
	eth.ImportKeystore()
}

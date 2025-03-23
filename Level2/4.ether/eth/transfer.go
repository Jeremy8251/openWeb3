package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Transfer(client *ethclient.Client, pvKey string) {
	// 发送交易需要发送方私钥对该交易签名
	privateKey, err := crypto.HexToECDSA(pvKey)
	if err != nil {
		log.Fatal(err)
	}
	// 从私钥派生公共地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 读取我们应该用于帐户交易的随机数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 设置转账的 ETH 数量
	value := big.NewInt(1000000000000000) // in wei (0.001 eth)
	gasLimit := uint64(21000)             // in units，ETH 转账的燃气应设上限为“21000”单位
	// gasPrice := big.NewInt(30000000000)   // in wei (30 gwei)燃气价格必须以 wei 为单位设定
	// SuggestGasPrice 函数，用于根据'x'个先前块来获得平均燃气价格
	gasPrice, err2 := client.SuggestGasPrice(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	// 接收方
	toAddress := common.HexToAddress("0x85c43dD33cabd2E270B171A2c2fCf9Fbb8D960Dd")

	// 发送 ETH 的数据字段为“nil”。 在与智能合约进行交互时，我们将使用数据字段，仅仅转账以太币是不需要数据字段的。
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	// 发件人的私钥对交易进行签名，需要链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// client 实例调用 SendTransaction 来将已签名的事务广播到整个网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa5c69d07d18eb3d7b2d0df4f21f12f86dac909278910233137cf768c61c9289f
}

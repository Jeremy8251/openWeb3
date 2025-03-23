package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 使用生成的Go合约代码执行合约方法
func ExecContractByGo(client *ethclient.Client, pvk string) {
	// storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	storeContract, err := LoadContract(client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(pvk)
	if err != nil {
		log.Fatal(err)
	}

	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("exec_save_key222"))
	copy(value[:], []byte("exec_save_value222"))
	fmt.Println("key:", string(key[:]))
	fmt.Println("value:", string(value[:]))
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainID:", chainID)
	// 写入智能合约需要我们用私钥来对交易事务进行签名
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())

	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := storeContract.Items(callOpt, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("valueInContract:", string(valueInContract[:]))
	fmt.Println("is value saving in contract equals to origin value:", valueInContract == value)
}

// 使用abi执行合约方法
func ExecContractByabi(client *ethclient.Client, pvk string) {
	// 获取私钥
	privateKey, err := crypto.HexToECDSA(pvk)
	if err != nil {
		log.Fatal(err)
	}
	// 从私钥实例获取公开地址：
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 使用地址获取地址的 nonce 值:
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 估算 gas 价格：
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 准备交易 calldata, 使用Store_sol_Store.abi字符串
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatal(err)
	}

	methodName := "setItem"
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("exec_save_key_use_abi_v1111"))
	copy(value[:], []byte("exec_save_value_use_abi_v1111"))
	input, err := contractABI.Pack(methodName, key, value)
	if err != nil {
		log.Fatal(err)
	}
	// 创建并签名交易：
	// chainID := big.NewInt(int64(11155111))
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainID:", chainID)

	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddr), big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送签名好的交易：
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 创建 call 查询
	callInput, err := contractABI.Pack("items", key)
	if err != nil {
		log.Fatal(err)
	}
	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}
	//解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 各种数据类型编码方式具体可以查看 abi 包中的 Pack 方法，返回数据解析查看 abi 包中的 UnpackIntoInterface 方法
	var unpacked [32]byte
	contractABI.UnpackIntoInterface(&unpacked, "items", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("value:", string(value[:]))
	fmt.Println("result:", string(result[:]))
	fmt.Println("unpacked:", string(unpacked[:]))
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)
}

// 不使用abi执行合约方法
func ExecContractBynone(client *ethclient.Client, pvk string) {
	// 获取私钥
	privateKey, err := crypto.HexToECDSA(pvk)
	if err != nil {
		log.Fatal(err)
	}
	// 从私钥实例获取公开地址：
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 使用地址获取地址的 nonce 值:
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 估算 gas 价格：
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 调用合约前的数据准备,不需要abi直接组装setItem
	methodSignature := []byte("setItem(bytes32,bytes32)")
	methodSelector := crypto.Keccak256(methodSignature)[:4]

	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("exec_save_key_no_use_abi_v3333"))
	copy(value[:], []byte("exec_save_value_no_use_abi_v3333"))

	// 组合调用数据input
	var input []byte
	input = append(input, methodSelector...)
	input = append(input, key[:]...)
	input = append(input, value[:]...)

	// 创建并签名交易：
	// chainID := big.NewInt(int64(11155111))
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chainID:", chainID)

	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddr), big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送签名好的交易：
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 创建 call 查询，直接组合方法items(bytes32)
	itemsSignature := []byte("items(bytes32)")
	itemsSelector := crypto.Keccak256(itemsSignature)[:4]

	var callInput []byte
	callInput = append(callInput, itemsSelector...)
	callInput = append(callInput, key[:]...)

	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}
	//解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var unpacked [32]byte
	copy(unpacked[:], result)
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)
	fmt.Println("value:", string(value[:]))
	fmt.Println("result:", string(result[:]))
	fmt.Println("unpacked:", string(unpacked[:]))
}

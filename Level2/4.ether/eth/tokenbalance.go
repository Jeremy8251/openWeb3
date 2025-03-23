package eth

import (
	"ether/token"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 代币账户余额
func QueryTokenBalance(client *ethclient.Client, accountAddress string) {
	//代币的地址
	tokenAddress := common.HexToAddress("0x6Fe2eC23aF78cC83C7Bf37dB4193AA1c11096fE5")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	//查询用户的代币余额
	address := common.HexToAddress(accountAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("account wei: %s\n", bal) // "wei: 74605500647408739782407023"
	//ERC20 智能合约的公共变量
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	totalSupply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("name: %s\n", name)               // name: DemoToken
	fmt.Printf("symbol: %s\n", symbol)           // symbol: DTK
	fmt.Printf("decimals: %v\n", decimals)       // decimals: 18
	fmt.Printf("totalSupply: %v\n", totalSupply) // totalSupply: 1000000000000000000000

	//简单的数学运算将余额转换为可读的十进制格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value)
	// 打印
	// account wei: 1000000000000000000000
	// name: DemoToken
	// symbol: DTK
	// decimals: 18
	// totalSupply: 1000000000000000000000
	// balance: 1000.00000
}

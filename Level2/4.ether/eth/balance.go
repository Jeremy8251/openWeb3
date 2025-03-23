package eth

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 传区块高度能读取指定区块时的账户余额，区块高度必须是 big.Int 类型:big.NewInt(blockNumber),可为nil
func QueryBalance(client *ethclient.Client, accountAddress string, blockNumber *big.Int) {
	account := common.HexToAddress(accountAddress)
	// 	blockNumber := big.NewInt(blockNumber)
	balance, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25729324269165216042
	ConvertToETH(balance)
	// 待处理的账户余额是多少
	pendingBalance, err2 := client.PendingBalanceAt(context.Background(), account)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(pendingBalance) // 25729324269165216042
}

func ConvertToETH(balance *big.Int) *big.Float {
	// 在 ETH 中它是_wei_。要读取 ETH 值，您必须做计算 wei/10^18。因为我们正在处理大数，我们得导入原生的 Go math 和 math/big 包。这是您做的转换。
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("Balance ETH =", ethValue) // 25.729324269165216041
	return ethValue
}

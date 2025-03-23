package eth

import (
	"ether/store"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractAddr = "0xA51041951519e5eB8c10Ad01C06e992072DE428D"

func LoadContract(client *ethclient.Client) (*store.Store, error) {
	instance, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	return instance, err
}

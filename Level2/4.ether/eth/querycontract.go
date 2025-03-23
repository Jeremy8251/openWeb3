package eth

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryContract(client *ethclient.Client) {
	instance, err := LoadContract(client)
	if err != nil {
		log.Fatal(err)
	}
	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("version =", version) // "1.0"
}

package eth

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func CreateWallet(client *ethclient.Client, defaultKey string) {

	//生成随机私钥
	// privateKey, err := crypto.GenerateKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//如果已经有了私钥的 Hex 字符串，也可以使用 HexToECDSA 方法恢复私钥：
	privateKey, err := crypto.HexToECDSA(defaultKey) // defaultKey等于随机私钥
	if err != nil {
		log.Fatal(err)
	}

	//使用 FromECDSA 方法将其转换为字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x' b6cf3f5d293b0da2e39f3aba91c5df25764846eb57920a28d6ed612cd58b09f2
	//由于公钥是从私钥派生的，因此 go-ethereum 的加密私钥具有一个返回公钥的 Public 方法
	publicKey := privateKey.Public()
	//将其转换为十六进制的过程
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04' 9d78a3555a0ef4b21f8a9e9e43065995112dbb9a606fbb3839e416d06eee395664ddbf4d9bcf54ff17b12c70f6a7e7149a32d6e125ad48f6715c95789ada875f
	// 生成你经常看到的公共地址,它接受一个 ECDSA 公钥，并返回公共地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) //0x2158A3aa3d5df4c64D909742A42EeE3E43AC677D
	// 公共地址其实就是公钥的 Keccak-256 哈希，然后我们取最后 40 个字符（20 个字节）并用“0x”作为前缀。
	// 以下是使用 golang.org/x/crypto/sha3 的 Keccak256 函数手动完成的方法
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) //// 原长32位，截去12位，保留后20位,0x2158a3aa3d5df4c64d909742a42eee3e43ac677d
}

package eth

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateKeystore() {
	ks := keystore.NewKeyStore("./keystores", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("keystore密钥创建账号:", account.Address.Hex())
}

func ImportKeystore() {
	file := "./keystores/UTC--2025-03-21T14-53-17.384847800Z--d0691b37ad2410b399f770834b78e1ebb2c2e962"
	// 检查文件是否存在再读取
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatal("Wallet file not found")
	}

	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	password := "secret"
	// 更安全的密码获取方式（如从环境变量读取）
	// password := os.Getenv("WALLET_PASSWORD")
	key, err := keystore.DecryptKey(jsonBytes, password)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	fmt.Println("私钥: ", hexutil.Encode(privateKeyBytes)[2:])
	publicKeyBytes := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	fmt.Println("公钥: ", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'剥离了0x和前2个字符04，它始终是EC前缀，不是必需的
	address := crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()
	fmt.Println("地址: ", address)
	// 调用Import方法，该方法接收keystore的JSON数据作为字节。第二个参数是用于加密私钥的口令。第三个参数是指定一个新的加密口令，
	// 但我们在示例中使用一样的口令。导入账户将允许您按期访问该账户，但它将生成新keystore文件

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("account: ", account.Address.Hex()) // 0x4aFF6791dF5e4659891d1D567dDF585F2FA56704

	// if err := os.Remove(file); err != nil {
	// 	log.Fatal(err)
	// }

}

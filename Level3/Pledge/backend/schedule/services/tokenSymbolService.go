package services

import (
	"encoding/json"
	"errors"
	"os"
	"pledge-backend/config"
	abifile "pledge-backend/contract/abi"
	"pledge-backend/db"
	"pledge-backend/log"
	"pledge-backend/schedule/models"
	"pledge-backend/utils"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type TokenSymbol struct{}

// 从区块链上获取代币的符号（symbol），并将其存储到数据库和 Redis 中
func NewTokenSymbol() *TokenSymbol {
	return &TokenSymbol{}
}

// UpdateContractSymbol get contract symbol
// 遍历数据库中的代币信息，逐个获取代币的符号，并更新到数据库和 Redis 中
func (s *TokenSymbol) UpdateContractSymbol() {
	var tokens []models.TokenInfo
	// 连接数据库，查询所有代币信息
	db.Mysql.Table("token_info").Find(&tokens)
	for _, t := range tokens {
		if t.Token == "" {
			log.Logger.Sugar().Error("UpdateContractSymbol token empty", t.Symbol, t.ChainId)
			continue
		}
		err := errors.New("")
		symbol := ""
		if t.ChainId == config.Config.TestNet.ChainId {
			// 如果代币在测试网上，则获取测试网的符号
			err, symbol = s.GetContractSymbolOnTestNet(t.Token, config.Config.TestNet.NetUrl)
		} else if t.ChainId == config.Config.MainNet.ChainId {
			// 如果代币在主网上，则获取主网的符号
			if t.AbiFileExist == 0 {
				// 如果主网的 ABI 文件本地不存在，则从远程获取并保存到本地
				err = s.GetRemoteAbiFileByToken(t.Token, t.ChainId)
				if err != nil {
					log.Logger.Sugar().Error("UpdateContractSymbol GetRemoteAbiFileByToken err ", t.Symbol, t.ChainId, err)
					continue
				}
			}
			err, symbol = s.GetContractSymbolOnMainNet(t.Token, config.Config.MainNet.NetUrl)
		} else {
			log.Logger.Sugar().Error("UpdateContractSymbol chain_id err ", t.Symbol, t.ChainId)
			continue
		}
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractSymbol err ", t.Symbol, t.ChainId, err)
			continue
		}
		// 检查 Redis 中的符号数据是否需要更新
		hasNewData, err := s.CheckSymbolData(t.Token, t.ChainId, symbol)
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractSymbol CheckSymbolData err ", err)
			continue
		}

		if hasNewData {
			// 如果需要更新，则保存到数据库
			err = s.SaveSymbolData(t.Token, t.ChainId, symbol)
			if err != nil {
				log.Logger.Sugar().Error("UpdateContractSymbol SaveSymbolData err ", err)
				continue
			}
		}
	}
}

// GetRemoteAbiFileByToken get and save remote abi file on main net
// 从远程下载代币的 ABI 文件，并存储到本地
func (s *TokenSymbol) GetRemoteAbiFileByToken(token, chainId string) error {
	// 构造远程请求 URL
	// url := "https://api.bscscan.com/api?module=contract&action=getabi&apikey=HJ3WS4N88QJ6S7PQ8D89BD49IZIFP1JFER&address=" + token

	url := "https://api-sepolia.etherscan.io/api?module=contract&action=getabi&address=" + token
	// 发送 HTTP 请求获取 ABI 文件
	res, err := utils.HttpGet(url, map[string]string{})
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	// 解析和格式化 ABI 数据
	resStr := s.FormatAbiJsonStr(string(res))

	abiJson := models.AbiJson{}
	err = json.Unmarshal([]byte(resStr), &abiJson)
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}

	if abiJson.Status != "1" {
		log.Logger.Sugar().Error("get remote abi file failed: status 0 ", resStr)
		return errors.New("get remote abi file failed: status 0 ")
	}

	// marshal and format
	abiJsonBytes, err := json.MarshalIndent(abiJson.Result, "", "\t")
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	// 存储到本地文件
	newAbiFile := abifile.GetCurrentAbPathByCaller() + "/" + token + ".abi"

	err = os.WriteFile(newAbiFile, abiJsonBytes, 0777)
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	// 更新数据库中的 ABI
	err = db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).Updates(map[string]interface{}{
		"abi_file_exist": 1,
	}).Debug().Error
	if err != nil {
		return err
	}
	return nil
}

// FormatAbiJsonStr format the abi string
// 格式化 ABI 字符串，去除多余的转义字符
func (s *TokenSymbol) FormatAbiJsonStr(result string) string {
	resStr := strings.Replace(result, `\`, ``, -1)
	resStr = strings.Replace(result, `\"`, `"`, -1)
	resStr = strings.Replace(resStr, `"[{`, `[{`, -1)
	resStr = strings.Replace(resStr, `}]"`, `}]`, -1)
	return resStr
}

// GetContractSymbolOnMainNet get contract symbol on main net
func (s *TokenSymbol) GetContractSymbolOnMainNet(token, network string) (error, string) {
	ethereumConn, err := ethclient.Dial(network)
	if nil != err {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	abiStr, err := abifile.GetAbiByToken(token)
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	contract, err := bind.NewBoundContract(common.HexToAddress(token), parsed, ethereumConn, ethereumConn, ethereumConn), nil
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}

	res := make([]interface{}, 0)
	err = contract.Call(nil, &res, "symbol")
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", err)
		return err, ""
	}

	return nil, res[0].(string)
}

// GetContractSymbolOnTestNet get contract symbol on test net
// 通过以太坊客户端与智能合约交互，调用 symbol 方法获取代币符
func (s *TokenSymbol) GetContractSymbolOnTestNet(token, network string) (error, string) {
	// 连接区块链网络
	ethereumConn, err := ethclient.Dial(network)
	if nil != err {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	// 加载 ABI 文件
	abiStr, err := abifile.GetAbiByToken("erc20")
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	// 解析 ABI 文件
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}
	// 创建合约实例
	contract, err := bind.NewBoundContract(common.HexToAddress(token), parsed, ethereumConn, ethereumConn, ethereumConn), nil
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}

	res := make([]interface{}, 0)
	// 调用合约的 symbol 方法获取代币符号
	err = contract.Call(nil, &res, "symbol")
	if err != nil {
		log.Logger.Sugar().Error("GetContractSymbolOnMainNet err ", token, err)
		return err, ""
	}

	return nil, res[0].(string)
}

// CheckSymbolData Saving symbol data to redis if it has new symbol
// 检查 Redis 中的符号数据是否需要更新，如果需要，则保存到 Redis 中
func (s *TokenSymbol) CheckSymbolData(token, chainId, symbol string) (bool, error) {
	// 从 Redis 获取符号数据
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, err := db.RedisGet(redisKey)
	if len(redisTokenInfoBytes) <= 0 {
		// 如果 Redis 中没有数据，则从数据库中获取符号数据
		err = s.CheckTokenInfo(token, chainId)
		if err != nil {
			log.Logger.Error(err.Error())
		}
		err = db.RedisSet(redisKey, models.RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Symbol:  symbol,
		}, 0)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}
	} else {
		redisTokenInfo := models.RedisTokenInfo{}
		err = json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}
		// 比较符号是否一致
		if redisTokenInfo.Symbol == symbol {
			return false, nil
		}
		// 更新 Redis 数据
		redisTokenInfo.Symbol = symbol
		err = db.RedisSet(redisKey, redisTokenInfo, 0)
		if err != nil {
			log.Logger.Error(err.Error())
			return true, err
		}
	}
	return true, nil
}

// CheckTokenInfo  Insert token information if it was not in mysql
func (s *TokenSymbol) CheckTokenInfo(token, chainId string) error {
	tokenInfo := models.TokenInfo{}
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo = models.TokenInfo{}
			nowDateTime := utils.GetCurDateTimeFormat()
			tokenInfo.Token = token
			tokenInfo.ChainId = chainId
			tokenInfo.UpdatedAt = nowDateTime
			tokenInfo.CreatedAt = nowDateTime
			err = db.Mysql.Table("token_info").Create(tokenInfo).Debug().Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// SaveSymbolData Saving symbol data to mysql if it has new symbol
// 将符号数据保存到数据库中
func (s *TokenSymbol) SaveSymbolData(token, chainId, symbol string) error {
	nowDateTime := utils.GetCurDateTimeFormat()
	// 更新数据库记录
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=? ", token, chainId).Updates(map[string]interface{}{
		"symbol":     symbol,
		"updated_at": nowDateTime,
	}).Debug().Error
	if err != nil {
		log.Logger.Sugar().Error("UpdateContractSymbol SaveSymbolData err ", err)
		return err
	}

	return nil
}

package services

import (
	"encoding/json"
	"math/big"
	"pledge-backend/config"
	"pledge-backend/contract/bindings"
	"pledge-backend/db"
	"pledge-backend/log"
	"pledge-backend/schedule/models"
	"pledge-backend/utils"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/**
用于从区块链获取质押池的基础信息和动态数据，并将其存储到数据库和 Redis 中。
它通过以太坊客户端与智能合约交互，结合数据库和缓存，实现了质押池数据的高效同步和更新检测。
*/
// 作为服务的核心逻辑封装
type poolService struct{}

// 创建 poolService 实例
func NewPool() *poolService {
	return &poolService{}
}

// 用于更新所有质押池的信息
func (s *poolService) UpdateAllPoolInfo() {
	// 传入了测试网的合约地址、网络 URL 和链 ID
	s.UpdatePoolInfo(config.Config.TestNet.PledgePoolToken, config.Config.TestNet.NetUrl, config.Config.TestNet.ChainId)
	// 主网
	// s.UpdatePoolInfo(config.Config.MainNet.PledgePoolToken, config.Config.MainNet.NetUrl, config.Config.MainNet.ChainId)

}

// 更新指定质押池的基础信息和动态数据
func (s *poolService) UpdatePoolInfo(contractAddress, network, chainId string) {

	log.Logger.Sugar().Info("UpdatePoolInfo ", contractAddress+" "+network)
	//连接区块链网络
	ethereumConn, err := ethclient.Dial(network)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	//创建质押池智能合约的实例
	pledgePoolToken, err := bindings.NewPledgePoolToken(common.HexToAddress(contractAddress), ethereumConn)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	// 获取合约数据
	// borrowFee借款费率
	borrowFee, err := pledgePoolToken.PledgePoolTokenCaller.BorrowFee(nil)

	// lendFee借出费率
	lendFee, err := pledgePoolToken.PledgePoolTokenCaller.LendFee(nil)

	//poolLength质押池数量
	pLength, err := pledgePoolToken.PledgePoolTokenCaller.PoolLength(nil)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	// 遍历质押池,获取其基础信息和动态数据
	for i := 0; i <= int(pLength.Int64())-1; i++ {

		log.Logger.Sugar().Info("UpdatePoolInfo ", i)
		poolId := utils.IntToString(i + 1)
		// 获取基础信息
		baseInfo, err := pledgePoolToken.PledgePoolTokenCaller.PoolBaseInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Info("UpdatePoolInfo PoolBaseInfo err", poolId, err)
			continue
		}

		_, borrowToken := models.NewTokenInfo().GetTokenInfo(baseInfo.BorrowToken.String(), chainId)
		_, lendToken := models.NewTokenInfo().GetTokenInfo(baseInfo.LendToken.String(), chainId)

		lendTokenJson, _ := json.Marshal(models.LendToken{
			LendFee:    lendFee.String(),
			TokenLogo:  lendToken.Logo,
			TokenName:  lendToken.Symbol,
			TokenPrice: lendToken.Price,
		})
		borrowTokenJson, _ := json.Marshal(models.BorrowToken{
			BorrowFee:  borrowFee.String(),
			TokenLogo:  borrowToken.Logo,
			TokenName:  borrowToken.Symbol,
			TokenPrice: borrowToken.Price,
		})

		poolBase := models.PoolBase{
			SettleTime:             baseInfo.SettleTime.String(),
			PoolId:                 utils.StringToInt(poolId),
			ChainId:                chainId,
			EndTime:                baseInfo.EndTime.String(),
			InterestRate:           baseInfo.InterestRate.String(),
			MaxSupply:              baseInfo.MaxSupply.String(),
			LendSupply:             baseInfo.LendSupply.String(),
			BorrowSupply:           baseInfo.BorrowSupply.String(),
			MartgageRate:           baseInfo.MartgageRate.String(),
			LendToken:              baseInfo.LendToken.String(),
			LendTokenSymbol:        lendToken.Symbol,
			LendTokenInfo:          string(lendTokenJson),
			BorrowToken:            baseInfo.BorrowToken.String(),
			BorrowTokenSymbol:      borrowToken.Symbol,
			BorrowTokenInfo:        string(borrowTokenJson),
			State:                  utils.IntToString(int(baseInfo.State)),
			SpCoin:                 baseInfo.SpCoin.String(),
			JpCoin:                 baseInfo.JpCoin.String(),
			AutoLiquidateThreshold: baseInfo.AutoLiquidateThreshold.String(),
		}

		hasInfoData, byteBaseInfoStr, baseInfoMd5Str := s.GetPoolMd5(&poolBase, "base_info:pool_"+chainId+"_"+poolId)
		if !hasInfoData || (baseInfoMd5Str != byteBaseInfoStr) { // have new data
			//tokenInfo
			// 基础信息存储
			err = models.NewPoolBase().SavePoolBase(chainId, poolId, &poolBase)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolBase err ", chainId, poolId)
			}
			// 将数据的 MD5 值存储到 Redis 中，用于去重和缓存
			_ = db.RedisSet("base_info:pool_"+chainId+"_"+poolId, baseInfoMd5Str, 60*30) //The expiration time is set to prevent hsah collision
		}
		// 获取动态数据
		dataInfo, err := pledgePoolToken.PledgePoolTokenCaller.PoolDataInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Info("UpdatePoolInfo PoolBaseInfo err", poolId, err)
			continue
		}
		//计算质押池数据的 MD5 值，并与 Redis 中的缓存值进行比较
		hasPoolData, byteDataInfoStr, dataInfoMd5Str := s.GetPoolMd5(&poolBase, "data_info:pool_"+chainId+"_"+poolId)
		if !hasPoolData || (dataInfoMd5Str != byteDataInfoStr) { // have new data
			poolData := models.PoolData{
				PoolId:                 poolId,
				ChainId:                chainId,
				FinishAmountBorrow:     dataInfo.FinishAmountBorrow.String(),
				FinishAmountLend:       dataInfo.FinishAmountLend.String(),
				LiquidationAmounBorrow: dataInfo.LiquidationAmounBorrow.String(),
				LiquidationAmounLend:   dataInfo.LiquidationAmounLend.String(),
				SettleAmountBorrow:     dataInfo.SettleAmountBorrow.String(),
				SettleAmountLend:       dataInfo.SettleAmountLend.String(),
			}
			// 动态数据存储
			err = models.NewPoolData().SavePoolData(chainId, poolId, &poolData)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolData err ", chainId, poolId)
			}
			// 将数据的 MD5 值存储到 Redis 中，用于去重和缓存
			_ = db.RedisSet("data_info:pool_"+chainId+"_"+poolId, dataInfoMd5Str, 60*30) //The expiration time is set to prevent hsah collision
		}
	}
}

// //计算质押池数据的 MD5 值，并与 Redis 中的缓存值进行比较
func (s *poolService) GetPoolMd5(baseInfo *models.PoolBase, key string) (bool, string, string) {
	baseInfoBytes, _ := json.Marshal(baseInfo)
	baseInfoMd5Str := utils.Md5(string(baseInfoBytes))
	resInfoBytes, _ := db.RedisGet(key)
	if len(resInfoBytes) > 0 {
		return true, strings.Trim(string(resInfoBytes), `"`), baseInfoMd5Str
	} else {
		return false, strings.Trim(string(resInfoBytes), `"`), baseInfoMd5Str
	}
}

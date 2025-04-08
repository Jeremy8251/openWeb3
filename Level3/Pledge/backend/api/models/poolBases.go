package models

import (
	"encoding/json"
	"pledge-backend/db"
)

// 质押池的基础信息
type PoolBaseInfo struct {
	PoolID                 int             `json:"pool_id"`
	AutoLiquidateThreshold string          `json:"autoLiquidateThreshold"`
	BorrowSupply           string          `json:"borrowSupply"`
	BorrowToken            string          `json:"borrowToken"`
	BorrowTokenInfo        BorrowTokenInfo `json:"borrowTokenInfo"`
	EndTime                string          `json:"endTime"`
	InterestRate           string          `json:"interestRate"`
	JpCoin                 string          `json:"jpCoin"`
	LendSupply             string          `json:"lendSupply"`
	LendToken              string          `json:"lendToken"`
	LendTokenInfo          LendTokenInfo   `json:"lendTokenInfo"`
	MartgageRate           string          `json:"martgageRate"`
	MaxSupply              string          `json:"maxSupply"`
	SettleTime             string          `json:"settleTime"`
	SpCoin                 string          `json:"spCoin"`
	State                  string          `json:"state"`
}

// 数据库中的 poolbases 表结构
type PoolBases struct {
	Id                     int    `json:"-" gorm:"column:id;primaryKey"`
	PoolID                 int    `json:"pool_id" gorm:"column:pool_id;"`
	AutoLiquidateThreshold string `json:"autoLiquidateThreshold" gorm:"column:auto_liquidata_threshold;"`
	BorrowSupply           string `json:"borrowSupply" gorm:"column:borrow_supply;"`
	BorrowToken            string `json:"borrowToken" gorm:"column:pool_id;"`
	BorrowTokenInfo        string `json:"borrowTokenInfo" gorm:"column:borrow_token_info;"`
	EndTime                string `json:"endTime" gorm:"end_time;"`
	InterestRate           string `json:"interestRate" gorm:"column:interest_rate;"`
	JpCoin                 string `json:"jpCoin" gorm:"column:jp_coin;"`
	LendSupply             string `json:"lendSupply" gorm:"column:lend_supply;"`
	LendToken              string `json:"lendToken" gorm:"column:lend_token;"`
	LendTokenInfo          string `json:"lendTokenInfo" gorm:"column:lend_token_info;"`
	MartgageRate           string `json:"martgageRate" gorm:"column:martgage_rate;"`
	MaxSupply              string `json:"maxSupply" gorm:"column:max_supply;"`
	SettleTime             string `json:"settleTime" gorm:"column:settle_time;"`
	SpCoin                 string `json:"spCoin" gorm:"column:sp_coin;"`
	State                  string `json:"state" gorm:"column:state;"`
}

// 借贷代币
type BorrowTokenInfo struct {
	BorrowFee  string `json:"borrowFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

// 借贷代币的基础信息结构体
type LendTokenInfo struct {
	LendFee    string `json:"lendFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

// 质押池的基础信息结构体
type PoolBaseInfoRes struct {
	Index    int          `json:"index"`
	PoolData PoolBaseInfo `json:"pool_data"`
}

// 创建 PoolBases 实例
func NewPoolBases() *PoolBases {
	return &PoolBases{}
}

// 数据库表名 poolbases
func (p *PoolBases) TableName() string {
	return "poolbases"
}

// 根据链 ID 查询质押池的基础信息
func (p *PoolBases) PoolBaseInfo(chainId int, res *[]PoolBaseInfoRes) error {
	var poolBases []PoolBases
	// 查询数据库中的 poolbases 表，获取所有的质押池基础信息
	err := db.Mysql.Table("poolbases").Where("chain_id=?", chainId).Order("pool_id asc").Find(&poolBases).Debug().Error
	if err != nil {
		return err
	}

	for _, v := range poolBases {
		// 将 JSON 字符串解析为 BorrowTokenInfo 结构体
		borrowTokenInfo := BorrowTokenInfo{}
		_ = json.Unmarshal([]byte(v.BorrowTokenInfo), &borrowTokenInfo)
		// 将 JSON 字符串解析为 LendTokenInfo 结构体
		lendTokenInfo := LendTokenInfo{}
		_ = json.Unmarshal([]byte(v.LendTokenInfo), &lendTokenInfo)
		// 构造响应数据
		*res = append(*res, PoolBaseInfoRes{
			Index: v.PoolID - 1,
			PoolData: PoolBaseInfo{
				PoolID:                 v.PoolID,
				AutoLiquidateThreshold: v.AutoLiquidateThreshold,
				BorrowSupply:           v.BorrowSupply,
				BorrowToken:            v.BorrowToken,
				BorrowTokenInfo:        borrowTokenInfo,
				EndTime:                v.EndTime,
				InterestRate:           v.InterestRate,
				JpCoin:                 v.JpCoin,
				LendSupply:             v.LendSupply,
				LendToken:              v.LendToken,
				LendTokenInfo:          lendTokenInfo,
				MartgageRate:           v.MartgageRate,
				MaxSupply:              v.MaxSupply,
				SettleTime:             v.SettleTime,
				SpCoin:                 v.SpCoin,
				State:                  v.State,
			},
		})
	}
	return nil
}

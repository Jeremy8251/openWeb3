package services

import (
	"context"
	"fmt"
	"math/big"
	"pledge-backend/config"
	"pledge-backend/log"
	"pledge-backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

/*
*
监控区块链上某些账户或合约的余额。当余额低于设定的阈值时，系统会通过电子邮件发送警告通知。
它支持主网和测试网的余额监控，并提供了余额获取、阈值比较和邮件通知的功能
*/
type BalanceMonitor struct {
}

func NewBalanceMonitor() *BalanceMonitor {
	return &BalanceMonitor{}
}

// Monitor Sending email when balance is insufficient
// 监控余额不足时发送邮件通知
// 该函数会检查指定的合约地址在区块链上的余额，并与设定的阈值进行比较
func (s *BalanceMonitor) Monitor() {

	//check on bsc test-net 获取测试网账户的余额
	tokenPoolBalance, err := s.GetBalance(config.Config.TestNet.NetUrl, config.Config.TestNet.PledgePoolToken)
	thresholdPoolToken, ok := new(big.Int).SetString(config.Config.Threshold.PledgePoolTokenThresholdBnb, 10)
	//将余额与阈值进行比较，如果余额低于阈值，调用 EmailBody 方法生成邮件内容，
	if ok && (err == nil) && (tokenPoolBalance.Cmp(thresholdPoolToken) <= 0) {
		emailBody, err := s.EmailBody(config.Config.TestNet.PledgePoolToken, "TBNB", tokenPoolBalance.String(), thresholdPoolToken.String())
		if err != nil {
			log.Logger.Error(err.Error())
		} else {
			//通过 utils.SendEmail 发送通知。
			err = utils.SendEmail(emailBody, 2)
			if err != nil {
				log.Logger.Error(err.Error())
			}
		}
	}

	//check on bsc main-net
	// tokenPoolBalance, err = s.GetBalance(config.Config.MainNet.NetUrl, config.Config.MainNet.PledgePoolToken)
	// thresholdPoolToken, ok = new(big.Int).SetString(config.Config.Threshold.PledgePoolTokenThresholdBnb, 10)
	// if ok && (err == nil) && (tokenPoolBalance.Cmp(thresholdPoolToken) <= 0) {
	// 	emailBody, err := s.EmailBody(config.Config.MainNet.PledgePoolToken, "BNB", tokenPoolBalance.String(), thresholdPoolToken.String())
	// 	if err != nil {
	// 		log.Logger.Error(err.Error())
	// 	} else {
	// 		err = utils.SendEmail(emailBody, 2)
	// 		if err != nil {
	// 			log.Logger.Error(err.Error())
	// 		}
	// 	}
	// }
}

// GetBalance get balance of ERC20 token
// 获取指定地址的以太坊或ERC20代币余额
// netUrl: 区块链网络的URL，例如主网或测试网的URL
func (s *BalanceMonitor) GetBalance(netUrl, token string) (*big.Int, error) {
	// 连接到指定的区块链网络（主网或测试网）
	ethereumClient, err := ethclient.Dial(netUrl)
	if err != nil {
		log.Logger.Error(err.Error())
		return big.NewInt(0), err
	}
	defer ethereumClient.Close()
	// 调用 BalanceAt 方法获取指定地址的余额
	balance, err := ethereumClient.BalanceAt(context.Background(), common.HexToAddress(token), nil)
	if err != nil {
		log.Logger.Error(err.Error())
		return big.NewInt(0), err
	}

	return balance, err
}

// EmailBody email body
// 生成电子邮件的内容，包含余额不足的警告信息
func (s *BalanceMonitor) EmailBody(token, currency, balance, threshold string) ([]byte, error) {
	e18, err := decimal.NewFromString("1000000000000000000")
	if err != nil {
		return nil, err
	}

	balanceDeci, err := decimal.NewFromString(balance)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceDeci.Div(e18).String()

	thresholdDeci, err := decimal.NewFromString(threshold)
	if err != nil {
		return nil, err
	}

	thresholdStr := thresholdDeci.Div(e18).String()
	log.Logger.Sugar().Info("balance not enough ", token, " ", currency, " ", balanceStr, " ", thresholdStr)
	body := fmt.Sprintf(`<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;The balance of <strong><span style="color: rgb(255, 0, 0);"> %s </span></strong> is <strong>%s %s </strong>. Please recharge it in time. The current minimum balance limit is %s %s 
</p>`, token, balanceStr, currency, thresholdStr, currency)
	return []byte(body), nil
}

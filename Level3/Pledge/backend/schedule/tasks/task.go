package tasks

import (
	"pledge-backend/db"
	"pledge-backend/schedule/common"
	"pledge-backend/schedule/services"
	"time"

	"github.com/jasonlvhit/gocron"
)

func Task() {

	// get environment variables
	common.GetEnv()

	// flush redis db
	// 清空 Redis 数据库
	err := db.RedisFlushDB()
	if err != nil {
		panic("clear redis error " + err.Error())
	}

	//init task
	// 初始化任务,任务会在程序启动时立即执行一次
	services.NewPool().UpdateAllPoolInfo()           //更新所有质押池的信息
	services.NewTokenPrice().UpdateContractPrice()   //更新代币的价格信息
	services.NewTokenSymbol().UpdateContractSymbol() //更新代币的符号信息
	services.NewTokenLogo().UpdateTokenLogo()        //更新代币的 Logo 信息
	services.NewBalanceMonitor().Monitor()           //监控账户余额是否低于阈值
	// services.NewTokenPrice().SavePlgrPrice()
	services.NewTokenPrice().SavePlgrPriceTestNet() //将价格数据写入测试网的区块链合约

	//run pool task
	// 创建一个 gocron 调度器实例
	s := gocron.NewScheduler()
	// 设置调度器的时区为 UTC
	s.ChangeLoc(time.UTC)
	// 定时任务列表
	// 每 2 分钟执行一次，更新所有质押池的信息。
	_ = s.Every(2).Minutes().From(gocron.NextTick()).Do(services.NewPool().UpdateAllPoolInfo)
	// 每 1 分钟执行一次，更新代币的价格信息。
	_ = s.Every(1).Minute().From(gocron.NextTick()).Do(services.NewTokenPrice().UpdateContractPrice)
	// 每 2 小时执行一次，更新代币的符号信息。
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenSymbol().UpdateContractSymbol)
	// 每 2 小时执行一次，更新代币的 Logo 信息。
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenLogo().UpdateTokenLogo)
	// 每 30 分钟执行一次，监控账户余额是否低于阈值。
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewBalanceMonitor().Monitor)
	// 每 30 分钟执行一次，更新代币的价格信息。
	//_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPrice)
	// 每 30 分钟执行一次，将价格数据写入测试网的区块链合约。
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPriceTestNet)
	// 启动调度器，开始执行所有已注册的任务
	<-s.Start() // Start all the pending jobs

}

package kucoin

import (
	"pledge-backend/db"
	"pledge-backend/log"

	"github.com/Kucoin/kucoin-go-sdk"
)

// ApiKeyVersionV2 is v2 api key version指定 KuCoin API 的版本为 V2
const ApiKeyVersionV2 = "2"

// 当前的 PLGR 价格
var PlgrPrice = "0.0027"

// 用于传递实时价格数据的通道
var PlgrPriceChan = make(chan string, 2)

// 实时价格监听，并将最新的价格数据存储到 Redis
// 通过 KuCoin 的 WebSocket API 订阅 PLGR-USDT 的价格更新
func GetExchangePrice() {

	log.Logger.Sugar().Info("GetExchangePrice ")

	// get plgr price from redis
	//从 Redis 中获取 plgr_price 的值
	price, err := db.RedisGetString("plgr_price")
	if err != nil {
		// 获取失败，记录错误日志
		log.Logger.Sugar().Error("get plgr price from redis err ", err)
	} else {
		// 将价格更新到全局变量 PlgrPrice 中
		PlgrPrice = price
	}
	// 初始化 KuCoin API 服务, KuCoin将API密钥升级至2.0版本的验证逻辑
	s := kucoin.NewApiService(
		kucoin.ApiKeyOption("key"),
		kucoin.ApiSecretOption("secret"),
		kucoin.ApiPassPhraseOption("passphrase"),
		kucoin.ApiKeyVersionOption(ApiKeyVersionV2),
	)
	// 获取 WebSocket 的访问令牌
	rsp, err := s.WebSocketPublicToken()
	if err != nil {
		log.Logger.Error(err.Error()) // Handle error
		return
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 使用获取的令牌创建 WebSocket 客户端
	c := s.NewWebSocketClient(tk)
	// 调用 Connect 方法建立连接，返回两个通道：
	// mc：用于接收消息。
	// ec：用于接收错误。
	mc, ec, err := c.Connect()
	if err != nil {
		log.Logger.Sugar().Errorf("Error: %s", err.Error())
		return
	}

	// 创建订阅消息，订阅 PLGR-USDT 的价格更新
	ch := kucoin.NewSubscribeMessage("/market/ticker:PLGR-USDT", false)
	// 创建取消订阅消息
	uch := kucoin.NewUnsubscribeMessage("/market/ticker:PLGR-USDT", false)
	// 订阅失败，记录错误日志并退出
	if err := c.Subscribe(ch); err != nil {
		log.Logger.Error(err.Error()) // Handle error
		return
	}

	for {
		select {
		case err := <-ec:
			// 从错误通道 ec 接收到错误，停止 WebSocket 客户端并取消订阅
			c.Stop() // Stop subscribing the WebSocket feed
			log.Logger.Sugar().Errorf("Error: %s", err.Error())
			_ = c.Unsubscribe(uch)
			return
		case msg := <-mc:
			// 从消息通道 mc 接收价格更新消息。
			// 将消息解析为 TickerLevel1Model 类型，提取价格数据
			t := &kucoin.TickerLevel1Model{}
			if err := msg.ReadData(t); err != nil {
				log.Logger.Sugar().Errorf("Failure to read: %s", err.Error())
				return
			}
			// 将价格数据发送到 PlgrPriceChan 通道，并更新全局变量 PlgrPrice
			PlgrPriceChan <- t.Price
			PlgrPrice = t.Price
			//log.Logger.Sugar().Info("Price ", t.Price)
			// 将最新价格存储到 Redis 中，设置过期时间为 0（即不设置过期时间）
			_ = db.RedisSetString("plgr_price", PlgrPrice, 0)
		}
	}
}

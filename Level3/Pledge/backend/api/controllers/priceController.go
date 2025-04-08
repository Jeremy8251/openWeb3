package controllers

import (
	"net/http"
	"pledge-backend/api/models/ws"
	"pledge-backend/log"
	"pledge-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PriceController struct {
}

// 前端通过 WebSocket 连接到服务器，接收实时价格更新
func (c *PriceController) NewPrice(ctx *gin.Context) {
	// 使用 defer 和 recover 捕获运行时错误，避免程序因未处理的异常而崩溃
	defer func() {
		recoverRes := recover()
		if recoverRes != nil {
			log.Logger.Sugar().Error("new price recover ", recoverRes)
		}
	}()
	// 使用 websocket.Upgrader 将 HTTP 请求升级为 WebSocket 连接
	conn, err := (&websocket.Upgrader{
		// 设置读写缓冲区大小为 1024 字节
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second, //设置握手超时时间为 5 秒
		CheckOrigin: func(r *http.Request) bool { //Cross domain
			return true //允许跨域请求
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// 升级失败，记录错误日志并返回
		log.Logger.Sugar().Error("websocket request err:", err)
		return
	}

	randomId := ""
	remoteIP, ok := ctx.RemoteIP() //获取客户端的 IP 地址
	if ok {
		// 如果获取到 IP 地址，则将其转换为字符串并替换 "." 为 "_"，然后与随机字符串拼接
		randomId = strings.Replace(remoteIP.String(), ".", "_", -1) + "_" + utils.GetRandomString(23)
	} else {
		// 如果获取不到 IP 地址，则仅使用随机字符串作为 ID
		randomId = utils.GetRandomString(32)
	}
	// 创建一个 ws.Server 实例，表示当前客户端的 WebSocket 连接
	server := &ws.Server{
		Id:       randomId,
		Socket:   conn,
		Send:     make(chan []byte, 800),
		LastTime: time.Now().Unix(),
	}
	// 启动一个协程，调用 server.ReadAndWrite 方法，处理 WebSocket 的读写操作
	go server.ReadAndWrite()
}

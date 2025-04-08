package ws

import (
	"encoding/json"
	"errors"
	"pledge-backend/api/models/kucoin"
	"pledge-backend/config"
	"pledge-backend/log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const SuccessCode = 0
const PongCode = 1
const ErrorCode = -1

// 一个 WebSocket 客户端连接
type Server struct {
	sync.Mutex
	Id       string          // 客户端的唯一标识符
	Socket   *websocket.Conn // WebSocket 连接对象
	Send     chan []byte     // 用于发送消息的通道
	LastTime int64           // last send time 记录最后一次心跳时间，用于心跳检测
}

type ServerManager struct {
	Servers    sync.Map     //使用 sync.Map 管理所有活跃的客户端连接
	Broadcast  chan []byte  // 用于广播消息的通道
	Register   chan *Server // 注册新连接的通道
	Unregister chan *Server // 注销连接的通道
}

// 发送给客户端的消息结构，包含状态码和数据
type Message struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

// WebSocket 客户端管理器，负责管理所有连接的 WebSocket 客户端
var Manager = ServerManager{}

// 超时时间
var UserPingPongDurTime = config.Config.Env.WssTimeoutDuration // seconds
// 向客户端发送消息
func (s *Server) SendToClient(data string, code int) {
	s.Lock()
	defer s.Unlock()
	// 将消息封装为 Message 结构并序列化为 JSON
	dataBytes, err := json.Marshal(Message{
		Code: code,
		Data: data,
	})
	// 使用 WebSocket 的 WriteMessage 方法发送消息
	err = s.Socket.WriteMessage(websocket.TextMessage, dataBytes)
	if err != nil {
		// 发送失败，记录错误日志
		log.Logger.Sugar().Error(s.Id+" SendToClient err ", err)
	}
}

// 处理 WebSocket 消息的读写逻辑
func (s *Server) ReadAndWrite() {

	errChan := make(chan error)
	// 将客户端连接存储到 ServerManager 的 Servers 中
	Manager.Servers.Store(s.Id, s)
	// 在连接关闭时，删除客户端并释放资源
	defer func() {
		Manager.Servers.Delete(s)
		_ = s.Socket.Close()
		close(s.Send)
	}()

	//write
	go func() {
		for {
			select {
			// 从 Send 通道中读取消息
			case message, ok := <-s.Send:
				if !ok {
					errChan <- errors.New("write message error")
					return
				}
				// 消息发送给客户端
				s.SendToClient(string(message), SuccessCode)
			}
		}
	}()

	//read
	go func() {
		for {
			// 从 WebSocket 连接中读取消息
			_, message, err := s.Socket.ReadMessage()
			if err != nil {
				log.Logger.Sugar().Error(s.Id+" ReadMessage err ", err)
				errChan <- err
				return
			}

			//update heartbeat time 如果收到 ping 消息，更新心跳时间并回复 pong
			if string(message) == "ping" || string(message) == `"ping"` || string(message) == "'ping'" {
				s.LastTime = time.Now().Unix() // 更新心跳时间
				s.SendToClient("pong", PongCode)
			}
			continue

		}
	}()

	//check heartbeat 定期检查心跳时间，如果超时则关闭连接
	for {
		select {
		case <-time.After(time.Second):
			if time.Now().Unix()-s.LastTime >= UserPingPongDurTime {
				s.SendToClient("heartbeat timeout", ErrorCode)
				return
			}
		case err := <-errChan:
			log.Logger.Sugar().Error(s.Id, " ReadAndWrite returned ", err)
			return
		}
	}
}

// 实时监听价格更新并将其推送给所有已连接的 WebSocket 客户端
func StartServer() {
	log.Logger.Info("WsServer start")
	for {
		select {
		// 接收来自Kucoin的价格更新
		case price, ok := <-kucoin.PlgrPriceChan:
			if ok {
				// 成功接收到新的价格数据
				Manager.Servers.Range(func(key, value interface{}) bool {
					// 将 value 转换为 *Server 类型，表示一个 WebSocket 服务器实例
					// 调用 Server 的 SendToClient 方法，将接收到的价格数据和状态码
					value.(*Server).SendToClient(price, SuccessCode)
					return true
				})
			}
		}
	}
}

package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// 定义远程调用方法
type Hello struct {
}

func (h Hello) SayHello(req string, res *string) error {
	fmt.Println("req = ", req)
	*res = "你好" + req
	return nil
}

func main() {
	// 1.注册rpc服务
	err1 := rpc.RegisterName("hello", new(Hello))
	if err1 != nil {
		fmt.Println(err1)

	}
	// 2.监听端口
	listener, err2 := net.Listen("tcp", "127.0.0.1:8080")
	if err2 != nil {
		fmt.Println(err2)
	}
	// 3.应用退出时关闭监听端口
	defer listener.Close()
	for {
		fmt.Println("开始建立连接")
		// 4.建立连接
		conn, err3 := listener.Accept()
		if err3 != nil {
			fmt.Println(err3)
		}
		// 5.绑定服务
		rpc.ServeConn(conn)
	}

}

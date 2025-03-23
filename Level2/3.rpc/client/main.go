package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 连接微服务
	conn, err1 := rpc.Dial("tcp", "127.0.0.1:8080")
	if err1 != nil {
		fmt.Println(err1)
	}
	// 关闭连接
	defer conn.Close()

	// 调用微服务方法
	var reply string
	err2 := conn.Call("hello.SayHello", "我是客户端-aaa", &reply)
	if err2 != nil {
		fmt.Println(err2)
	}
	// 获取微服务返回的值
	fmt.Println(reply)

}

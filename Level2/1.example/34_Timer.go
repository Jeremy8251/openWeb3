package main

import (
	"fmt"
	"time"
)

func testTimer() {
	//‌创建定时器‌：timer1会在2秒后触发。
	timer1 := time.NewTimer(2 * time.Second)

	//‌阻塞等待触发‌：主goroutine通过<-timer1.C等待定时器的通道C，2秒后触发，打印"Timer 1 fired"。
	<-timer1.C
	fmt.Println("Timer 1 fired")

	//创建定时器‌：timer2会在1秒后触发。
	timer2 := time.NewTimer(time.Second)
	//启动匿名goroutine‌：在新的goroutine中等待timer2触发，触发后打印"Timer 2 fired"。
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	//停止定时器‌：主goroutine立即调用timer2.Stop()。
	stop2 := timer2.Stop()
	if stop2 {
		//如果成功停止（定时器未触发），Stop()返回true，打印"Timer 2 stopped"。
		fmt.Println("Timer 2 stopped")
	}
	//主goroutine休眠2秒，确保程序不会立即退出，给其他goroutine执行机会（但此处定时器已停止，不会有输出）。
	time.Sleep(2 * time.Second)
}

package main

import (
	"fmt"
	"time"
)

func testTicker() {

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	//这里我们使用通道内建的 select，等待每 500ms 到达一次的值。
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

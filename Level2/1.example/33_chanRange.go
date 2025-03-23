package main

import "fmt"
//一个非空的通道也是可以关闭的， 并且，通道中剩下的值仍然可以被接收到
func testChanRange() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)
	for elem := range queue {
		fmt.Println(elem)
	}
}

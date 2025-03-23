package main

import "fmt"

func testChannels() {
	message := make(chan string)

	go func() {
		message <- "ping"
		fmt.Println("func ping")
	}()

	msg := <-message
	fmt.Println(msg)
}

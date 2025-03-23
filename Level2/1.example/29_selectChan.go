package main

import (
	"fmt"
	"time"
)

func testSelectChan() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			c1 <- "one"
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(2 * time.Second)
			c2 <- "two"
		}

	}()

	// for i := 0; i < 2; i++
	for {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		case <-time.After(10 * time.Second):
			fmt.Println("received is timeout")
			return
		}

	}

}

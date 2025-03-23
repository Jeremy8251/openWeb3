package main

import (
	"fmt"
	"sync"
	"time"
)

func workerfunc(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func testWaitGroup() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			workerfunc(i)
		}()

	}
	wg.Wait()
}

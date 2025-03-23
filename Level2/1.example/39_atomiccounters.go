package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func testAtom() {

	var ops uint64

	var wg sync.WaitGroup
	//启动 50 个协程，并且每个协程会将计数器递增 1000 次。
	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				//使用 AddUint64 来让计数器自动增加， 使用 & 语法给定 ops 的内存地址
				atomic.AddUint64(&ops, 1) //5000
				//非原子的 ops++ 来增加计数器， 由于多个协程会互相干扰，运行时值会改变
				//ops++ //ops: 46557 || ops: 50142
			}
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("ops:", ops)
}

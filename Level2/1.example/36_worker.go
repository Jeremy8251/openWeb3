package main

import (
	"fmt"
	"time"
)

func workefunc(id int, jobs <-chan int, results chan<- int) {
	fmt.Println("worker jobs = ", len(jobs))
	for j := range jobs { //启动3个worker协程（此时jobs通道为空，worker阻塞在for j := range jobs）
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func testWorker() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go workefunc(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		fmt.Println("numJobs started  job", j)
		jobs <- j
	}
	close(jobs) //close(jobs)触发worker退出for range循环

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

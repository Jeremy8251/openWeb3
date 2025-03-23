package main

import "fmt"

// 变参函数
func testMultiParams() {
	sum(1, 2)
	sum(1, 2, 3)
	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
func sum(nums ...int) {
	fmt.Print(nums)
	total := 0
	for _, v := range nums {
		total += v
	}
	fmt.Println("=", total)
}

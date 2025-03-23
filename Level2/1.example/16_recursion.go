package main

import "fmt"
//  递归
func testRecursion() {
	fmt.Println("1: ", fact(7))

	var fib func(n int) int
	fib = fact
	fmt.Println("2: ", fact(7))

	fib = func(n int) int {
		if n < 2 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	fmt.Println("1: ", fib(7))

}
func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

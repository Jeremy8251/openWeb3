package main

import "fmt"

// 闭包
func testClosure() {
	// nextInt := intSeq()
	// fmt.Println(nextInt())
	// fmt.Println(nextInt())
	// fmt.Println(nextInt())

	// newInts := intSeq()
	// fmt.Println(newInts())

	nextInt2 := intSeq2()
	fmt.Println(nextInt2(10))
	fmt.Println(nextInt2(10))
}
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func intSeq2() func(int) int {
	i := 0
	return func(x int) int {
		i = i + x
		return i
	}
}

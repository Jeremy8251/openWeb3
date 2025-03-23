package main

import "fmt"

// 多返回值
func testMultiReturn() {
	a, b := vals()
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)

	_, res := vals()
	fmt.Println("res = ", res)
}
func vals() (int, int) {
	return 3, 7
}

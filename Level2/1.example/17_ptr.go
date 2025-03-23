package main

import "fmt"

//  指针
func testPtr() {
	i := 1
	fmt.Println("initial:", i)

	zeroval(i)
	fmt.Println("zeroval:", i)

	zeroPtr(&i)
	fmt.Println("zeroPtr:", i)

	fmt.Println("pointer:", &i)
}
func zeroval(n int) {
	n = 0
}

func zeroPtr(nptr *int) {
	*nptr = 0
}

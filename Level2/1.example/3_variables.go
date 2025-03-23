package main

import "fmt"

func testVariables() {
	var a = "initial"
	fmt.Println("a = ", a)

	var b, c int = 1, 2
	fmt.Println("b=", b)
	fmt.Println("c =", c)

	var d = true
	fmt.Println("d =", d)

	var e int
	fmt.Println("e =", e)

	f := "apple"
	fmt.Println("f =", f)
}

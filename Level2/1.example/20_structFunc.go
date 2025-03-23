package main

import "fmt"

type rect struct {
	width, heigth int
}

//Go 语言允许你在调用方法时使用自动转换机制，即使接收者类型不匹配，Go 也会自动进行转换。
func (r *rect) area() int {

	return r.width * r.heigth
}

func (r rect) perim() int {

	return 2*r.width + 2*r.heigth
}

func testStructFunc() {
	r := rect{width: 10, heigth: 5}
	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perim:", rp.perim())
}

package main

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect2 struct {
	width, heigth float64
}

type circle struct {
	radius float64
}

func (r rect2) area() float64 {
	return r.width * r.heigth
}
func (r rect2) perim() float64 {
	return 2*r.width + 2*r.heigth
}
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}
func meansure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}
func testInterface() {
	meansure(rect2{width: 3, heigth: 4})
	meansure(circle{radius: 5})
}

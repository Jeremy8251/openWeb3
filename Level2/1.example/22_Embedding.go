package main

import "fmt"

//Go支持对于结构体(struct)和接口(interfaces)的 嵌入(embedding)
// 以表达一种更加无缝的 组合(composition) 类型

type base struct {
	num int
}

type container struct {
	base
	str  string
}

func (b base) describe() string {
	return fmt.Sprintf("base with num=%v", b.num)
}

func testEmbedding() {
	co := container{
		base: base{num: 1},
		str:  "my name",
	}
	fmt.Printf("co={num: %v, str: %v}\n", co.num, co.str)
	fmt.Println("also num:", co.base.num)
	fmt.Println("describe:", co.describe())

	type describer interface {
		describe() string
	}
	var d describer = co
	fmt.Println("describer:", d.describe())
}

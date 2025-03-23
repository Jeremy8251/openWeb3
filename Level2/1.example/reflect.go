package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string `json:"name" db:"name"`
	Age   int    `json:"age" db:"age"`
	Score int    `json:"score" db:"score"`
}

func (s Student) GetName() string {
	return s.Name
}
func (s *Student) SetName(name string) {
	s.Name = name
}

func changeStruct(a interface{}) {
	t := reflect.TypeOf(a)
	if t.Kind() != reflect.Struct && t.Elem().Kind() != reflect.Struct {
		fmt.Println("不是结构体")
		return
	} else {
		fmt.Println("是结构体")
	}
	fmt.Println("t.Kind:", t.Kind() != reflect.Struct)
	fmt.Println("Type:", t)
	fmt.Printf("类型名称：%v, 类型：%v\n", t.Name(), t.Kind())
	if t.Kind() == reflect.Struct {
		// 值类型
		typeCount := t.NumField()
		fmt.Println("字段数：", typeCount)
	} else {
		// 指针类型
		typeCount := t.Elem().NumField()
		fmt.Println("Elem字段数：", typeCount)
	}

	v := reflect.ValueOf(a)
	fmt.Printf("值名称：%v, 值类型：%v\n", v, v.Type())

	fmt.Println("方法数：", t.NumMethod())
	m := t.Method(0)
	fmt.Printf("方法名称：%v, 方法类型：%v\n", m.Name, m.Type)

	m, ok := t.MethodByName("GetName")
	if !ok {
		fmt.Println("没有这个方法")
	}
	tcall := m.Func.Call([]reflect.Value{reflect.ValueOf(a)})
	fmt.Println("tcall = ", tcall)

	name := v.MethodByName("GetName").Call(nil)
	fmt.Println("v name = ", name)

	if t.Kind() == reflect.Pointer {
		v.MethodByName("SetName").Call([]reflect.Value{reflect.ValueOf("李四")})
	}
	name = v.MethodByName("GetName").Call(nil)
	fmt.Println("new name = ", name)

	nameSet := v.Elem().FieldByName("Name")
	nameSet.SetString("小王")

	ageSet := v.Elem().FieldByName("Age")
	ageSet.SetInt(60)

}
func testReflect() {

	stu := Student{
		Name:  "张三",
		Age:   20,
		Score: 90,
	}
	changeStruct(&stu)
	fmt.Println(stu)

}

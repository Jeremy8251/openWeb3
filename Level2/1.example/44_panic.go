package main

import (
	"errors"
	"fmt"
)

func testPanic() {
	// divideByZero(10, 0)
	read()
}

func divideByZero(a int, b int) int {
	defer func() {
		error := recover()
		if error != nil {
			fmt.Println("error:", error) //error: runtime error: integer divide by zero
		}
	}()
	// panic("have a problem")
	return a / b
}

func read() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:", err)
		}
	}()
	err := readFile("XX.go")
	if err != nil {
		panic(err)
	}
}
func readFile(fileName string) error {
	if fileName == "main.go" {
		return nil
	} else {
		return errors.New("filename is error")
	}
	// _, err := os.Create("/tmp/file")
	// if err != nil {
	// 	panic(err)
	// }
}

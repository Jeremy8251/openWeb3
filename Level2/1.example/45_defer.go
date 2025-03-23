package main

import (
	"fmt"
	"os"
)

func testDefer() {
	// f := createFile("/tmp/defer.txt")
	// defer closeFile(f)
	// writeFile(f)
	// fmt.Println("test1 = ", test1()) //0
	// fmt.Println("test2 = ", test2()) //1
	// fmt.Println("test3 = ", test3()) //5
	fmt.Println("test4 = ", test4()) //5
}

func test1() int {
	var a int
	defer func() {
		a++
	}()
	return a
}

func test2() (a int) {
	// var a int
	defer func() {
		a++
	}()
	return a
}

func test3() (y int) {
	a := 5
	defer func() {
		a++
	}()
	return a
}

func test4() (x int) {
	defer func(x int) {
		fmt.Println("defer x = ", x)
		x++
	}(x) // 内部不影响外部x=0
	fmt.Println("test4 x = ", x)
	return 5
}

func createFile(s string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(s)
	if err != nil {
		panic(err)
	}
	return f
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func writeFile(f *os.File) {
	fmt.Println("writing")
	fmt.Fprintln(f, "data")
}

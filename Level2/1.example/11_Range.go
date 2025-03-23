package main

import "fmt"

func testRange() {

	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum = ", sum)

	for index, num := range nums {
		if num == 3 {
			fmt.Println("index = ", index)
		}
	}

	kvs := map[string]string{"a": "apple", "b": "banana"}
	for key, val := range kvs {
		fmt.Printf("%s -> %s\n", key, val)

	}
	for k, _ := range kvs {
		fmt.Println("key = ", k)
	}

	for i, c := range "学习go" {
		fmt.Println(i, c)
		fmt.Printf("index = %d, char = %c \n", i, c)
	}

}

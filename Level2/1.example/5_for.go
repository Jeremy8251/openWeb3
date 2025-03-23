package main

import "fmt"

func testFor() {
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i += 1
	}
	fmt.Println("====")
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}
	fmt.Println("====")
	for {
		fmt.Println("loop")
		break
	}
	fmt.Println("====")
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}

}

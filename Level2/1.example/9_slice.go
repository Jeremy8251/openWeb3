package main

import "fmt"

func testSlice() {
	s := make([]int, 3, 6)
	fmt.Println("initial, s =", s) //[0 0 0]
	s[1] = 2
	fmt.Println("after set position 1, s =", s) //[0 2 0]

	s2 := append(s, 4)
	fmt.Println("after append, s2 length:", len(s2))   //4
	fmt.Println("after append, s2 capacity:", cap(s2)) //6
	fmt.Println("after append, s =", s)                //[0 2 0]
	fmt.Println("after append, s2 =", s2)              //[0 2 0 4]

	s[0] = 1024
	fmt.Println("after set position 0, s =", s)   //[1024 2 0]
	fmt.Println("after set position 0, s2 =", s2) //[1024 2 0 4]

	appendInFunc(s)
	fmt.Println("after append in func, s =", s)   //[1024 2 512]?
	fmt.Println("after append in func, s2 =", s2) //[1024 2 512 1022]

	st := append(s2, 2048)
	fmt.Println("after append2048 in func, s =", s)   //[1024 2 512]
	fmt.Println("after append2048 in func, s2 =", s2) //[1024 2 512 1022]
	fmt.Println("after append2048 in func, st =", st) // [1024 2 512 1022 2048]

	appendInFunc2(s2)
	fmt.Println("after appendInFunc2 in func, s =", s)        //[1024 2 512]
	fmt.Println("after appendInFunc2 in func, s2 =", s2)      //[1024 2 512 1022]
	fmt.Println("after appendInFunc2 in func, s2: =", s2[2:]) //[512 1022]
	fmt.Println("after appendInFunc2 in func, st =", st)      //[1024 2 512 1022 4096]
	fmt.Println("after appendInFunc2 in func, st: =", st[2:]) //[512 1022 4096]

	s[2] = 8192
	fmt.Println("after [2]=8192 in func, s =", s)        //[1024 2 8192]
	fmt.Println("after [2]=8192 in func, s2 =", s2)      //[1024 2 8192 1022]
	fmt.Println("after [2]=8192 in func, s2: =", s2[2:]) //[8192 1022]
	fmt.Println("after [2]=8192 in func, st =", st)      //[1024 2 8192 1022 4096]
	fmt.Println("after [2]=8192 in func, st: =", st[2:]) //[8192 1022 4096]
}

func appendInFunc(param []int) {
	param = append(param, 1022)
	fmt.Println("in func, param =", param) //[1024 2 0 1022]
	param[2] = 512
	fmt.Println("set position 2 in func, param =", param) //[1024 2 512 1022]
}

func appendInFunc2(param []int) {
	param = append(param, 4096)
	fmt.Println("in func, param =", param) //[1024 2 512 1022 4096]
}

func testSlice2() {
	s := make([]string, 3)
	fmt.Println("emp = ", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	fmt.Println("len:", len(s))

	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("append:", s)

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("copy:", c)

	l := s[2:5]
	fmt.Println("l:", l)

	l = s[:5]
	fmt.Println("l2:", l)

	l = s[2:]
	fmt.Println("l3:", l)

	t := []string{"g", "h", "l"}
	fmt.Println("t:", t)

	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j

		}

	}
	fmt.Println("twoD: ", twoD)

}

// 排序
func testSlice3() {
	var numSlice = []int{9, 8, 7, 6, 5, 4}
	for i := 0; i < len(numSlice); i++ {
		for j := i + 1; j < len(numSlice); j++ {
			if numSlice[i] > numSlice[j] {
				temp := numSlice[i]
				numSlice[i] = numSlice[j]
				numSlice[j] = temp
			}
		}
		fmt.Println(numSlice)

	}
}

package main

import (
	"fmt"
	"sort"
)

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func testSortByFun() {
	fruits := []string{"peach", "banana", "kiwi"}
	//实现了 sort.Interface 接口的 Len、Less 和 Swap 方法
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}

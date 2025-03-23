package main

import (
	"fmt"
	"unicode/utf8"
)

// 计算字符串中有多少rune，
func testStringRune() {
	const s = "สวัสดี" //泰语中的单词 “hello”
	fmt.Println("Len:", len(s))
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i])
	}
	fmt.Println()
	fmt.Println("Rune count:", utf8.RuneCountInString(s))

	for idx, runeValue := range s {
		fmt.Printf("%#U starts at %d\n", runeValue, idx)

	}
	fmt.Println("\nUsing DecodeRuneInString")
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%#U starts at %d\n", runeValue, i)
		w = width
		exaimineRune(runeValue)
	}

}

func exaimineRune(r rune) {
	if r == 't' {
		fmt.Println("found tee")

	} else if r == 'ส' {
		fmt.Println("found so sua")
	}
}

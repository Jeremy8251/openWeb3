package models

import (
	"fmt"
	"time"
)

// 时间戳转换成日期函数
func UnixToTime(timestamp int64) string {
	fmt.Println("timestamp = ", timestamp)
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func Println(str1 string, str2 string) string {
	return str1 + str2
}

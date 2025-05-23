package main

import (
	"fmt"
	"time"
)

func timeStamp() {
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()

	fmt.Println("now = ", now)

	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)

	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}

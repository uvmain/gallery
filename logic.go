package main

import (
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func GetUnixTimeString() string {
	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)

	nanoTime := time.Now().Nanosecond()
	nanoTimeString := strconv.Itoa(nanoTime)
	return unixTimeString + nanoTimeString
}

package main

import (
	"strconv"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	_ "modernc.org/sqlite"
)

func GenerateSlug() string {
	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)

	nanoTime := time.Now().Nanosecond()
	nanoTimeString := strconv.Itoa(nanoTime)
	return unixTimeString + nanoTimeString
}

func ToTitle(inputString string) string {
	toTitle := cases.Title(language.English)
	return toTitle.String(inputString)
}

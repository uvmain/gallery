package logic

import (
	"fmt"
	"gallery/core/config"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}

var (
	bootTime     time.Time
	bootTimeOnce sync.Once
)

func GetBootTime() time.Time {
	bootTimeOnce.Do(func() {
		bootTime = time.Now().UTC().Truncate(time.Second)
	})
	return bootTime
}

func GenerateSlug() string {
	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)

	nanoTime := time.Now().Nanosecond()
	nanoTimeString := strconv.Itoa(nanoTime)
	return unixTimeString + nanoTimeString
}

func CreateDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		log.Printf("Creating directory: %s", directoryPath)
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Printf("Error creating directory%s: %s", directoryPath, err)
		} else {
			log.Printf("Directory created: %s", directoryPath)
		}
	} else {
		log.Printf("Directory already exists: %s", directoryPath)
	}
}

func GetDirContents(directoryPath string) ([]string, error) {
	var foundFiles []string

	absPath, _ := filepath.Abs(directoryPath)
	CreateDir(absPath)

	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error opening directory %s: %s", directoryPath, err)
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, validExt := range config.ImageExtensions {
				if ext == validExt {
					foundFiles = append(foundFiles, path)
					break
				}
			}
		}
		return nil
	})
	log.Printf("Found: %d images in %s", len(foundFiles), directoryPath)
	return foundFiles, err
}

func TernaryString(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
}

func StringArraySortUnique(arrayToSort []string) []string {
	slices.Sort(arrayToSort)
	arrayToSort = slices.Compact(arrayToSort)
	return arrayToSort
}

var dateFormats = []string{
	time.RFC3339,                    // "2025-01-24T18:15:21Z"
	"2006-01-02T15:04",              // "2025-01-29T17:21"
	"2006-01-02 15:04:05 -0700 MST", // "2018-06-23 16:05:18 +0100 BST"
	"2006-01-02 15:04:05 -0700",     // "2018-06-23 16:05:18 +0100"
	"2006-01-02 15:04",              // "2020-01-31T00:00"
}

func FormatTimeToString(dateString string) (string, error) {
	for _, layout := range dateFormats {
		t, err := time.Parse(layout, dateString)
		if err == nil {
			formattedTime := t.Format("2006-01-02 15:04:05")
			return formattedTime, nil
		}
	}
	log.Printf("unsupported date format: %s", dateString)
	return "0000-00-00 00:00:00", fmt.Errorf("unsupported date format: %s", dateString)
}

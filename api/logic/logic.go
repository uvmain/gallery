package logic

import (
	"log"
	"os"
	"path/filepath"
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
			for _, validExt := range ImageExtensions {
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

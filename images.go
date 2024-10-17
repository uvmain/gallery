package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

var FoundFiles []string

func GetImageDirContents() ([]string, error) {
	var foundFiles []string

	err := filepath.Walk(ImagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error creating thumbnails directory: %s", err)
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
	log.Printf("Found: %d images", len(foundFiles))
	FoundFiles = foundFiles

	var photoInt int = 0
	log.Printf(`Found: %s`, FoundFiles[photoInt])
	getExif(FoundFiles[photoInt])
	GenerateThumbnail(FoundFiles[photoInt], GetUnixTimeString())

	return foundFiles, err
}

func getExif(filepath string) {

	ExposureModes := map[int]string{
		0: "Not defined",
		1: "Manual",
		2: "Normal program",
		3: "Aperture priority",
		4: "Shutter priority",
		5: "Creative program",
		6: "Action program",
		7: "Portrait mode",
		8: "Landscape mode",
	}

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", filepath, err)
	}
	exifData, _ := exif.Decode(f)

	exposureModeTag, _ := exifData.Get("ExposureProgram")
	exposureModeInt, _ := exposureModeTag.Int(0)
	exposureMode := ExposureModes[exposureModeInt]
	dateTaken, _ := exifData.DateTime()

	log.Printf(`Exposure Mode: %s`, exposureMode)
	log.Printf(`Date Taken: %s`, dateTaken)
}

// func insertExifIfNeeded(filepath string) {

// }

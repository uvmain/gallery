package main

import (
	"log"
	"os"
	"path/filepath"
)

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
	return foundFiles, err
}

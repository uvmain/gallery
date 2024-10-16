package main

import (
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func CreateThumbnailsDir() {
	if _, err := os.Stat(ThumbnailDirectory); os.IsNotExist(err) {
		log.Println("Creating thumbnails directory")
		err := os.MkdirAll(ThumbnailDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating thumbnails directory")
		} else {
			log.Println("Thumbnails directory created")
		}
	} else {
		log.Println("Thumbnail directory already exists")
	}
}

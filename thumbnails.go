package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	_ "modernc.org/sqlite"
)

func CreateThumbnailsDir() {
	if _, err := os.Stat(ThumbnailDirectory); os.IsNotExist(err) {
		log.Println("Creating thumbnails directory")
		err := os.MkdirAll(ThumbnailDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating thumbnails directory %s", err)
		} else {
			log.Println("Thumbnails directory created")
		}
	} else {
		log.Println("Thumbnail directory already exists")
	}
}

func GenerateThumbnail(imageFile string, slug string) error {
	source, err := imaging.Open(imageFile)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	width, height := 0, 0

	if source.Bounds().Max.X > source.Bounds().Max.Y {
		width = 0
		height = int(ThumbnailMaxPixels)
	} else {
		width = int(ThumbnailMaxPixels)
		height = 0
	}

	thumbnailPath := filepath.Join(ThumbnailDirectory, slug) + ".jpeg"

	thumbnailImage := imaging.Resize(source, width, height, imaging.Lanczos)

	err = imaging.Save(thumbnailImage, thumbnailPath)
	if err != nil {
		log.Fatalf("Error creating thumbnail: %s", err)
		return err
	}
	return nil
}

// func GenerateThumbnails() error {
// 	filesToCheck := GetExistingMetadataFilePaths()
// 	for _, imagePath := range  {
// 		if ext == validExt {
// 			foundFiles = append(foundFiles, path)
// 			break
// 		}
// 	}
// 	return nil
// }

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	_ "modernc.org/sqlite"
)

func CreateOptimisedDir() {
	if _, err := os.Stat(OptimisedDirectory); os.IsNotExist(err) {
		log.Println("Creating Optimised directory")
		err := os.MkdirAll(OptimisedDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating Optimised directory %s", err)
		} else {
			log.Println("Optimised directory created")
		}
	} else {
		log.Println("Optimised directory already exists")
	}
}

func GenerateOptimised(imageFile string, slug string) error {
	source, err := imaging.Open(imageFile)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	width, height := 0, 0

	if source.Bounds().Max.X > source.Bounds().Max.Y {
		width = 0
		height = int(OptimisedMaxPixels)
	} else {
		width = int(OptimisedMaxPixels)
		height = 0
	}

	OptimisedPath := filepath.Join(OptimisedDirectory, slug) + ".jpeg"

	OptimisedImage := imaging.Resize(source, width, height, imaging.Lanczos)

	err = imaging.Save(OptimisedImage, OptimisedPath)
	if err != nil {
		log.Fatalf("Error creating Optimised: %s", err)
		return err
	}
	return nil
}

// func GenerateOptimised() error {
// 	filesToCheck := GetExistingMetadataFilePaths()
// 	for _, imagePath := range  {
// 		if ext == validExt {
// 			foundFiles = append(foundFiles, path)
// 			break
// 		}
// 	}
// 	return nil
// }

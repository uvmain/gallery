package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
)

func optimisedAlreadyExists(slug string) bool {
	optimisedPath := filepath.Join(OptimisedDirectory, (slug + "." + ImageFormat))
	if _, err := os.Stat(optimisedPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateOptimised(imageFile string, slug string) {

	if optimisedAlreadyExists(slug) {
		return
	}

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

	var optimisedPath string

	if ImageFormat == "jpeg" || ImageFormat == "jpg" {
		optimisedPath = filepath.Join(OptimisedDirectory, slug) + ".jpeg"
		optimisedImage := imaging.Resize(source, width, height, imaging.Lanczos)
		err = imaging.Save(optimisedImage, optimisedPath)
	} else if ImageFormat == "webp" {
		optimisedPath = filepath.Join(OptimisedDirectory, slug) + ".webp"
		optimisedImage := imaging.Resize(source, width, height, imaging.Lanczos)

		f, _ := os.Create(optimisedPath)
		defer f.Close()

		webp.Encode(f, optimisedImage)
	}

	if err != nil {
		log.Printf("Error creating optimised: %s", err)
	}
	log.Printf("Optimised created for %s: %s", imageFile, optimisedPath)
}

func getOptimisedDirContents() ([]string, error) {
	var foundOptimised []string

	err := filepath.Walk(OptimisedDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error opening Optimised directory: %s", err)
			return err
		}
		if !info.IsDir() {
			foundOptimised = append(foundOptimised, path)
		}
		return nil
	})
	FoundOptimised = foundOptimised
	log.Printf("Found %d optimised", len(foundOptimised))
	return foundOptimised, err
}

func populateOptimised() {
	for _, row := range GetExistingMetadataFilePaths() {
		slug := row.slug
		filePath := row.filePath
		fileName := row.fileName
		imageFullPath := filepath.Join(filePath, fileName)
		generateOptimised(imageFullPath, slug)
	}
}

func deleteExtraneousOptimised() {
	optimisedDirContents, _ := getOptimisedDirContents()
	for _, optimised := range optimisedDirContents {
		ext := strings.Split(filepath.Ext(optimised), ".")[1]
		if ext != ImageFormat {
			deleteOptimisedByFilename(optimised)
		} else {
			slug := strings.TrimSuffix(filepath.Base(optimised), filepath.Ext(optimised))
			_, err := GetMetadataBySlug(slug)
			if err != nil {
				log.Println(slug)
				deleteOptimisedByFilename(optimised)
			}
		}
	}
}

func deleteOptimisedByFilename(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Printf("Error deleting optimised %s: %s", filename, err)
		return
	}
	log.Printf("Optimised %s deleted", filename)
}

func InitialiseOptimised() {
	CreateDir(OptimisedDirectory)
	getOptimisedDirContents()
	populateOptimised()
	deleteExtraneousOptimised()
}

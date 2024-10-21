package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
)

func thumbnailAlreadyExists(slug string) bool {
	thumbnailPath := filepath.Join(ThumbnailDirectory, (slug + "." + ImageFormat))
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateThumbnail(imageFile string, slug string) {
	if thumbnailAlreadyExists(slug) {
		return
	}

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

	var thumbnailPath string

	if ImageFormat == "jpeg" || ImageFormat == "jpg" {
		thumbnailPath = filepath.Join(ThumbnailDirectory, slug) + ".jpeg"
		thumbnailImage := imaging.Resize(source, width, height, imaging.Lanczos)
		err = imaging.Save(thumbnailImage, thumbnailPath)
	} else if ImageFormat == "webp" {
		thumbnailPath = filepath.Join(ThumbnailDirectory, slug) + ".webp"
		thumbnailImage := imaging.Resize(source, width, height, imaging.Lanczos)

		f, _ := os.Create(thumbnailPath)
		defer f.Close()

		webp.Encode(f, thumbnailImage)
	}

	if err != nil {
		log.Printf("Error creating thumbnail: %s", err)
	}
	log.Printf("Thumbnail created for %s: %s", imageFile, thumbnailPath)
}

func getThumbnailDirContents() ([]string, error) {
	var foundThumbnails []string

	err := filepath.Walk(ThumbnailDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error opening Thumbnails directory: %s", err)
			return err
		}
		if !info.IsDir() {
			foundThumbnails = append(foundThumbnails, path)
		}
		return nil
	})
	FoundThumbnails = foundThumbnails
	log.Printf("Found %d thumbnails", len(foundThumbnails))
	return foundThumbnails, err
}

func populateThumbnails() {
	for _, row := range GetExistingMetadataFilePaths() {
		slug := row.slug
		filePath := row.filePath
		fileName := row.fileName
		imageFullPath := filepath.Join(filePath, fileName)
		go generateThumbnail(imageFullPath, slug)
	}
}

func deleteExtraneousThumbnails() {
	thumbnailDirContents, _ := getThumbnailDirContents()
	for _, thumbnail := range thumbnailDirContents {
		ext := strings.Split(filepath.Ext(thumbnail), ".")[1]
		if ext != ImageFormat {
			deleteThumbnailByFilename(thumbnail)
		} else {
			slug := strings.TrimSuffix(filepath.Base(thumbnail), filepath.Ext(thumbnail))
			_, err := GetMetadataBySlug(slug)
			if err != nil {
				log.Println(slug)
				deleteThumbnailByFilename(thumbnail)
			}
		}
	}
}

func deleteThumbnailByFilename(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Printf("Error deleting thumbnail %s: %s", filename, err)
		return
	}
	log.Printf("Thumbnail %s deleted", filename)
}

func InitialiseThumbnails() {
	CreateDir(ThumbnailDirectory)
	getThumbnailDirContents()
	populateThumbnails()
	deleteExtraneousThumbnails()
}

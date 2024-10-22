package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
)

func thumbnailAlreadyExists(slug string) bool {
	thumbnailPath := filepath.Join(ThumbnailDirectory, (slug + ".webp"))
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

	defer func() {
		source = nil
		runtime.GC()
	}()

	width, height := 0, 0

	if source.Bounds().Max.X > source.Bounds().Max.Y {
		width = 0
		height = int(ThumbnailMaxPixels)
	} else {
		width = int(ThumbnailMaxPixels)
		height = 0
	}

	thumbnailPath := filepath.Join(ThumbnailDirectory, slug) + ".webp"
	thumbnailImage := imaging.Resize(source, width, height, imaging.Lanczos)

	f, err := os.Create(thumbnailPath)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer f.Close()

	err = webp.Encode(f, thumbnailImage)
	if err != nil {
		log.Printf("Error encoding image: %s", err)
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
	log.Printf("Found %d thumbnails", len(foundThumbnails))
	return foundThumbnails, err
}

func populateThumbnails() {
	for _, row := range GetExistingMetadataFilePaths() {
		slug := row.slug
		filePath := row.filePath
		fileName := row.fileName
		imageFullPath := filepath.Join(filePath, fileName)
		generateThumbnail(imageFullPath, slug)
	}
}

func deleteExtraneousThumbnails() {
	thumbnailDirContents, _ := getThumbnailDirContents()
	for _, thumbnail := range thumbnailDirContents {
		ext := strings.Split(filepath.Ext(thumbnail), ".")[1]
		if ext != "webp" {
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

func GetThumbnailBySlug(slug string) ([]byte, error) {
	thumbnailPath := filepath.Join(ThumbnailDirectory, slug+".webp")

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return nil, err
	}

	thumbnailBlob, err := os.ReadFile(thumbnailPath)

	if err != nil {
		log.Printf("Error reading thumbnail for slug %s: %s", slug, err)
		return nil, err
	}

	return thumbnailBlob, nil
}

func InitialiseThumbnails() {
	CreateDir(ThumbnailDirectory)
	deleteExtraneousThumbnails()
	populateThumbnails()
}

package main

import (
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"photogallery/database"
	"photogallery/logic"
	"photogallery/types"
	"runtime"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

var wgThumbnails sync.WaitGroup

func thumbnailAlreadyExists(slug string) bool {
	thumbnailPath := filepath.Join(logic.ThumbnailDirectory, (slug + ".jpeg"))
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateThumbnail(imageFile string, slug string) {

	defer wgThumbnails.Done()

	if thumbnailAlreadyExists(slug) {
		return
	}

	source, err := imaging.Open(imageFile)
	if err != nil {
		log.Printf("Failed to open image: %v", err)
	}

	defer func() {
		source = nil
		runtime.GC()
	}()

	width, height := 0, 0

	if source.Bounds().Max.X > source.Bounds().Max.Y {
		width = 0
		height = int(logic.ThumbnailMaxPixels)
	} else {
		width = int(logic.ThumbnailMaxPixels)
		height = 0
	}

	thumbnailPath := filepath.Join(logic.ThumbnailDirectory, slug) + ".jpeg"
	thumbnailImage := imaging.Resize(source, width, height, imaging.Lanczos)

	f, err := os.Create(thumbnailPath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	defer f.Close()

	jpeg.Encode(f, thumbnailImage, nil)
	if err != nil {
		log.Printf("Error encoding image: %s", err)
	}

	log.Printf("Thumbnail created for %s: %s", imageFile, thumbnailPath)
}

func getThumbnailDirContents() ([]string, error) {
	var foundThumbnail []string

	err := filepath.Walk(logic.ThumbnailDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error opening Thumbnail directory: %s", err)
			return err
		}
		if !info.IsDir() {
			foundThumbnail = append(foundThumbnail, path)
		}
		return nil
	})
	log.Printf("Found %d thumbnail", len(foundThumbnail))
	return foundThumbnail, err
}

func populateThumbnails() {
	numWorkers := runtime.NumCPU() / 2 // Use half the available CPU cores
	if numWorkers < 1 {
		numWorkers = 1
	}
	workerPool := make(chan struct{}, numWorkers)
	for _, row := range database.GetExistingMetadataFilePaths() {
		workerPool <- struct{}{} // Block if the pool is full
		wgThumbnails.Add(1)
		go func(row types.MetadataFile) {
			defer func() { <-workerPool }()
			slug := row.Slug
			filePath := row.FilePath
			fileName := row.FileName
			imageFullPath := filepath.Join(filePath, fileName)
			generateThumbnail(imageFullPath, slug)
		}(row)
	}

	wgThumbnails.Wait()
}

func deleteExtraneousThumbnails() {
	thumbnailDirContents, _ := getThumbnailDirContents()
	for _, thumbnail := range thumbnailDirContents {
		ext := strings.Split(filepath.Ext(thumbnail), ".")[1]
		if ext != "jpeg" {
			deleteThumbnailByFilename(thumbnail)
		} else {
			slug := strings.TrimSuffix(filepath.Base(thumbnail), filepath.Ext(thumbnail))
			_, err := database.GetMetadataBySlug(slug)
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
	thumbnailPath := filepath.Join(logic.ThumbnailDirectory, slug+".jpeg")
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		log.Printf("Thumbnail file does not exist: %s", thumbnailPath)
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
	logic.CreateDir(logic.ThumbnailDirectory)
	deleteExtraneousThumbnails()
	populateThumbnails()
}

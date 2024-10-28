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

var wgOptimised sync.WaitGroup

func optimisedAlreadyExists(slug string) bool {
	optimisedPath := filepath.Join(logic.OptimisedDirectory, (slug + ".jpeg"))
	if _, err := os.Stat(optimisedPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateOptimised(imageFile string, slug string) {

	defer wgOptimised.Done()

	if optimisedAlreadyExists(slug) {
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
		height = int(logic.OptimisedMaxPixels)
	} else {
		width = int(logic.OptimisedMaxPixels)
		height = 0
	}

	optimisedPath := filepath.Join(logic.OptimisedDirectory, slug) + ".jpeg"
	optimisedImage := imaging.Resize(source, width, height, imaging.Lanczos)

	f, err := os.Create(optimisedPath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	defer f.Close()

	jpeg.Encode(f, optimisedImage, nil)
	if err != nil {
		log.Printf("Error encoding image: %s", err)
	}

	log.Printf("Optimised created for %s: %s", imageFile, optimisedPath)
}

func getOptimisedDirContents() ([]string, error) {
	var foundOptimised []string

	err := filepath.Walk(logic.OptimisedDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error opening Optimised directory: %s", err)
			return err
		}
		if !info.IsDir() {
			foundOptimised = append(foundOptimised, path)
		}
		return nil
	})
	log.Printf("Found %d optimised", len(foundOptimised))
	return foundOptimised, err
}

func populateOptimised() {
	numWorkers := runtime.NumCPU() / 2 // Use half the available CPU cores
	if numWorkers < 1 {
		numWorkers = 1
	}
	workerPool := make(chan struct{}, numWorkers)
	for _, row := range database.GetExistingMetadataFilePaths() {
		workerPool <- struct{}{} // Block if the pool is full
		wgOptimised.Add(1)
		go func(row types.MetadataFile) {
			defer func() { <-workerPool }()
			slug := row.Slug
			filePath := row.FilePath
			fileName := row.FileName
			imageFullPath := filepath.Join(filePath, fileName)
			generateOptimised(imageFullPath, slug)
		}(row)
	}

	wgOptimised.Wait()
}

func deleteExtraneousOptimised() {
	optimisedDirContents, _ := getOptimisedDirContents()
	for _, optimised := range optimisedDirContents {
		ext := strings.Split(filepath.Ext(optimised), ".")[1]
		if ext != "jpeg" {
			deleteOptimisedByFilename(optimised)
		} else {
			slug := strings.TrimSuffix(filepath.Base(optimised), filepath.Ext(optimised))
			_, err := database.GetMetadataBySlug(slug)
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

func GetOptimisedBySlug(slug string) ([]byte, error) {
	optimisedPath := filepath.Join(logic.OptimisedDirectory, slug+".jpeg")
	if _, err := os.Stat(optimisedPath); os.IsNotExist(err) {
		log.Printf("Optimised file does not exist: %s", optimisedPath)
		metadata, _ := database.GetMetadataBySlug(slug)
		filePath, _ := filepath.Abs(filepath.Join(metadata.FilePath, metadata.FileName))
		err = generateOptimisedBeforeReady(filePath, slug)
		if err != nil {
			log.Printf("Error getting optimised for slug %s: %s", slug, err)
			return nil, err
		}
	}
	optimisedBlob, err := os.ReadFile(optimisedPath)
	if err != nil {
		log.Printf("Error reading optimised for slug %s: %s", slug, err)
		return nil, err
	}
	return optimisedBlob, nil
}

func generateOptimisedBeforeReady(imageFile string, slug string) error {

	source, err := imaging.Open(imageFile)
	if err != nil {
		log.Printf("Failed to open image: %v", err)
		return err
	}
	width, height := 0, 0

	if source.Bounds().Max.X > source.Bounds().Max.Y {
		width = 0
		height = int(logic.OptimisedMaxPixels)
	} else {
		width = int(logic.OptimisedMaxPixels)
		height = 0
	}

	optimisedPath := filepath.Join(logic.OptimisedDirectory, slug) + ".jpeg"
	optimisedImage := imaging.Resize(source, width, height, imaging.Lanczos)

	f, err := os.Create(optimisedPath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer f.Close()

	jpeg.Encode(f, optimisedImage, nil)
	if err != nil {
		log.Printf("Error encoding image: %s", err)
		return err
	}

	log.Printf("Optimised created for %s: %s", imageFile, optimisedPath)
	return nil
}

func InitialiseOptimised() {
	logic.CreateDir(logic.OptimisedDirectory)
	deleteExtraneousOptimised()
	go populateOptimised()
}

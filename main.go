package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	LoadEnv()

	DatabaseDirectory, _ = filepath.Abs(os.Getenv("DATA_PATH"))
	ThumbnailDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "thumbnails"))
	OptimisedDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "optimised"))
	ImagePath, _ = filepath.Abs(os.Getenv("IMAGE_PATH"))
	ImageExtensions = strings.Split(os.Getenv("IMAGE_FILES"), ",")
	value, _ := strconv.ParseUint(os.Getenv("THUMBNAIL_MAX_PIXELS"), 10, 64)
	ThumbnailMaxPixels = uint(value)
	value, _ = strconv.ParseUint(os.Getenv("OPTIMISED_MAX_PIXELS"), 10, 64)
	OptimisedMaxPixels = uint(value)

	Database = InitialiseDatabase()
	CreateThumbnailsDir()
	CreateOptimisedDir()
	FoundFiles, _ = GetImageDirContents()
	FoundMetadataFiles = GetExistingMetadataFilePaths()

	log.Printf("Found: %d source images", len(FoundFiles))
	log.Printf(`Found: %d metadata rows`, len(FoundMetadataFiles))

	InitialiseAllMetadata()

	// imageMetadata := GetSourceMetadataForImagePath(FoundFiles[50])
	// InsertMetadataRow(imageMetadata)
	// GenerateThumbnail(FoundFiles[photoInt], imageMetadata.slug)
}

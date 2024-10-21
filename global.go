package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var Database *sql.DB
var DatabaseDirectory string
var ThumbnailDirectory string
var OptimisedDirectory string
var ImagePath string
var ImageExtensions []string
var ImageFormat string
var ThumbnailMaxPixels uint
var OptimisedMaxPixels uint
var FoundMetadataFiles []MetadataFile
var FoundFiles []string
var FoundThumbnails []string
var FoundOptimised []string
var WorkerCount int

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	} else {
		log.Println(".env file loaded")
	}

	DatabaseDirectory, _ = filepath.Abs(os.Getenv("DATA_PATH"))
	ThumbnailDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "thumbnails"))
	OptimisedDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "optimised"))
	ImagePath, _ = filepath.Abs(os.Getenv("IMAGE_PATH"))
	ImageFormat = os.Getenv("IMAGE_FORMAT")
	ImageExtensions = strings.Split(os.Getenv("IMAGE_FILES"), ",")
	u, _ := strconv.ParseUint(os.Getenv("THUMBNAIL_MAX_PIXELS"), 10, 64)
	ThumbnailMaxPixels = uint(u)
	u, _ = strconv.ParseUint(os.Getenv("OPTIMISED_MAX_PIXELS"), 10, 64)
	OptimisedMaxPixels = uint(u)
	WorkerCount, _ = strconv.Atoi(os.Getenv("WORKER_COUNT"))
}

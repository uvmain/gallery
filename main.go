package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var Database *sql.DB
var DatabaseDirectory string
var ThumbnailDirectory string
var ImagePath string
var ImageExtensions []string

func main() {
	loadEnv()

	DatabaseDirectory, _ = filepath.Abs(os.Getenv("DATA_PATH"))
	ThumbnailDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "thumbnails"))
	ImagePath, _ = filepath.Abs(os.Getenv("IMAGE_PATH"))
	ImageExtensions = strings.Split(os.Getenv("IMAGE_FILES"), ",")

	Database = InitialiseDatabase()
	CreateMetadataTable()
	CreateThumbnailsDir()
	GetImageDirContents()

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	} else {
		log.Println(".env file loaded")
	}
}

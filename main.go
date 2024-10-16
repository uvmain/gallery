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

func main() {
	loadEnv()

	databaseDirectory, _ := filepath.Abs(os.Getenv("DATA_PATH"))
	thumbnailPath, _ := filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "thumbnails"))
	imagePath, _ := filepath.Abs(os.Getenv("IMAGE_PATH"))
	imageExtensions := strings.Split(os.Getenv("IMAGE_FILES"), ",")

	initDB(databaseDirectory)
	createThumbnailsDir(thumbnailPath)
	traverseImageDir(imagePath, imageExtensions)

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		log.Println(".env file loaded")
	}
}

func initDB(databaseDirectory string) (*sql.DB, error) {
	dbPath := filepath.Join(databaseDirectory, "sqlite.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Creating database file")
		// Create the database directory
		err := os.MkdirAll(databaseDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating database directory")
		} else {
			log.Println("Database directory created")
		}

		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Error creating database file: %s", err)
			return nil, err
		} else {
			log.Println("Database file created")
		}
		file.Close()
	}

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Fatalf("Error opening database file: %s", err)
		return nil, err
	} else {
		log.Println("Database file opened")
	}

	query := `CREATE TABLE IF NOT EXISTS metadata (
			slug TEXT PRIMARY KEY,
			path TEXT,
			fileName TEXT,
			title TEXT,
			dateTaken DATETIME,
			dateUploaded DATETIME,
			cameraModel TEXT,
			lensModel TEXT,
			aperture TEXT,
			shutterSpeed TEXT,
			flashStatus TEXT,
			focusLength TEXT,
			iso TEXT,
			exposureMode TEXT,
			whiteBalance TEXT
	);`

	_, err = db.Exec(query)

	if err != nil {
		log.Fatalf("Error creating metadata table: %s", err)
		return nil, err
	} else {
		log.Println("Metadata table created")
	}

	return db, err
}

func createThumbnailsDir(dataPath string) error {
	thumbnailsDir := filepath.Join(dataPath, "thumbnails")
	if _, err := os.Stat(thumbnailsDir); os.IsNotExist(err) {
		return os.Mkdir(thumbnailsDir, os.ModePerm)
	}
	return nil
}

func traverseImageDir(imagePath string, imageExtensions []string) ([]string, error) {
	var foundFiles []string

	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, validExt := range imageExtensions {
				if ext == validExt {
					foundFiles = append(foundFiles, path)
					break
				}
			}
		}
		return nil
	})
	return foundFiles, err
}

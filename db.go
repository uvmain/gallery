package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitialiseDatabase() *sql.DB {
	dbPath := filepath.Join(DatabaseDirectory, "sqlite.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Creating database file")
		// Create the database directory
		err := os.MkdirAll(DatabaseDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating database directory: %s", err)
		} else {
			log.Println("Database directory created")
		}

		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Error creating database file: %s", err)
			return nil
		} else {
			log.Println("Database file created")
		}
		file.Close()
	} else {
		log.Println("Database already exists")
	}

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Fatalf("Error opening database file: %s", err)
		return nil
	} else {
		log.Println("Database file opened")
	}

	return db
}

func CreateMetadataTable() {
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

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='metadata'"
	var name string
	checkError := Database.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("Metadata table already exists")
	} else {
		_, err := Database.Query(query)

		if err != nil {
			log.Fatalf("Error creating metadata table: %s", err)
		} else {
			log.Println("Metadata table created")
		}
	}
}

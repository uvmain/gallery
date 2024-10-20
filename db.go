package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type MetadataFile struct {
	slug     string
	filePath string
	fileName string
}

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

	CreateAlbumsTable(db)
	CreateMetadataTable(db)
	Database = db
	return db
}

func CreateAlbumsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS albums (
		name TEXT PRIMARY KEY
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='albums'"

	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("Albums table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating albums table: %s", err)
		} else {
			log.Println("Albums table created")
		}
	}
}

func CreateMetadataTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS metadata (
		slug TEXT PRIMARY KEY,
		filePath TEXT,
		fileName TEXT,
		title TEXT,
		dateTaken DATETIME,
		dateUploaded DATETIME,
		cameraMake TEXT,
		cameraModel TEXT,
		lensMake TEXT,
		lensModel TEXT,
		fStop TEXT,
		shutterSpeed TEXT,
		flashStatus TEXT,
		focalLength TEXT,
		iso TEXT,
		exposureMode TEXT,
		whiteBalance TEXT,
		albums TEXT
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='metadata'"
	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("Metadata table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating metadata table: %s", err)
		} else {
			log.Println("Metadata table created")
		}
	}
}

func GetExistingMetadataFilePaths() []MetadataFile {
	foundMetadataFiles := []MetadataFile{}

	query := "SELECT slug, filePath, fileName FROM metadata"
	rows, err := Database.Query(query)
	if err != nil {
		log.Fatalf("Failed to fetch rows from metadata table: %s", err)
	}

	for rows.Next() {
		var slug string
		var filePath string
		var fileName string
		err = rows.Scan(&slug, &filePath, &fileName)
		if err != nil {
			log.Fatal(err)
		}

		rowResult := MetadataFile{
			slug:     slug,
			filePath: filePath,
			fileName: fileName,
		}

		foundMetadataFiles = append(foundMetadataFiles, rowResult)
	}

	FoundMetadataFiles = foundMetadataFiles
	return foundMetadataFiles
}

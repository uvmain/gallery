package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type MetadataFile struct {
	filePath string
	fileName string
}

var FoundMetadataFiles []MetadataFile = []MetadataFile{}

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

func GetExistingMetadataFilePaths() {
	FoundMetadataFiles = []MetadataFile{}

	query := "SELECT filePath, fileName FROM metadata"
	rows, err := Database.Query(query)
	if err != nil {
		log.Fatalf("Failed to fetch rows from metadata table: %s", err)
	}

	for rows.Next() {
		var filePath string
		var fileName string
		err = rows.Scan(&filePath, &fileName)
		if err != nil {
			log.Fatal(err)
		}

		rowResult := MetadataFile{
			filePath: filePath,
			fileName: fileName,
		}

		FoundMetadataFiles = append(FoundMetadataFiles, rowResult)
	}
}

func InsertMetadataRow(imageMetadata ImageMetadata) error {
	checkQuery := `SELECT COUNT(*) FROM metadata WHERE filePath = ? AND fileName = ?;`
	var count int
	err := Database.QueryRow(checkQuery, imageMetadata.filePath, imageMetadata.fileName).Scan(&count)
	if err != nil {
		log.Printf("error checking existing row for %s: %v", imageMetadata.fileName, err)
		return err
	}

	if count > 0 {
		log.Printf("Metadata row already exists, skipping insert: %s\n", imageMetadata.fileName)
		return nil
	}

	insertQuery := `INSERT INTO metadata (
			slug, filePath, fileName, title, dateTaken, dateUploaded,
			cameraMake, cameraModel, lensMake, lensModel, fStop, shutterSpeed,
			flashStatus, focalLength, iso, exposureMode, whiteBalance, albums
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = Database.Exec(
		insertQuery,
		imageMetadata.slug, imageMetadata.filePath, imageMetadata.fileName,
		imageMetadata.title, imageMetadata.dateTaken, imageMetadata.dateUploaded,
		imageMetadata.cameraMake, imageMetadata.cameraModel, imageMetadata.lensMake,
		imageMetadata.lensModel, imageMetadata.fStop, imageMetadata.shutterSpeed,
		imageMetadata.flashStatus, imageMetadata.focalLength, imageMetadata.iso,
		imageMetadata.exposureMode, imageMetadata.whiteBalance, imageMetadata.albums,
	)
	if err != nil {
		log.Printf("error inserting metadata row: %s", err)
		return err
	}

	log.Printf("Metadata row inserted successfully for %s", imageMetadata.fileName)
	return nil
}

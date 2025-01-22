package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"photogallery/logic"
	"photogallery/types"

	_ "modernc.org/sqlite"
)

var Database *sql.DB

func GetDB() *sql.DB {
	return Database
}

func Initialise() *sql.DB {
	logic.CreateDir(logic.DatabaseDirectory)

	dbPath := filepath.Join(logic.DatabaseDirectory, "sqlite.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Creating database file")

		file, err := os.Create(dbPath)
		if err != nil {
			log.Printf("Error creating database file: %s", err)
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
		log.Printf("Error opening database file: %s", err)
		return nil
	} else {
		log.Println("Database file opened")
	}

	Database = db

	createMetadataTable(db)
	InitialiseAlbums(db)
	InitialiseLinks(db)
	return db
}

func createMetadataTable(db *sql.DB) {
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
		exposureTime TEXT,
		flashStatus TEXT,
		focalLength TEXT,
		iso TEXT,
		exposureMode TEXT,
		whiteBalance TEXT,
		WhiteBalanceMode TEXT
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='metadata'"
	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("metadata table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating metadata table: %s", err)
		} else {
			log.Println("metadata table created")
		}
	}
}

func GetExistingMetadataFilePaths() []types.MetadataFile {
	foundMetadataFiles := []types.MetadataFile{}

	query := "SELECT slug, filePath, fileName FROM metadata"
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Failed to fetch rows from metadata table: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var slug string
		var filePath string
		var fileName string
		err = rows.Scan(&slug, &filePath, &fileName)
		if err != nil {
			log.Println(err)
		}

		rowResult := types.MetadataFile{
			Slug:     slug,
			FilePath: filePath,
			FileName: fileName,
		}

		foundMetadataFiles = append(foundMetadataFiles, rowResult)
	}

	return foundMetadataFiles
}

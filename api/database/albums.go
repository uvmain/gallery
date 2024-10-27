package database

import (
	"database/sql"
	"log"
)

func createAlbumsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS albums (
		name TEXT PRIMARY KEY,
		dateCreated DATETIME,
		coverSlug TEXT
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='albums'"

	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("albums table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating albums table: %s", err)
		} else {
			log.Println("albums table created")
		}
	}
}

func createAlbumLinksTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS album_links (
    albumName TEXT,
    imageSlug TEXT,
    FOREIGN KEY (albumName) REFERENCES albums(name),
    FOREIGN KEY (imageSlug) REFERENCES metadata(slug),
    PRIMARY KEY (albumName, imageSlug)
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='album_links'"

	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("album_links table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating album_links table: %s", err)
		} else {
			log.Println("album_links table created")
		}
	}
}

func InitialiseAlbums(db *sql.DB) {
	createAlbumsTable(db)
	createAlbumLinksTable(db)
}

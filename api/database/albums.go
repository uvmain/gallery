package database

import (
	"database/sql"
	"log"
	"photogallery/logic"
	"photogallery/types"
	"time"
)

func createAlbumsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS albums (
		slug TEXT PRIMARY KEY,
		name TEXT,
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
    albumSlug TEXT,
    imageSlug TEXT,
    FOREIGN KEY (albumSlug) REFERENCES albums(slug),
    FOREIGN KEY (imageSlug) REFERENCES metadata(slug),
    PRIMARY KEY (albumSlug, imageSlug)
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

func GetAllAlbums() []types.Album {
	var albums []types.Album

	query := `SELECT slug, name, dateCreated, coverSlug FROM albums ORDER BY datecreated DESC;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []types.Album{}
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		var name string
		var dateCreated string
		var coverSlug string
		err = rows.Scan(&slug, &name, &dateCreated, &coverSlug)
		if err != nil {
			log.Fatal(err)
		}

		rowResult := types.Album{
			Slug:        slug,
			Name:        name,
			DateCreated: dateCreated,
			CoverSlug:   coverSlug,
		}

		albums = append(albums, rowResult)
	}
	return albums
}

func InsertAlbumRow(album types.Album) error {
	stmt, err := Database.Prepare(`INSERT INTO albums (
		slug, name, dateCreated, coverSlug
	) VALUES (?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		logic.GenerateSlug(), album.Name, time.Now(), album.CoverSlug,
	)
	if err != nil {
		log.Printf("error inserting album row: %s", err)
		return err
	}

	log.Printf("Album row inserted successfully for %s", album.Name)
	return nil
}

func InitialiseAlbums(db *sql.DB) {
	createAlbumsTable(db)
	createAlbumLinksTable(db)
}

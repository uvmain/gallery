package database

import (
	"database/sql"
	"log"
	"photogallery/types"
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

func GetAllAlbums() []types.Album {
	var albums []types.Album

	query := `SELECT name, dateCreated, coverSlug FROM albums ORDER BY datecreated DESC;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []types.Album{}
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var dateCreated string
		var coverSlug string
		err = rows.Scan(&name, &dateCreated, &coverSlug)
		if err != nil {
			log.Fatal(err)
		}

		rowResult := types.Album{
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
		name, dateCreated, coverSlug
	) VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		album.Name, album.DateCreated, album.CoverSlug,
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

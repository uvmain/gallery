package database

import (
	"database/sql"
	"gallery/core/logic"
	"gallery/core/types"
	"log"
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

func GetAlbum(slug string) (types.Album, error) {
	var album types.Album
	query := `SELECT slug, name, dateCreated, coverSlug FROM albums where slug = ?;`
	err := Database.QueryRow(query, slug).Scan(
		&album.Slug,
		&album.Name,
		&album.DateCreated,
		&album.CoverSlug,
	)
	if err != nil {
		return types.Album{}, err
	}
	return album, nil
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
			log.Println(err)
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

func DeleteAlbumRow(albumSlug string) error {
	stmt, err := Database.Prepare(`DELETE FROM albums where slug = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(albumSlug)
	if err != nil {
		log.Printf("error deleting album row: %s", err)
		return err
	}

	log.Printf("Album row deleted successfully for albumSlug %s", albumSlug)
	return nil
}

func UpdateAlbumCover(albumSlug string, coverSlug string) error {
	stmt, err := Database.Prepare(`UPDATE albums set coverSlug = ? where slug = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(coverSlug, albumSlug)
	if err != nil {
		log.Printf("error updating cover for album row: %s", err)
		return err
	}

	log.Printf("Cover updated to %s successfully for albumSlug %s", coverSlug, albumSlug)
	return nil
}

func UpdateAlbumName(albumSlug string, albumName string) error {
	stmt, err := Database.Prepare(`UPDATE albums set name = ? where slug = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(albumName, albumSlug)
	if err != nil {
		log.Printf("error updating name for album row: %s", err)
		return err
	}

	log.Printf("Album Name updated to %s successfully for albumSlug %s", albumName, albumSlug)
	return nil
}

func InitialiseAlbums(db *sql.DB) {
	createAlbumsTable(db)
}

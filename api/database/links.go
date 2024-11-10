package database

import (
	"database/sql"
	"log"
	"photogallery/types"
)

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

func GetAlbumLinks(slug string) ([]string, error) {
	var links []string
	query := `SELECT imageSlug FROM album_links where albumSlug = ?;`
	rows, err := Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var imageSlug string
		err = rows.Scan(&imageSlug)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, imageSlug)
	}
	return links, nil
}

func GetImageLinks(slug string) ([]string, error) {
	var links []string
	query := `SELECT albumSlug FROM album_links where imageSlug = ?;`
	rows, err := Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var albumSlug string
		err = rows.Scan(&albumSlug)
		if err != nil {
			log.Fatal(err)
		}

		links = append(links, albumSlug)
	}
	return links, nil
}

func GetAllLinks() []types.Link {
	var links []types.Link

	query := `SELECT albumSlug, imageSlug FROM album_links;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []types.Link{}
	}
	defer rows.Close()

	for rows.Next() {
		var albumSlug string
		var imageSlug string
		err = rows.Scan(&albumSlug, &imageSlug)
		if err != nil {
			log.Fatal(err)
		}

		rowResult := types.Link{
			AlbumSlug: albumSlug,
			ImageSlug: imageSlug,
		}

		links = append(links, rowResult)
	}
	return links
}

func InsertLinkRow(link types.Link) error {
	stmt, err := Database.Prepare(`INSERT INTO albums (
		albumSlug, imageSlug
	) VALUES (?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		link.AlbumSlug, link.ImageSlug,
	)
	if err != nil {
		log.Printf("error inserting album link row: %s", err)
		return err
	}

	log.Printf("Album link row inserted successfully for %s %s", link.AlbumSlug, link.ImageSlug)
	return nil
}

func DeleteLinkRow(link types.Link) error {
	stmt, err := Database.Prepare(`DELETE FROM albums where albumSlug = ? and imageSlug = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(link.AlbumSlug, link.ImageSlug)
	if err != nil {
		log.Printf("error deleting album row: %s", err)
		return err
	}

	log.Printf("Album row deleted successfully for %s %s", link.AlbumSlug, link.ImageSlug)
	return nil
}

func InitialiseLinks(db *sql.DB) {
	createAlbumLinksTable(db)
}

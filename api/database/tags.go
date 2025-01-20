package database

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"photogallery/types"
	"slices"
	"strconv"
	"strings"
)

func createTagsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS tags (
		tag TEXT,
		imageSlug TEXT,
		FOREIGN KEY (imageSlug) REFERENCES metadata(slug),
		PRIMARY KEY (tag, imageSlug)
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='tags'"

	var name string
	checkError := db.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("tags table already exists")
	} else {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating tags table: %s", err)
		} else {
			log.Println("tags table created")
		}
	}
}

func GetAllTags() ([]string, error) {
	var tags []string
	query := `SELECT DISTINCT tag FROM tags;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal(err)
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

func GetTagsForSlug(slug string) ([]string, error) {
	var tags []string
	query := `SELECT tag FROM tags where imageSlug = ?`
	rows, err := Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func GetSlugsForTag(tag string) ([]string, error) {
	var slugs []string
	query := `SELECT tag FROM tags where tag = ?`
	rows, err := Database.Query(query, tag)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			log.Fatal(err)
		}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

func InsertTagsRow(tag types.Tag) error {
	stmt, err := Database.Prepare(`INSERT INTO tags (tag, imageSlug) VALUES (?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.Tag, tag.ImageSlug)
	if err != nil {
		log.Printf("error inserting tag row: %s", err)
		return err
	}

	log.Printf("Tag row inserted successfully for %s %s", tag.Tag, tag.ImageSlug)
	return nil
}

func DeleteTagsRow(tag types.Tag) error {
	stmt, err := Database.Prepare(`DELETE FROM tags WHERE tag = ? AND imageSlug = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.Tag, tag.ImageSlug)
	if err != nil {
		log.Printf("error deleting tag row: %s", err)
		return err
	}

	log.Printf("Tag row deleted successfully for %s %s", tag.Tag, tag.ImageSlug)
	return nil
}

func CreateTagsOnUpload(tags types.TagsUpload) error {
	log.Printf("adding tags for %s", tags.ImageSlug)
	checkQuery := `SELECT title, cameraMake, cameraModel, lensMake, lensModel, fStop, flashStatus, focalLength, iso, exposureMode FROM metadata where slug = ?;`
	var title string
	var cameraMake string
	var cameraModel string
	var lensMake string
	var lensModel string
	var fStop string
	var flashStatus string
	var focalLength string
	var iso string
	var exposureMode string
	err := Database.QueryRow(checkQuery, tags.ImageSlug).Scan(
		&title,
		&cameraMake,
		&cameraModel,
		&lensMake,
		&lensModel,
		&fStop,
		&flashStatus,
		&focalLength,
		&iso,
		&exposureMode,
	)

	if err != nil {
		log.Printf("Error creating initial tags: %s", err)
		return err
	}

	var newTags []string
	newTags = append(newTags, tags.Tags...)
	newTags = append(newTags, strings.Split(title, " ")...)
	newTags = append(newTags, cameraMake)
	newTags = append(newTags, cameraModel)
	newTags = append(newTags, lensMake)
	newTags = append(newTags, lensModel)

	fStop = getFStopOrFocalLength(fStop, "fStop")
	newTags = append(newTags, fStop)

	if strings.Contains(flashStatus, "Fired") {
		newTags = append(newTags, "flash:fired")
	}
	focalLength = getFStopOrFocalLength(focalLength, "focalLength")
	newTags = append(newTags, focalLength)

	if iso != "unknown" {
		iso = fmt.Sprintf("iso:%s", iso)
		newTags = append(newTags, iso)
	}

	if exposureMode != "unknown" {
		newTags = append(newTags, exposureMode)
	}

	slices.Sort(newTags)
	newTags = slices.Compact(newTags)

	for _, tag := range newTags {
		if len(tag) > 2 {
			var newTag = types.Tag{
				Tag:       strings.ToLower(tag),
				ImageSlug: tags.ImageSlug,
			}
			InsertTagsRow(newTag)
		}
	}
	return nil
}

func InitialiseTags(db *sql.DB) {
	createTagsTable(db)
}

func getFStopOrFocalLength(value string, kind string) string {
	var first, second float64

	if value != "unknown" {
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			var err error
			first, err = strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return ""
			}
			second, err = strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return ""
			}
		}
	} else {
		first, second = 1, 1
	}

	result := first / second
	if math.IsNaN(result) {
		return ""
	}
	if kind == "fStop" {
		return fmt.Sprintf("ƒ/%.1f", result)
	}
	if kind == "focalLength" {
		// Round the result and return as an integer
		return fmt.Sprintf("%dmm", int(result))
	}
	return ""
}

package database

import (
	"fmt"
	"log"
	"math"
	"photogallery/logic"
	"photogallery/types"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func InitialiseTags() {
	createTagsTable()
	populateTags()
}

func createTagsTable() {
	query := `CREATE TABLE IF NOT EXISTS tags (
		tag TEXT,
		imageSlug TEXT,
		FOREIGN KEY (imageSlug) REFERENCES metadata(slug),
		PRIMARY KEY (tag, imageSlug)
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='tags'"

	var name string
	checkError := Database.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("tags table already exists")
	} else {
		_, err := Database.Exec(query)
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
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Println(err)
		}

		tags = append(tags, strings.ToLower(tag))
	}

	dimensionTags := []string{"landscape", "portrait", "square", "panoramic"}
	tags = append(tags, dimensionTags...)

	titles := getAllTitleTags()
	tags = append(tags, titles...)

	albums := getAllAlbumTags()
	tags = append(tags, albums...)

	tags = logic.StringArraySortUnique(tags)

	return tags, nil
}

func getAllTitleTags() []string {
	checkQuery := `SELECT DISTINCT title FROM metadata;`
	var titles []string
	rows, _ := Database.Query(checkQuery)
	defer rows.Close()

	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			log.Println(err)
		} else {
			titleRegexp := regexp.MustCompile(`[ \-_]+`) // Matches [" ", "-", "_"]
			titleArray := titleRegexp.Split(title, -1)
			for _, titleTag := range titleArray {
				if len(titleTag) > 2 {
					titles = append(titles, strings.ToLower(titleTag))
				}
			}
		}
	}
	titles = logic.StringArraySortUnique(titles)
	return titles
}

func getAllAlbumTags() []string {
	checkQuery := `SELECT DISTINCT name FROM albums;`
	var names []string
	rows, _ := Database.Query(checkQuery)
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Println(err)
		} else {
			nameRegexp := regexp.MustCompile(`[ \-_]+`) // Matches [" ", "-", "_"]
			nameArray := nameRegexp.Split(name, -1)
			for _, nameTag := range nameArray {
				if len(nameTag) > 2 {
					names = append(names, strings.ToLower(nameTag))
				}
			}
		}
	}
	names = logic.StringArraySortUnique(names)
	return names
}

func GetTagsForSlug(slug string) ([]string, error) {
	var tags []string

	// image tags|metadata
	query := `SELECT tag FROM tags where imageSlug = ?;`
	rows, err := Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Println(err)
		}
		tags = append(tags, tag)
	}

	// image title
	checkQuery := `SELECT title FROM metadata WHERE slug = ?;`
	var title string
	err = Database.QueryRow(checkQuery, slug).Scan(&title)
	if err != nil {
		log.Printf("Query failed: %v", err)
	} else {
		titleRegexp := regexp.MustCompile(`[ \-_]+`) // Matches [" ", "-", "_"]
		titleArray := titleRegexp.Split(title, -1)
		tags = append(tags, titleArray...)
	}

	// albums
	query = "select name from albums join album_links on albums.slug == album_links.albumSlug where imageSlug = ?"
	var albumTitles []string
	rows, err = Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var albumTitle string
		err = rows.Scan(&albumTitle)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
		}
		albumTitle = strings.ToLower(albumTitle)
		albumTitleSplit := strings.Split(albumTitle, " ")

		albumTitles = append(albumTitles, albumTitleSplit...)
	}
	tags = append(tags, albumTitles...)

	// dimensions
	query = "select orientation, panoramic from dimensions where imageSlug = ?"
	rows, err = Database.Query(query, slug)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var orientation string
		var panoramic bool
		err = rows.Scan(&orientation, &panoramic)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
		}
		tags = append(tags, orientation)
		if panoramic {
			tags = append(tags, "panoramic")
		}
	}

	foundTags := []string{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if len(tag) > 2 && tag != "unknown" {
			foundTags = append(foundTags, tag)
		}
	}

	foundTags = logic.StringArraySortUnique(foundTags)

	return foundTags, nil
}

func GetSlugsForTag(tag string) ([]string, error) {
	var slugs []string

	// tags|metadata
	query := `SELECT imageSlug FROM tags where tag = ?;`
	rows, err := Database.Query(query, tag)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			log.Printf("Query failed: %v", err)
		}
		slugs = append(slugs, slug)
	}

	// image titles
	likeQuery := `SELECT slug FROM metadata where lower(title) like lower(?);`
	likePattern := fmt.Sprintf("%%%s%%", tag) // Add % wildcards around tag
	rows, err = Database.Query(likeQuery, likePattern)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
		}
		slugs = append(slugs, slug)
	}

	// albums
	likeQuery = `SELECT slug FROM albums where lower(name) like lower(?);`
	likePattern = fmt.Sprintf("%%%s%%", tag) // Add % wildcards around tag
	rows, err = Database.Query(likeQuery, likePattern)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var albumSlug string
		err = rows.Scan(&albumSlug)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
		}
		albumSlugs, _ := GetAlbumLinks(albumSlug)
		slugs = append(slugs, albumSlugs...)
	}

	// dimensions
	query = `SELECT imageSlug FROM dimensions where orientation = ?;`
	rows, err = Database.Query(query, tag)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		err = rows.Scan(&slug)
		if err != nil {
			log.Println(err)
		}
		slugs = append(slugs, slug)
	}

	// panoramic
	if tag == "panoramic" {
		query = `SELECT imageSlug FROM dimensions where panoramic = 1;`
		rows, err = Database.Query(query)
		if err != nil {
			log.Printf("Query failed: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var slug string
			err = rows.Scan(&slug)
			if err != nil {
				log.Println(err)
			}
			slugs = append(slugs, slug)
		}
	}

	slugs = logic.StringArraySortUnique(slugs)
	return slugs, nil
}

func InsertTagsRow(tag types.Tag) error {
	stmt, err := Database.Prepare(`INSERT INTO tags (tag, imageSlug) VALUES (?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(strings.ToLower(tag.Tag), tag.ImageSlug)
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
	checkQuery := `SELECT cameraMake, cameraModel, lensMake, lensModel, fStop, flashStatus, focalLength, iso, exposureMode FROM metadata where slug = ?;`
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
		&cameraMake, &cameraModel, &lensMake, &lensModel,
		&fStop, &flashStatus, &focalLength, &iso, &exposureMode,
	)

	if err != nil {
		log.Printf("Error creating initial tags: %s", err)
		return err
	}

	var newTags []string
	newTags = append(newTags, tags.Tags...)

	cameraMake = logic.TernaryString(cameraMake == "none", "", cameraMake)
	cameraModel = logic.TernaryString(cameraModel == "none" || len(cameraModel) < 4, "", cameraModel)
	lensMake = logic.TernaryString(lensMake == "none", "", lensMake)
	lensModel = logic.TernaryString(lensModel == "none" || len(lensMake) < 4, "", lensModel)

	newTags = append(newTags, cameraMake, cameraModel, lensMake, lensModel)

	fStop = getFStopOrFocalLength(fStop, "fStop")
	newTags = append(newTags, fStop)

	if strings.Contains(flashStatus, "Fired") {
		newTags = append(newTags, "flash fired")
	}
	focalLength = getFStopOrFocalLength(focalLength, "focalLength")
	newTags = append(newTags, focalLength)

	if iso != "unknown" {
		iso = fmt.Sprintf("iso %s", iso)
		newTags = append(newTags, iso)
	}

	if exposureMode != "unknown" {
		newTags = append(newTags, exposureMode)
	}

	newTags = logic.StringArraySortUnique(newTags)

	for _, tag := range newTags {
		tag = strings.TrimSpace(tag)
		if len(tag) > 2 && tag != "unknown" {
			var newTag = types.Tag{
				Tag:       strings.ToLower(tag),
				ImageSlug: tags.ImageSlug,
			}
			InsertTagsRow(newTag)
		}
	}
	return nil
}

func populateTags() {
	slugs, err := GetAllSlugs()

	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}

	existingSlugs := GetTaggedSlugs()

	slugsToInsert := []string{}
	for _, slug := range slugs {
		if !slices.Contains(existingSlugs, slug) {
			slugsToInsert = append(slugsToInsert, slug)
		}
	}

	for _, slug := range slugsToInsert {
		var newTag = types.TagsUpload{
			Tags:      []string{},
			ImageSlug: slug,
		}
		CreateTagsOnUpload(newTag)
	}
}

func GetTaggedSlugs() []string {
	var slugs []string
	query := `SELECT DISTINCT imageSlug FROM tags;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var imageSlug string
		err = rows.Scan(&imageSlug)
		if err != nil {
			log.Println(err)
		}

		slugs = append(slugs, imageSlug)
	}
	return slugs
}

func getFStopOrFocalLength(value string, kind string) string {
	if value == "unknown" {
		return ""
	}

	var first, second float64

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

	result := first / second
	if math.IsNaN(result) {
		return ""
	}
	if kind == "fStop" {
		return logic.TernaryString(int(result) > 0, fmt.Sprintf("f%.1f", result), "")
	}
	if kind == "focalLength" {
		return logic.TernaryString(int(result) > 0, fmt.Sprintf("%dmm", int(result)), "")
	}
	return ""
}

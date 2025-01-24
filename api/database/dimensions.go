package database

import (
	"log"
	"path/filepath"
	"photogallery/types"
	"runtime"
	"slices"

	"github.com/disintegration/imaging"
)

func InitialiseDimensions() {
	createDimensionsTable()
	populateDimensions()
}

func GetSourceDimensionsForSlug(slug string) (types.DimensionsRow, error) {
	metadata, err := GetMetadataBySlug(slug)
	if err != nil {
		log.Printf("Error getting metadata: %s", err)
		return types.DimensionsRow{}, err
	}

	imagePath := filepath.Join(metadata.FilePath, metadata.FileName)

	source, err := imaging.Open(imagePath)
	if err != nil {
		log.Printf("Failed to open image: %v", err)
	}

	defer func() {
		source = nil
		runtime.GC()
	}()

	width, height := source.Bounds().Max.X, source.Bounds().Max.Y

	var orientation string
	if float64(source.Bounds().Max.X) > float64(source.Bounds().Max.Y)*1.05 {
		orientation = "landscape"
	} else if float64(source.Bounds().Max.Y) > float64(source.Bounds().Max.X)*1.05 {
		orientation = "portrait"
	} else {
		orientation = "square"
	}

	panoramic := false
	if width > 0 && height > 0 && (width >= (height*2) || height >= (width*2)) {
		panoramic = true
	}

	dimensions := types.DimensionsRow{
		ImageSlug:   slug,
		Width:       width,
		Height:      height,
		Orientation: orientation,
		Panoramic:   panoramic,
	}

	return dimensions, nil
}

func createDimensionsTable() {
	query := `CREATE TABLE IF NOT EXISTS dimensions (
		imageSlug TEXT,
		width INTEGER,
		height INTEGER,
		orientation TEXT,
		panoramic TEXT,
		FOREIGN KEY (imageSlug) REFERENCES metadata(slug),
		PRIMARY KEY (imageSlug)
	);`

	checkQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='dimensions'"

	var name string
	checkError := Database.QueryRow(checkQuery).Scan(&name)

	if checkError == nil {
		log.Println("dimensions table already exists")
	} else {
		_, err := Database.Exec(query)
		if err != nil {
			log.Printf("Error creating dimensions table: %s", err)
		} else {
			log.Println("dimensions table created")
		}
	}
}

func GetDimensionedSlugs() ([]string, error) {
	var slugs []string
	query := `SELECT DISTINCT imageSlug FROM dimensions;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []string{}, err
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
	return slugs, nil
}

func CreateDimsensionsOnUpload(slug string) {
	dimensions, err := GetSourceDimensionsForSlug(slug)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	InsertDimensionsRow(dimensions)
}

func populateDimensions() {
	slugs, err := GetAllSlugs()

	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}

	existingSlugs, err := GetDimensionedSlugs()

	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}

	slugsToInsert := []string{}
	for _, slug := range slugs {
		if !slices.Contains(existingSlugs, slug) {
			slugsToInsert = append(slugsToInsert, slug)
		}
	}

	for _, slug := range slugsToInsert {
		dimensions, err := GetSourceDimensionsForSlug(slug)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		InsertDimensionsRow(dimensions)
	}
}

func InsertDimensionsRow(dimensions types.DimensionsRow) error {
	log.Printf("adding dimensions for %s", dimensions.ImageSlug)
	stmt, err := Database.Prepare(`INSERT INTO dimensions (imageSlug, width, height, orientation, panoramic) VALUES (?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dimensions.ImageSlug, dimensions.Width, dimensions.Height, dimensions.Orientation, dimensions.Panoramic)
	if err != nil {
		log.Printf("error inserting tag row: %s", err)
		return err
	}

	log.Printf("Dimensions row inserted successfully for %s", dimensions.ImageSlug)
	return nil
}

func GetAllDimensions() ([]types.DimensionsRow, error) {
	var dimensions []types.DimensionsRow
	query := `SELECT imageSlug, width, height, orientation, panoramic FROM dimensions;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []types.DimensionsRow{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dimension types.DimensionsRow
		err = rows.Scan(&dimension.ImageSlug, &dimension.Width, &dimension.Height, &dimension.Orientation, &dimension.Panoramic)
		if err != nil {
			log.Println(err)
		}
		dimensions = append(dimensions, dimension)
	}
	return dimensions, nil
}

func GetDimensionForSlug(slug string) (types.DimensionsRow, error) {
	var dimension types.DimensionsRow
	query := `SELECT imageSlug, width, height, orientation, panoramic FROM dimensions where imageSlug = ?;`
	err := Database.QueryRow(query, slug).Scan(
		&dimension.ImageSlug,
		&dimension.Width,
		&dimension.Height,
		&dimension.Orientation,
		&dimension.Panoramic,
	)
	if err != nil {
		return types.DimensionsRow{}, err
	}
	return dimension, nil
}

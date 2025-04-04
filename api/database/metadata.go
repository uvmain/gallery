package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"photogallery/exif"
	"photogallery/logic"
	"photogallery/types"

	_ "modernc.org/sqlite"
)

func deleteMetadataRowByFile(filePath string, fileName string) error {

	deleteStatement := `DELETE FROM metadata where filePath = ? AND fileName = ?;`

	fullFilePath := filepath.Join(filePath, fileName)

	_, err := Database.Exec(deleteStatement, filePath, fileName)
	if err != nil {
		log.Printf("Error deleting metadata for %s: %s", fullFilePath, err)
		return err
	}

	log.Printf("Metadata row deleted successfully for %s", fullFilePath)
	return nil
}

func getMetadataRowsToDelete() []types.MetadataFile {
	results := []types.MetadataFile{}

	filesMap := make(map[string]bool)
	foundFiles, _ := logic.GetDirContents(logic.ImagePath)
	for _, v := range foundFiles {
		filesMap[v] = true
	}

	foundMetadataFiles := GetExistingMetadataFilePaths()
	for _, v := range foundMetadataFiles {
		fullFilePath := filepath.Join(v.FilePath, v.FileName)
		if !filesMap[fullFilePath] {
			result := types.MetadataFile{
				Slug:     v.Slug,
				FilePath: v.FilePath,
				FileName: v.FileName,
			}
			results = append(results, result)
		}
	}

	return results
}

func deleteExtraneousMetadata() {
	metadataToDelete := getMetadataRowsToDelete()

	for _, file := range metadataToDelete {
		filePath := file.FilePath
		fileName := file.FileName
		deleteMetadataRowByFile(filePath, fileName)
	}
}

func InitialiseMetadata() {
	GetExistingMetadataFilePaths()
	populateMetadata()
	deleteExtraneousMetadata()
}

func insertMetadataRow(imageMetadata types.ImageMetadata) error {

	parsedDateTaken, err := logic.FormatTimeToString(imageMetadata.DateTaken.String())
	if err != nil {
		return err
	}

	stmt, err := Database.Prepare(`INSERT INTO metadata (
		slug, filePath, fileName, title, dateTaken, dateUploaded,
		cameraMake, cameraModel, lensMake, lensModel, fStop, exposureTime,
		flashStatus, focalLength, iso, exposureMode, whiteBalance, WhiteBalanceMode
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		imageMetadata.Slug, imageMetadata.FilePath, imageMetadata.FileName,
		imageMetadata.Title, parsedDateTaken, imageMetadata.DateUploaded,
		imageMetadata.CameraMake, imageMetadata.CameraModel, imageMetadata.LensMake,
		imageMetadata.LensModel, imageMetadata.FStop, imageMetadata.ExposureTime,
		imageMetadata.FlashStatus, imageMetadata.FocalLength, imageMetadata.ISO,
		imageMetadata.ExposureMode, imageMetadata.WhiteBalance, imageMetadata.WhiteBalanceMode,
	)
	if err != nil {
		log.Printf("error inserting metadata row: %s", err)
		return err
	}

	log.Printf("Metadata row inserted successfully for %s", imageMetadata.FileName)
	return nil
}

func populateMetadata() {
	foundFiles, _ := logic.GetDirContents(logic.ImagePath)
	for _, file := range foundFiles {
		checkQuery := `SELECT COUNT(*) FROM metadata WHERE filePath = ? AND fileName = ?;`
		filePath := filepath.Dir(file)
		fileName := filepath.Base(file)
		var count int
		err := Database.QueryRow(checkQuery, filePath, fileName).Scan(&count)
		if err != nil {
			log.Printf("error checking existing row for %s: %v", fileName, err)
		} else if count > 0 {
			log.Printf("Metadata row already exists, skipping insert: %s\n", fileName)
		} else {
			imageMetadata := exif.GetSourceMetadataForImagePath(file)
			insertMetadataRow(imageMetadata)
		}
	}
}

func GetMetadataBySlug(slug string) (types.ImageMetadataWithDimensions, error) {
	var row types.ImageMetadataWithDimensions
	query := `SELECT slug, filePath, fileName, title, dateTaken, dateUploaded, cameraMake, cameraModel, lensMake, lensModel, fStop, exposureTime, flashStatus, focalLength, iso, exposureMode, whiteBalance, whiteBalanceMode FROM metadata WHERE slug = ?;`

	err := Database.QueryRow(query, slug).Scan(
		&row.Slug,
		&row.FilePath,
		&row.FileName,
		&row.Title,
		&row.DateTaken,
		&row.DateUploaded,
		&row.CameraMake,
		&row.CameraModel,
		&row.LensMake,
		&row.LensModel,
		&row.FStop,
		&row.ExposureTime,
		&row.FlashStatus,
		&row.FocalLength,
		&row.ISO,
		&row.ExposureMode,
		&row.WhiteBalance,
		&row.WhiteBalanceMode,
	)
	if err != nil {
		return types.ImageMetadataWithDimensions{}, err
	}

	dimensions, err := GetDimensionForSlug(slug)
	if err == nil {
		row.Width = dimensions.Width
		row.Height = dimensions.Height
		row.Orientation = dimensions.Orientation
		row.Panoramic = dimensions.Panoramic
	}

	return row, nil
}

func CheckMetadataByFileNameExists(filename string) bool {
	query := `SELECT slug FROM metadata WHERE fileName = ?;`
	var slug string

	err := Database.QueryRow(query, filename).Scan(&slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Println("Database error:", err)
		return false
	}
	return true
}

func GetSlugsWithDimensions() ([]types.SlugWithDimensions, error) {
	var slugsWithDimensions []types.SlugWithDimensions
	slugs, err := GetSlugsOrderedByDateTaken()
	if err != nil {
		log.Printf("Failed to get slugs: %v", err)
		return []types.SlugWithDimensions{}, err
	}

	for _, slug := range slugs {
		dimension, err := GetDimensionForSlug(slug)
		if err != nil {
			log.Printf("Failed to get dimensions: %v", err)
		} else {
			slugWithDimensions := types.SlugWithDimensions{
				Slug:   slug,
				Width:  dimension.Width,
				Height: dimension.Height,
			}
			slugsWithDimensions = append(slugsWithDimensions, slugWithDimensions)
		}
	}

	return slugsWithDimensions, nil
}

func GetSlugsOrderedByDateTaken() ([]string, error) {
	var slugs []string

	query := `SELECT slug FROM metadata order by dateTaken desc`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		slugs = append(slugs, slug)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return slugs, nil
}

func GetSlugsOrderedRandom() ([]string, error) {
	var slugs []string

	query := `SELECT slug FROM metadata ORDER BY RANDOM() DESC;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		slugs = append(slugs, slug)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return slugs, nil
}

func GetOriginalImageBlobBySlug(slug string) ([]byte, error) {
	metadata, _ := GetMetadataBySlug(slug)
	filePath, _ := filepath.Abs(filepath.Join(metadata.FilePath, metadata.FileName))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Original file does not exist: %s:  %s", filePath, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading original image for slug %s: %s", slug, err)
		return nil, err
	}
	return blob, nil
}

func UpdateMetadataBySlug(slug string, updates map[string]interface{}) error {
	query := "UPDATE metadata SET "
	params := []interface{}{}
	i := 1
	for field, value := range updates {
		query += fmt.Sprintf("%s = ?", field)
		if i < len(updates) {
			query += ", "
		}
		params = append(params, value)
		i++
	}
	query += " WHERE slug = ?"
	params = append(params, slug)

	stmt, err := Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)

	if err == nil {
		log.Printf("Metadata updated for %s, %s", slug, updates)
	}
	return err
}

func DeleteMetadataBySlug(slug string) error {
	query := "DELETE from metadata WHERE slug = ?"

	stmt, err := Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(slug)

	if err == nil {
		log.Printf("Metadata deleted for %s", slug)
	}
	return err
}

func PopulateMetadataForUpload(fileName string) (string, error) {
	filePath := filepath.Join(logic.ImageDirectory, fileName)

	checkQuery := `SELECT COUNT(*) FROM metadata WHERE filePath = ? AND fileName = ?;`

	var count int
	err := Database.QueryRow(checkQuery, filePath, fileName).Scan(&count)
	if err != nil {
		log.Printf("error checking existing row for %s: %v", fileName, err)
		return "", err
	} else if count > 0 {
		log.Printf("Metadata row already exists, skipping insert: %s\n", fileName)
		return "", errors.New("metadata already exists")
	} else {
		imageMetadata := exif.GetSourceMetadataForImagePath(filePath)
		insertMetadataRow(imageMetadata)
		return imageMetadata.Slug, nil
	}
}

func GetAllSlugs() ([]string, error) {
	var slugs []string

	query := `SELECT slug FROM metadata;`
	rows, err := Database.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		slugs = append(slugs, slug)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return slugs, nil
}

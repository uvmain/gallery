package main

import (
	"log"
	"path/filepath"

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

func getMetadataRowsToDelete() []MetadataFile {
	results := []MetadataFile{}

	filesMap := make(map[string]bool)
	for _, v := range FoundFiles {
		filesMap[v] = true
	}

	for _, v := range FoundMetadataFiles {
		fullFilePath := filepath.Join(v.filePath, v.fileName)
		if !filesMap[fullFilePath] {
			result := MetadataFile{
				slug:     v.slug,
				filePath: v.filePath,
				fileName: v.fileName,
			}
			results = append(results, result)
		}
	}

	return results
}

func deleteExtraneousMetadata() {
	metadataToDelete := getMetadataRowsToDelete()

	for _, file := range metadataToDelete {
		filePath := file.filePath
		fileName := file.fileName
		deleteMetadataRowByFile(filePath, fileName)
	}
}

func InitialiseMetadata() {
	GetExistingMetadataFilePaths()
	populateMetadata()
	deleteExtraneousMetadata()
}

func insertMetadataRow(imageMetadata ImageMetadata) error {

	insertQuery := `INSERT INTO metadata (
			slug, filePath, fileName, title, dateTaken, dateUploaded,
			cameraMake, cameraModel, lensMake, lensModel, fStop, shutterSpeed,
			flashStatus, focalLength, iso, exposureMode, whiteBalance, albums
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := Database.Exec(
		insertQuery,
		imageMetadata.slug, imageMetadata.filePath, imageMetadata.fileName,
		imageMetadata.title, imageMetadata.dateTaken, imageMetadata.dateUploaded,
		imageMetadata.cameraMake, imageMetadata.cameraModel, imageMetadata.lensMake,
		imageMetadata.lensModel, imageMetadata.fStop, imageMetadata.shutterSpeed,
		imageMetadata.flashStatus, imageMetadata.focalLength, imageMetadata.iso,
		imageMetadata.exposureMode, imageMetadata.whiteBalance, imageMetadata.albums,
	)
	if err != nil {
		log.Printf("error inserting metadata row: %s", err)
		return err
	}

	log.Printf("Metadata row inserted successfully for %s", imageMetadata.fileName)
	return nil
}

func populateMetadata() {
	for _, file := range FoundFiles {
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
			imageMetadata := GetSourceMetadataForImagePath(file)
			insertMetadataRow(imageMetadata)
		}
	}
}

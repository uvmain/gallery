package main

import (
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func DeleteMetadataRowByFile(filePath string, fileName string) error {

	deleteStatement := `DELETE FROM metadata where filePath = ? AND fileName = ?;`

	_, err := Database.Exec(deleteStatement, filePath, fileName)
	if err != nil {
		log.Printf("error deleting metadata row: %s", err)
		return err
	}

	log.Printf("Metadata row deleted successfully for %s", fileName)
	return nil
}

func DeleteMetadataRowBySlug(slug string) error {

	deleteStatement := `DELETE FROM metadata where slug = ?;`

	_, err := Database.Exec(deleteStatement, slug)
	if err != nil {
		log.Printf("error deleting metadata row: %s", err)
		return err
	}

	log.Printf("Metadata row deleted successfully for slug %s", slug)
	return nil
}

func GetMetadataRowsToDelete() []MetadataFile {
	results := []MetadataFile{}

	filesMap := make(map[string]bool)
	for _, v := range FoundFiles {
		filesMap[v] = true
	}

	for _, v := range FoundMetadataFiles {
		fullFilePath := filepath.Join(v.filePath, v.fileName)
		if !filesMap[fullFilePath] {
			result := MetadataFile{
				filePath: v.filePath,
				fileName: v.fileName,
			}
			results = append(results, result)
		}
	}

	return results
}

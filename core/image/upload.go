package image

import (
	"gallery/core/config"
	"gallery/core/database"
	"gallery/core/optimised"
	"gallery/core/thumbnails"
	"gallery/core/types"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) string {
	fileName := fileHeader.Filename
	ext := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, ext) + strings.ToLower(ext)

	log.Printf("Uploading: %s", fileName)
	saveOriginalImage(file, fileName)
	filePath := filepath.Join(config.ImageDirectory, fileName)
	slug, _ := database.PopulateMetadataForUpload(fileName)
	thumbnails.GenerateThumbnail(filePath, slug)
	optimised.GenerateOptimised(filePath, slug)

	tags := types.TagsUpload{
		Tags:      []string{},
		ImageSlug: slug,
	}
	database.CreateTagsOnUpload(tags)
	database.CreateDimsensionsOnUpload(slug)
	return slug
}

func saveOriginalImage(file multipart.File, filename string) {

	filePath := filepath.Join(config.ImageDirectory, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("failed to parse uploaded file: %v", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Printf("failed to write uploaded file: %v", err)
		return
	}

	log.Printf("Uploaded file to %s", filePath)
}

func DeleteOriginalImage(filename string) error {

	existing := database.CheckMetadataByFileNameExists(filename)
	if existing {
		log.Printf("Original image %s is used by existing metadata, skipping deletion", filename)
		return nil
	}

	filePath := filepath.Join(config.ImageDirectory, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Original file does not exist for filename %s", filename)
		return err
	}
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error deleting original image %s: %s", filename, err)
		return err
	}
	log.Printf("Original image %s deleted", filename)

	return err
}

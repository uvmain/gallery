package image

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"photogallery/database"
	"photogallery/logic"
	"photogallery/optimised"
	"photogallery/thumbnails"
	"photogallery/types"
	"strings"
)

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) string {
	fileName := fileHeader.Filename
	ext := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, ext) + strings.ToLower(ext)

	log.Printf("Uploading: %s", fileName)
	saveOriginalImage(file, fileName)
	filePath := filepath.Join(logic.ImageDirectory, fileName)
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

	filePath := filepath.Join(logic.ImageDirectory, filename)
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

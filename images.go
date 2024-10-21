package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

var ExposureModes = map[int]string{
	0: "unknown",
	1: "Manual",
	2: "Normal program",
	3: "Aperture priority",
	4: "Shutter priority",
	5: "Creative program",
	6: "Action program",
	7: "Portrait mode",
	8: "Landscape mode",
}

type ImageMetadata struct {
	Slug         string    `json:"slug"`
	FilePath     string    `json:"filePath"`
	FileName     string    `json:"fileName"`
	Title        string    `json:"title"`
	DateTaken    time.Time `json:"dateTaken"`
	DateUploaded time.Time `json:"dateUploaded"`
	CameraMake   string    `json:"cameraMake"`
	CameraModel  string    `json:"cameraModel"`
	LensMake     string    `json:"lensMake"`
	LensModel    string    `json:"lensModel"`
	FStop        string    `json:"fStop"`
	ShutterSpeed string    `json:"shutterSpeed"`
	FlashStatus  string    `json:"flashStatus"`
	FocalLength  string    `json:"focalLength"`
	ISO          string    `json:"iso"`
	ExposureMode string    `json:"exposureMode"`
	WhiteBalance string    `json:"whiteBalance"`
	Albums       string    `json:"albums"`
}

func GetImageDirContents() ([]string, error) {
	var foundFiles []string

	err := filepath.Walk(ImagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error opening Images directory: %s", err)
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, validExt := range ImageExtensions {
				if ext == validExt {
					foundFiles = append(foundFiles, path)
					break
				}
			}
		}
		return nil
	})
	FoundFiles = foundFiles
	log.Printf("Found: %d source images", len(FoundFiles))
	return foundFiles, err
}

func GetSourceMetadataForImagePath(imagePath string) ImageMetadata {
	f, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", imagePath, err)
	}
	defer f.Close()

	filePath := filepath.Dir(imagePath)
	fileName := filepath.Base(imagePath)
	fileTitle := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	dateUploaded := time.Now()

	exifData, err := exif.Decode(f)
	if err != nil {
		log.Printf("Error decoding exif data for %s: %s", imagePath, err)
		return defaultImageMetadata(filePath, fileName, fileTitle, dateUploaded)
	}

	var exposureMode string
	tag := getExifTag(exifData, exif.ExposureProgram)
	if tag != nil {
		exposureMode = ExposureModes[getTagIntOrDefault(tag, 0)]
	} else {
		exposureMode = ExposureModes[0]
	}

	dateTaken, _ := getExifDate(exifData)

	cameraMake := getExifStringOrDefault(exifData, exif.Make, "unknown")
	cameraModel := getExifStringOrDefault(exifData, exif.Model, "unknown")
	lensModel := getExifStringOrDefault(exifData, exif.LensModel, "unknown")
	lensMake := getExifStringOrDefault(exifData, exif.LensMake, cameraMake)

	fStop := getExifStringOrDefault(exifData, exif.FNumber, "unknown")
	shutterSpeed := getExifStringOrDefault(exifData, exif.ShutterSpeedValue, "unknown")
	flashStatus := getExifStringOrDefault(exifData, exif.Flash, "unknown")
	focalLength := getExifStringOrDefault(exifData, exif.FocalLength, "unknown")
	iso := getExifStringOrDefault(exifData, exif.ISOSpeedRatings, "unknown")
	whiteBalance := getExifStringOrDefault(exifData, exif.WhiteBalance, "unknown")

	imageMetadata := ImageMetadata{
		Slug:         GenerateSlug(),
		FilePath:     filePath,
		FileName:     fileName,
		Title:        fileTitle,
		DateTaken:    dateTaken,
		DateUploaded: dateUploaded,
		CameraMake:   cameraMake,
		CameraModel:  cameraModel,
		LensMake:     lensMake,
		LensModel:    lensModel,
		FStop:        fStop,
		ShutterSpeed: shutterSpeed,
		FlashStatus:  flashStatus,
		FocalLength:  focalLength,
		ISO:          iso,
		ExposureMode: exposureMode,
		WhiteBalance: whiteBalance,
		Albums:       "[]",
	}

	return imageMetadata
}

func getExifTag(exifData *exif.Exif, field exif.FieldName) *tiff.Tag {
	tag, err := exifData.Get(field)
	if err != nil {
		return nil
	}
	return tag
}

func getTagIntOrDefault(tag *tiff.Tag, defaultValue int) int {
	tagString := tag.String()
	tagInt, err := strconv.Atoi(tagString)
	if err != nil {
		return defaultValue
	}
	return tagInt
}

func getExifStringOrDefault(exifData *exif.Exif, field exif.FieldName, defaultValue string) string {
	tag := getExifTag(exifData, field)
	if tag == nil {
		return defaultValue
	}
	tagString, _ := tag.StringVal()
	return tagString
}

func getExifDate(exifData *exif.Exif) (time.Time, error) {
	dateTaken, err := exifData.DateTime()
	if err != nil {
		tag := getExifTag(exifData, exif.DateTime)
		if tag == nil {
			log.Printf("no valid DateTime tag found")
			return time.Time{}, err
		}
		dateTaken, err = time.Parse(time.RFC1123, tag.String())
		if err != nil {
			log.Printf("failed to parse DateTime: %v", err)
			return time.Time{}, err
		}
	}
	return dateTaken, nil
}

func defaultImageMetadata(filePath, fileName, fileTitle string, dateUploaded time.Time) ImageMetadata {
	return ImageMetadata{
		Slug:         GenerateSlug(),
		FilePath:     filePath,
		FileName:     fileName,
		Title:        fileTitle,
		DateTaken:    dateUploaded,
		DateUploaded: dateUploaded,
		CameraMake:   "unknown",
		CameraModel:  "unknown",
		LensMake:     "unknown",
		LensModel:    "unknown",
		FStop:        "unknown",
		ShutterSpeed: "unknown",
		FlashStatus:  "unknown",
		FocalLength:  "unknown",
		ISO:          "unknown",
		ExposureMode: "unknown",
		WhiteBalance: "unknown",
		Albums:       "[]",
	}
}

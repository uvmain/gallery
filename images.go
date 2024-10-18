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

var FoundFiles []string

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
	slug         string
	filePath     string
	fileName     string
	title        string
	dateTaken    time.Time
	dateUploaded time.Time
	cameraMake   string
	cameraModel  string
	lensMake     string
	lensModel    string
	fStop        string
	shutterSpeed string
	flashStatus  string
	focalLength  string
	iso          string
	exposureMode string
	whiteBalance string
	albums       string
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

	dateTaken := getExifDateOrDefault(exifData, dateUploaded)

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
		slug:         GenerateSlug(),
		filePath:     filePath,
		fileName:     fileName,
		title:        fileTitle,
		dateTaken:    dateTaken,
		dateUploaded: dateUploaded,
		cameraMake:   cameraMake,
		cameraModel:  cameraModel,
		lensMake:     lensMake,
		lensModel:    lensModel,
		fStop:        fStop,
		shutterSpeed: shutterSpeed,
		flashStatus:  flashStatus,
		focalLength:  focalLength,
		iso:          iso,
		exposureMode: exposureMode,
		whiteBalance: whiteBalance,
		albums:       "[]",
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
	return tag.String()
}

func getExifDateOrDefault(exifData *exif.Exif, defaultValue time.Time) time.Time {
	dateTaken, err := exifData.DateTime()
	if err != nil {
		tag := getExifTag(exifData, exif.DateTime)
		if tag == nil {
			return defaultValue
		}
		dateTaken, err = time.Parse(time.RFC1123, tag.String())
		if err != nil {
			return defaultValue
		}
	}
	return dateTaken
}

func defaultImageMetadata(filePath, fileName, fileTitle string, dateUploaded time.Time) ImageMetadata {
	return ImageMetadata{
		slug:         GenerateSlug(),
		filePath:     filePath,
		fileName:     fileName,
		title:        fileTitle,
		dateTaken:    dateUploaded,
		dateUploaded: dateUploaded,
		cameraMake:   "unknown",
		cameraModel:  "unknown",
		lensMake:     "unknown",
		lensModel:    "unknown",
		fStop:        "unknown",
		shutterSpeed: "unknown",
		flashStatus:  "unknown",
		focalLength:  "unknown",
		iso:          "unknown",
		exposureMode: "unknown",
		whiteBalance: "unknown",
		albums:       "[]",
	}
}

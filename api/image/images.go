package image

import (
	"log"
	"os"
	"path/filepath"
	"photogallery/logic"
	"photogallery/types"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func GetSourceMetadataForImagePath(imagePath string) types.ImageMetadata {
	f, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Error opening file %s: %s", imagePath, err)
	}
	defer f.Close()

	var tag *tiff.Tag

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
	tag = getExifTag(exifData, exif.ExposureProgram)
	if tag != nil {
		exposureMode = types.ExposureModes[getTagIntOrDefault(tag, 0)]
	} else {
		exposureMode = types.ExposureModes[0]
	}

	var whiteBalance string
	tag = getExifTag(exifData, exif.WhiteBalance)
	if tag != nil {
		tagString := tag.String()
		if tagString == "1" {
			whiteBalance = "Manual"
		} else if tagString == "0" {
			whiteBalance = "Auto"
		}
	} else {
		whiteBalance = "unknown"
	}

	var flashStatus string
	tag = getExifTag(exifData, exif.Flash)
	if tag != nil {
		flashStatus = types.FlashModes[getTagIntOrDefault(tag, 0)]
	} else {
		flashStatus = "unknown"
	}

	var whiteBalanceMode string
	tag = getExifTag(exifData, exif.LightSource)
	if tag != nil {
		whiteBalanceMode = types.WhiteBalanceModes[getTagIntOrDefault(tag, 0)]
	} else {
		whiteBalanceMode = types.WhiteBalanceModes[0]
	}

	dateTaken, _ := getExifDate(exifData)

	cameraMake := getExifStringOrDefault(exifData, exif.Make, "unknown")
	cameraModel := getExifStringOrDefault(exifData, exif.Model, "unknown")
	lensModel := getExifStringOrDefault(exifData, exif.LensModel, "unknown")
	lensMake := getExifStringOrDefault(exifData, exif.LensMake, cameraMake)

	fStop := getExifStringOrDefault(exifData, exif.FNumber, "unknown")
	shutterSpeed := getExifStringOrDefault(exifData, exif.ShutterSpeedValue, "unknown")
	exposureTime := getExifStringOrDefault(exifData, exif.ExposureTime, shutterSpeed)
	focalLength := getExifStringOrDefault(exifData, exif.FocalLength, "unknown")
	iso := getExifStringOrDefault(exifData, exif.ISOSpeedRatings, "unknown")

	imageMetadata := types.ImageMetadata{
		Slug:             logic.GenerateSlug(),
		FilePath:         filePath,
		FileName:         fileName,
		Title:            fileTitle,
		DateTaken:        dateTaken,
		DateUploaded:     dateUploaded,
		CameraMake:       cameraMake,
		CameraModel:      cameraModel,
		LensMake:         lensMake,
		LensModel:        lensModel,
		FStop:            fStop,
		ExposureTime:     exposureTime,
		FlashStatus:      flashStatus,
		FocalLength:      focalLength,
		ISO:              iso,
		ExposureMode:     exposureMode,
		WhiteBalance:     whiteBalance,
		WhiteBalanceMode: whiteBalanceMode,
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
	tagString := tag.String()
	tagString = strings.ReplaceAll(tagString, "\"", "")
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

func defaultImageMetadata(filePath, fileName, fileTitle string, dateUploaded time.Time) types.ImageMetadata {
	return types.ImageMetadata{
		Slug:             logic.GenerateSlug(),
		FilePath:         filePath,
		FileName:         fileName,
		Title:            fileTitle,
		DateTaken:        dateUploaded,
		DateUploaded:     dateUploaded,
		CameraMake:       "unknown",
		CameraModel:      "unknown",
		LensMake:         "unknown",
		LensModel:        "unknown",
		FStop:            "unknown",
		ExposureTime:     "unknown",
		FlashStatus:      "unknown",
		FocalLength:      "unknown",
		ISO:              "unknown",
		ExposureMode:     "unknown",
		WhiteBalance:     "unknown",
		WhiteBalanceMode: "unknown",
	}
}

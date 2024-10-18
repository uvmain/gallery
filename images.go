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
			log.Fatalf("Error creating thumbnails directory: %s", err)
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
	return foundFiles, err
}

func GetSourceMetadataForImagePath(imagePath string) ImageMetadata {

	f, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", imagePath, err)
	}
	defer f.Close()

	exifData, err := exif.Decode(f)
	if err != nil {
		log.Fatalf("Error decoding exif data for %s: %s", imagePath, err)
	}

	filePath := filepath.Dir(imagePath)
	fileName := filepath.Base(imagePath)
	fileTitle := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	dateUploaded := time.Now()

	var tag *tiff.Tag

	var exposureMode string
	tag, err = exifData.Get(exif.ExposureProgram)
	if err != nil {
		exposureMode = ExposureModes[0]
	} else {
		exposureModeTag := tag.String()
		exposureModeInt, err := strconv.Atoi(exposureModeTag)
		if err != nil {
			exposureMode = ExposureModes[0]
		} else {
			exposureMode = ExposureModes[exposureModeInt]
		}
	}

	dateTaken, err := exifData.DateTime()
	if err != nil {
		value, getErr := exifData.Get(exif.DateTime)
		dateTaken, _ = time.Parse(time.RFC1123, value.String())
		if getErr != nil {
			dateTaken = time.Now()
		}
	}

	var cameraMake string
	tag, err = exifData.Get(exif.Make)
	if err != nil {
		cameraMake = "unknown"
	} else {
		cameraMake = ToTitle(tag.String())
	}

	var cameraModel string
	tag, err = exifData.Get(exif.Model)
	if err != nil {
		cameraModel = "unknown"
	} else {
		cameraModel = tag.String()
	}

	var lensModel string
	tag, err = exifData.Get(exif.LensModel)
	if err != nil {
		lensModel = "unknown"
	} else {
		lensModel = tag.String()
	}

	var lensMake string
	tag, err = exifData.Get(exif.LensMake)
	if err != nil {
		if (lensModel != "unknown") && (cameraMake != "unknown") {
			lensMake = cameraMake
		} else {
			lensMake = "unknown"
		}
	} else {
		lensMake = ToTitle(tag.String())
	}

	var fStop string
	tag, err = exifData.Get(exif.FNumber)
	if err != nil {
		fStop = "unknown"
	} else {
		fStop = tag.String()
	}

	var shutterSpeed string
	tag, err = exifData.Get(exif.ShutterSpeedValue)
	if err != nil {
		shutterSpeed = "unknown"
	} else {
		shutterSpeed = tag.String()
	}

	var flashStatus string
	tag, err = exifData.Get(exif.Flash)
	if err != nil {
		flashStatus = "unknown"
	} else {
		flashStatus = tag.String()
	}

	var focalLength string
	tag, err = exifData.Get(exif.FocalLength)
	if err != nil {
		focalLength = "unknown"
	} else {
		focalLength = tag.String()
	}

	var iso string
	tag, err = exifData.Get(exif.ISOSpeedRatings)
	if err != nil {
		iso = "unknown"
	} else {
		iso = tag.String()
	}

	var whiteBalance string
	tag, err = exifData.Get(exif.WhiteBalance)
	if err != nil {
		whiteBalance = "unknown"
	} else {
		whiteBalance = tag.String()
	}

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
		albums:       "[test, test2]",
	}

	return imageMetadata
}

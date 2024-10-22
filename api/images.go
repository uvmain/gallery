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

var WhiteBalanceModes = map[int]string{
	0:   "Unknown",
	1:   "Daylight",
	2:   "Fluorescent",
	3:   "Tungsten",
	4:   "Flash",
	9:   "Fine Weather",
	10:  "Cloudy Weather",
	11:  "Shade",
	12:  "Daylight Fluorescent",
	13:  "Day White Fluorescent",
	14:  "Cool White Fluorescent",
	15:  "White Fluorescent",
	17:  "Standard Light A",
	18:  "Standard Light B",
	19:  "Standard Light C",
	20:  "D55",
	21:  "D65",
	22:  "D75",
	23:  "D50",
	24:  "ISO Studio Tungsten",
	255: "Other Light Source",
}

var FlashModes = map[int]string{
	0:  "No Flash",
	1:  "Fired",
	5:  "Fired, Return not detected",
	7:  "Fired, Return detected",
	8:  "On, Did not fire",
	9:  "On, Fired",
	11: "On, Return not detected",
	15: "On, Return detected",
	16: "Off, Did not fire",
	20: "Off, Did not fire, Return not detected",
	24: "Auto, Did not fire",
	25: "Auto, Fired",
	29: "Auto, Fired, Return not detected",
	31: "Auto, Fired, Return detected",
	32: "No flash function",
	48: "Off, No flash function",
	65: "Fired, Red-eye reduction",
	69: "Fired, Red-eye reduction, Return not detected",
	71: "Fired, Red-eye reduction, Return detected",
	73: "On, Red-eye reduction",
	77: "On, Red-eye reduction, Return not detected",
	79: "On, Red-eye reduction, Return detected",
	80: "Off, Red-eye reduction",
	88: "Auto, Did not fire, Red-eye reduction",
	89: "Auto, Fired, Red-eye reduction",
	93: "Auto, Fired, Red-eye reduction, Return not detected",
	95: "Auto, Fired, Red-eye reduction, Return detected",
}

type ImageMetadata struct {
	Slug             string    `json:"slug"`
	FilePath         string    `json:"filePath"`
	FileName         string    `json:"fileName"`
	Title            string    `json:"title"`
	DateTaken        time.Time `json:"dateTaken"`
	DateUploaded     time.Time `json:"dateUploaded"`
	CameraMake       string    `json:"cameraMake"`
	CameraModel      string    `json:"cameraModel"`
	LensMake         string    `json:"lensMake"`
	LensModel        string    `json:"lensModel"`
	FStop            string    `json:"fStop"`
	ExposureTime     string    `json:"exposureTime"`
	FlashStatus      string    `json:"flashStatus"`
	FocalLength      string    `json:"focalLength"`
	ISO              string    `json:"iso"`
	ExposureMode     string    `json:"exposureMode"`
	WhiteBalance     string    `json:"whiteBalance"`
	WhiteBalanceMode string    `json:"whiteBalanceMode"`
	Albums           string    `json:"albums"`
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
	log.Printf("Found: %d source images", len(foundFiles))
	return foundFiles, err
}

func GetSourceMetadataForImagePath(imagePath string) ImageMetadata {
	f, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", imagePath, err)
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
		exposureMode = ExposureModes[getTagIntOrDefault(tag, 0)]
	} else {
		exposureMode = ExposureModes[0]
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
		flashStatus = FlashModes[getTagIntOrDefault(tag, 0)]
	} else {
		flashStatus = "unknown"
	}

	var whiteBalanceMode string
	tag = getExifTag(exifData, exif.LightSource)
	if tag != nil {
		whiteBalanceMode = WhiteBalanceModes[getTagIntOrDefault(tag, 0)]
	} else {
		whiteBalanceMode = WhiteBalanceModes[0]
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

	imageMetadata := ImageMetadata{
		Slug:             GenerateSlug(),
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
		Albums:           "[]",
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

func defaultImageMetadata(filePath, fileName, fileTitle string, dateUploaded time.Time) ImageMetadata {
	return ImageMetadata{
		Slug:             GenerateSlug(),
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
		Albums:           "[]",
	}
}

func GetOriginalImageBySlug(slug string) ([]byte, error) {
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

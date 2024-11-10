package logic

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

var DatabaseDirectory string
var ThumbnailDirectory string
var OptimisedDirectory string
var ImagePath string
var ImageExtensions []string
var ThumbnailMaxPixels uint
var OptimisedMaxPixels uint

func LoadEnv() {

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}

	ImagePath = os.Getenv("IMAGE_PATH")
	if ImagePath == "" {
		ImagePath = "./images"
	}

	DatabaseDirectory, _ = filepath.Abs(dataPath)
	ThumbnailDirectory, _ = filepath.Abs(filepath.Join(dataPath, "thumbnails"))
	OptimisedDirectory, _ = filepath.Abs(filepath.Join(dataPath, "optimised"))

	imageExtensions := os.Getenv("IMAGE_EXTENSIONS")
	if imageExtensions == "" {
		ImageExtensions = []string{".avif", ".bmp", ".gif", ".jpg", ".jpeg", ".png", ".webp"}
	} else {
		ImageExtensions = strings.Split(imageExtensions, ",")
	}

	u, _ := strconv.ParseUint(os.Getenv("THUMBNAIL_MAX_PIXELS"), 10, 64)
	if u > 0 {
		ThumbnailMaxPixels = uint(u)
	} else {
		ThumbnailMaxPixels = 500
	}

	u, _ = strconv.ParseUint(os.Getenv("OPTIMISED_MAX_PIXELS"), 10, 64)
	if u > 0 {
		OptimisedMaxPixels = uint(u)
	} else {
		OptimisedMaxPixels = 1280
	}
}

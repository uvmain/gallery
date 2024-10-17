package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	LoadEnv()

	DatabaseDirectory, _ = filepath.Abs(os.Getenv("DATA_PATH"))
	ThumbnailDirectory, _ = filepath.Abs(filepath.Join(os.Getenv("DATA_PATH"), "thumbnails"))
	ImagePath, _ = filepath.Abs(os.Getenv("IMAGE_PATH"))
	ImageExtensions = strings.Split(os.Getenv("IMAGE_FILES"), ",")
	value, _ := strconv.ParseUint(os.Getenv("THUMBNAIL_MAX_PIXELS"), 10, 64)
	ThumbnailMaxPixels = uint(value)

	Database = InitialiseDatabase()
	CreateThumbnailsDir()
	GetImageDirContents()
}

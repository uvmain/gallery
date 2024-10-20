package main

import (
	_ "modernc.org/sqlite"
)

func main() {
	LoadEnv()
	InitialiseDatabase()
	GetImageDirContents()
	InitialiseMetadata()
	InitialiseThumbnails()

	CreateOptimisedDir()

	// GenerateThumbnail(FoundFiles[photoInt], imageMetadata.slug)
}

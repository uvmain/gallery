package main

import (
	"gallery/core/config"
	"gallery/core/database"
	"gallery/core/optimised"
	"gallery/core/thumbnails"
)

func main() {
	config.LoadEnv()
	database.Initialise()
	database.InitialiseMetadata()
	database.InitialiseTags()
	database.InitialiseDimensions()
	thumbnails.InitialiseThumbnails()
	optimised.InitialiseOptimised()
	StartServer()
}

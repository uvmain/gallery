package main

import (
	"photogallery/database"
	"photogallery/logic"
	"photogallery/optimised"
	"photogallery/thumbnails"
)

func main() {
	logic.LoadEnv()
	database.Initialise()
	database.InitialiseMetadata()
	database.InitialiseTags()
	database.InitialiseDimensions()
	thumbnails.InitialiseThumbnails()
	optimised.InitialiseOptimised()
	StartServer()
}

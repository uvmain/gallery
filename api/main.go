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
	thumbnails.InitialiseThumbnails()
	optimised.InitialiseOptimised()
	StartServer()
}

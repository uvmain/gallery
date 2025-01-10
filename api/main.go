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
	thumbnails.InitialiseThumbnails()
	optimised.InitialiseOptimised()
	StartServer()
}

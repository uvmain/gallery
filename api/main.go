package main

import (
	"photogallery/database"
	"photogallery/logic"
)

func main() {
	logic.LoadEnv()
	database.Initialise()
	database.InitialiseMetadata()
	InitialiseThumbnails()
	InitialiseOptimised()

	StartServer()
}

package main

func main() {
	LoadEnv()
	InitialiseDatabase()
	GetImageDirContents()
	InitialiseMetadata()
	InitialiseThumbnails()
	InitialiseOptimised()

	StartServer()
}

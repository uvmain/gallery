package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var Database *sql.DB
var DatabaseDirectory string
var ThumbnailDirectory string
var OptimisedDirectory string
var ImagePath string
var ImageExtensions []string
var ThumbnailMaxPixels uint
var OptimisedMaxPixels uint

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	} else {
		log.Println(".env file loaded")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

func listDirRecursive(dir string) ([]string, error) {
	var files []string

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		filePath := filepath.Join(dir, fileInfo.Name())
		if !fileInfo.IsDir() {
			files = append(files, filePath)
		}

		if fileInfo.IsDir() {
			subDirFiles, err := listDirRecursive(filePath)
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}

	return files, nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <directory> <output-json-file>\n", os.Args[0])
	}

	dir := os.Args[1]
	outputFile := os.Args[2]

	// Get the directory contents recursively
	files, err := listDirRecursive(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}

	// Write the output to a JSON file
	jsonData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v\n", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v\n", err)
	}

	fmt.Printf("Directory contents successfully written to %s\n", outputFile)
}

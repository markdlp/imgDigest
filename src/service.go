package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/barasher/go-exiftool"
)

func GetDates(inputFolder string, fileType string) ([]string, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return nil, err
	}
	defer et.Close()

	dir := inputFolder + "/" + fileType

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileDates := []string{}
	for _, file := range files {
		fileInfos := et.ExtractMetadata(dir + `/` + file.Name())

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
				continue
			}

			for k, v := range fileInfo.Fields {
				fmt.Printf("[%v] %v\n", k, v)
			}
		}
	}

	return fileDates, nil
}

func ProcessFilesByType(inputFolder string) ([]string, error) {
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return nil, err
	}

	fileTypes := []string{}
	for _, file := range files {
		if file.IsDir() {
			if !slices.Contains(fileTypes, file.Name()) {
				fileTypes = append(fileTypes, file.Name())
			}
			continue
		} // Skip directories

		fileType := filepath.Ext(file.Name())
		filePath := filepath.Join(inputFolder, file.Name())
		subfolderPath := filepath.Join(inputFolder, fileType)

		// Create subfolder if it doesn't exist
		if err := os.MkdirAll(subfolderPath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("error creating folder %s: %w", subfolderPath, err)
		}

		// Move the file
		newFilePath := filepath.Join(subfolderPath, file.Name())
		if err := os.Rename(filePath, newFilePath); err != nil {
			return nil, fmt.Errorf("error moving file %s: %w", file.Name(), err)
		}
	}

	return fileTypes, nil
}

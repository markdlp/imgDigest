package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/barasher/go-exiftool"
)

func GetDates(inputFolder string, fileTypes []string) ([][]string, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return nil, err
	}
	defer et.Close()

	fileDates := make([][]string, len(fileTypes))
	for i, fileType := range fileTypes {

		fileDates[i] = append(fileDates[i], fileType)

		dir := inputFolder + `/` + fileType

		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			fileInfos := et.ExtractMetadata(dir + `/` + file.Name())

			for _, fileInfo := range fileInfos {
				if fileInfo.Err != nil {
					fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
					continue
				}

				fileDates[i] = append(fileDates[i],
					strings.ReplaceAll(fileInfo.Fields["CreateDate"].(string), ":", "-"))
			}
		}
	}

	return fileDates, nil
}

func setNames(inputFolder string, outputFolder string, fileDates [][]string) error {

	dupCount := 0 // never > 9
	for i := 1; i < len(fileDates); i++ {
		for j := 0; j < len(fileDates[i]); j++ {
			if fileDates[i][j] == fileDates[i][j-1][:len(fileDates[0])] {
				dupCount++
				fileDates[i][j] = fileDates[i][j] + `_` + fmt.Sprint(dupCount)
			} else {
				dupCount = 0
			}
		}
	}

	// Create outputfolder if it doesn't exist
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating folder %s: %w", outputFolder, err)
	}

	for i := 0; i < len(fileDates); i++ {

		subfolderPath := filepath.Join(outputFolder, fileDates[i][0])

		// Create outputfolder if it doesn't exist
		if err := os.MkdirAll(subfolderPath, os.ModePerm); err != nil {
			return fmt.Errorf("error creating folder %s: %w", subfolderPath, err)
		}

		files, err := os.ReadDir(inputFolder + `/` + fileDates[i][0])
		if err != nil {
			return err
		}

		for j := 1; j < len(fileDates[i]); j++ {

			newFilePath := filepath.Join(outputFolder, fileDates[i][0], fileDates[i][j])
			oldFilePath := filepath.Join(inputFolder, fileDates[i][0], files[j-1].Name())
			// Move the file
			if err := os.Rename(oldFilePath, newFilePath+fileDates[i][0]); err != nil {
				return fmt.Errorf("error moving file %s: %w", files[j-1].Name(), err)
			}
		}
	}
	return nil
}

func ProcessFilesByType(inputFolder string) ([]string, error) {
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return nil, err
	}

	fileTypes := []string{}
	for _, file := range files {

		fileType := filepath.Ext(file.Name())

		if !slices.Contains(fileTypes, fileType) {
			fileTypes = append(fileTypes, fileType)
		}

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

func compressFolder(inputFolder string) {
	file, err := os.Create("../output.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(inputFolder, walker)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetFiles(c *gin.Context) {

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	files := form.File["files"]

	for _, file := range files {
		log.Println(file.Filename)
		filename := filepath.Base(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, "../upload/"+filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

	tmpFolder := "../upload"

	fileTypes, err := ProcessFilesByType(tmpFolder)
	if err != nil {
		fmt.Println("Error organizing files:", err)
	} else {
		fmt.Println("Files organized successfully!")
	}

	fileDates, _ := GetDates(tmpFolder, fileTypes)
	setNames(tmpFolder, tmpFolder, fileDates)
}

func SendFile(c *gin.Context) {

	file, err := compressFolder("../upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Error zipping folder:", err.Error())
		return
	}

	// Set download headers
	c.Header("Content-Disposition", "attachment; filename=../output.zip")
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", len(file)))
	c.Writer.Write(file)
}

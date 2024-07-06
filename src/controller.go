package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var tmpFolder = "../upload"

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
		c.SaveUploadedFile(file, tmpFolder+"/"+filename)
	}
}

func SendFile(c *gin.Context) {

	fileTypes, err := ProcessFilesByType(tmpFolder)
	if err != nil {
		fmt.Println("Error organizing files:", err)
	} else {
		fmt.Println("Files organized successfully!")
	}

	fileDates, _ := GetDates(tmpFolder, fileTypes)
	setNames(tmpFolder, tmpFolder, fileDates)

	compressFolder(tmpFolder)

	file, err := os.ReadFile("../output.zip")
	if err != nil {
		c.String(http.StatusBadRequest, "Error zipping folder:", err.Error())
		return
	}

	// Set download headers
	c.Header("Content-Disposition", "attachment; filename=output.zip")
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", len(file)))
	c.Writer.Write(file)

	e := os.Remove("../output.zip")
	if e != nil {
		log.Fatal(e)
	}

	e = os.RemoveAll(tmpFolder)
	if e != nil {
		log.Fatal(e)
	}
}

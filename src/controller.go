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

	inputFolder := "../upload"

	fileTypes, err := ProcessFilesByType(inputFolder)
	if err != nil {
		fmt.Println("Error organizing files:", err)
	} else {
		fmt.Println("Files organized successfully!")
	}

	fileDates, _ := GetDates(inputFolder, fileTypes)
	setNames(inputFolder, "../download", fileDates)
	compressFolder("../download")

	c.Header("Content-Disposition", "attachment; filename=output.zip")
	http.ServeFile(c.Writer, c.Request, "../output.zip")
}

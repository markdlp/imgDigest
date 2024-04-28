package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func getPhotos(c *gin.Context) {

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	files := form.File["photos"]

	for _, file := range files {
		log.Println(file.Filename)
		filename := filepath.Base(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func main() {

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2", "10.0.0.0/8"})

	router.LoadHTMLFiles("templates/index.tmpl")

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.POST("/", getPhotos)

	router.Run("localhost:8080")
}

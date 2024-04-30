package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/gin-gonic/gin"
)

func getFiles(c *gin.Context) {

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
		c.SaveUploadedFile(file, "upload/"+filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func processFilesByType(inputFolder string) error {
	files, err := os.ReadDir("upload")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		} // Skip directories

		fileType := filepath.Ext(file.Name())

		filePath := filepath.Join(inputFolder, file.Name())
		subfolderPath := filepath.Join(inputFolder, fileType)

		// Create subfolder if it doesn't exist
		if err := os.MkdirAll(subfolderPath, os.ModePerm); err != nil {
			return fmt.Errorf("error creating folder %s: %w", subfolderPath, err)
		}

		// Move the file
		newFilePath := filepath.Join(subfolderPath, file.Name())
		if err := os.Rename(filePath, newFilePath); err != nil {
			return fmt.Errorf("error moving file %s: %w", file.Name(), err)
		}
	}

	return nil
}

func getDates(inputFolder string) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata("testdata/20190404_131804.jpg")

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

func main() {

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2", "10.0.0.0/8"})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Static("/", "public")

	router.POST("/upload", getFiles)

	inputFolder := "upload"

	err := processFilesByType(inputFolder)
	   	if err != nil {
	   		fmt.Println("Error organizing files:", err)
	   	} else {
	   		fmt.Println("Files organized successfully!")
	   	}
	

	getDates(inputFolder)
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

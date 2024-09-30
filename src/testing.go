package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testFolder = "../testing_upload"

func TestGetFiles(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	router.POST("/upload", GetFiles)

	// Prepare the test file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := os.CreateTemp("", "testfile*.txt")
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // Clean up the temp file after test

	// Write file content
	part, err := writer.CreateFormFile("files", filepath.Base(file.Name()))
	assert.NoError(t, err)
	part.Write([]byte("this is a test file"))

	writer.Close()

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check response code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Check if file is saved correctly in the tmpFolder
	savedFilePath := filepath.Join(testFolder, filepath.Base(file.Name()))
	_, err = os.Stat(savedFilePath)
	assert.NoError(t, err)
	os.Remove(savedFilePath) // Clean up the saved file after test
}

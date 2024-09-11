package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Static("/upload", "public")

	router.POST("/upload", GetFiles)
	router.GET("/download", SendFile)

	router.Run(":8080")
}

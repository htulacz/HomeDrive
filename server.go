package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadDir = "uploads"

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could't read the file"})
		return
	}

	path := filepath.Join(uploadDir, file.Filename)
	c.SaveUploadedFile(file, path)
	c.JSON(http.StatusOK, gin.H{"message": "Plik przesłany", "url": "/uploads/" + file.Filename})
}

func main() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	router := gin.Default()
	router.POST("/upload", upload)

	router.Run(":8080")
}

package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadDir = "upload"

func upload(c *gin.Context) {
	path_rest := c.Query("directories")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read the file"})
		return
	}

	path := filepath.Join(uploadDir, path_rest, file.Filename)
	c.SaveUploadedFile(file, path)
	c.JSON(http.StatusOK, gin.H{"message": "File sent", "url": "/uploads/" + file.Filename})
}

func makeDirectory(c *gin.Context) {
	dirName := c.PostForm("name")

	if dirName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No directory name"})
		return
	}

	fullPath := filepath.Join(uploadDir, dirName)
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create directory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Directory created", "path": fullPath})
}

func main() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	r := gin.Default()
	r.LoadHTMLGlob("site/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Hello"})
	})

	r.POST("/upload/*directories", upload)
	r.POST("makeDirectory/*parent_dir", makeDirectory)

	r.Run(":8080")
}

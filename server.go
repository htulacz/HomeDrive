package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

const uploadDir = "upload"

func upload(c *gin.Context) {
	path_rest := c.PostForm("directory")
	if path_rest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No directory specified"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read files"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files selected"})
		return
	}

	dirPath := filepath.Join(uploadDir, path_rest)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}

	var wg sync.WaitGroup
	var uploadedFiles []string
	var mutex sync.Mutex

	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {
			defer wg.Done()

			filePath := filepath.Join(dirPath, file.Filename)
			err := c.SaveUploadedFile(file, filePath)
			if err != nil {
				fmt.Println("Error saving file:", err)
				return
			}

			mutex.Lock()
			uploadedFiles = append(uploadedFiles, "/upload/"+path_rest+"/"+file.Filename)
			mutex.Unlock()
		}(file)
	}

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded", "files": uploadedFiles})
}

func getDirectories(c *gin.Context) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't read directories"})
		return
	}

	var directories []string
	for _, file := range files {
		directories = append(directories, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{"directories": directories})
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

func getContent(c *gin.Context) {
	directory := c.DefaultQuery("directory", "")
	if directory == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Directory name is required"})
		return
	}

	path := filepath.Join(uploadDir, directory)
	fmt.Println("Opening directory: ", path)

	files, err := os.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't open directory"})
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": fileNames})
}

func deletePhoto(c *gin.Context) {
	directory := c.DefaultQuery("directory", "")
	if directory == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Directory name is required"})
		return
	}

	photoName := c.DefaultQuery("photo_name", "")
	if photoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo name is required"})
		return
	}

	path := filepath.Join(uploadDir, directory, photoName)
	fmt.Println(path)
	err := os.Remove(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": photoName})
}

func main() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 200 << 20
	r.LoadHTMLGlob("index.html")
	r.Static("/upload", "./upload")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Upload File"})
	})

	r.POST("/upload", upload)
	r.POST("/create-directory", makeDirectory)
	r.GET("/directories", getDirectories)
	r.GET("/directoryContent", getContent)
	r.DELETE("/delete", deletePhoto)
	r.Run(":8080")
}

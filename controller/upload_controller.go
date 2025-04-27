package controller

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
    // "encoding/base64"
    // "fmt"
)

type UploadInput struct {
	Input struct {
		FileBase64 string `json:"file_base64"`
	} `json:"input"`
}

func UploadSingleHandler(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
        return
    }

    fileURL, err := saveUploadedFile(c,file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "file_path": fileURL,
    })
}


func UploadMultipleHandler(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
        return
    }

    files := form.File["files"] // key = files (array of files)
    var filePaths []string

    for _, file := range files {
        fileURL, err := saveUploadedFile(c,file)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save one of the files"})
            return
        }
        filePaths = append(filePaths, fileURL)
    }

    c.JSON(http.StatusOK, gin.H{
        "file_paths": filePaths,
    })
}

func saveUploadedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
    os.MkdirAll("./uploads", os.ModePerm)

    extension := filepath.Ext(file.Filename)
    newFileName := uuid.New().String() + "-" + time.Now().Format("20060102150405") + extension
    uploadPath := filepath.Join("./uploads", newFileName)

    if err := c.SaveUploadedFile(file, uploadPath); err != nil {
        return "", err
    }

    return "/uploads/" + newFileName, nil
}
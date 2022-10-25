package main

import (
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/disintegration/imaging"
)

func uploadResizeSingleFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := "http://localhost:8000/images/single/" + filename

	imageFile, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	src := imaging.Resize(imageFile, 1000, 0, imaging.Lanczos)
	err = imaging.Save(src, fmt.Sprintf("public/single/%v", filename))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"filepath": filePath})
}
func uploadSingleFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := "http://localhost:8000/images/single/" + filename

	out, err := os.Create("public/single/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, gin.H{"filepath": filePath})
}
func uploadResizeMultipleFile(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["images"]
	filePaths := []string{}
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath := "http://localhost:8000/images/multiple/" + filename

		filePaths = append(filePaths, filePath)
		readerFile, _ := file.Open()
		imageFile, _, err := image.Decode(readerFile)
		if err != nil {
			log.Fatal(err)
		}
		src := imaging.Resize(imageFile, 1000, 0, imaging.Lanczos)
		err = imaging.Save(src, fmt.Sprintf("public/multiple/%v", filename))
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
}
func uploadMultipleFile(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["images"]
	filePaths := []string{}
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath := "http://localhost:8000/images/multiple/" + filename

		filePaths = append(filePaths, filePath)
		out, err := os.Create("./public/multiple/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		readerFile, _ := file.Open()
		_, err = io.Copy(out, readerFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"filepath": filePaths})
}
func init() {
	if _, err := os.Stat("public/single"); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("public/single", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	if _, err := os.Stat("public/multiple"); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("public/multiple", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "success", "message": "How to Upload Single and Multiple Files in Golang"})
	})

	router.POST("/upload/single", uploadResizeSingleFile)
	router.POST("/upload/multiple", uploadResizeMultipleFile)
	router.StaticFS("/images", http.Dir("public"))
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	router.Run(":8000")
}

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	images := loadImages()
	boundary := "--boundary"

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.GET("/mjpeg", func(c *gin.Context) {
		c.Header("content-type", "multipart/x-mixed-replace; boundary="+boundary)
		imageIdx := 0
		for {
			c.Writer.Write([]byte(fmt.Sprintf("\r\n--%s\r\n", boundary)))
			c.Writer.Write([]byte("Content-Type: image/jpeg\r\n"))
			c.Writer.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(images[imageIdx]))))
			c.Writer.Write(images[imageIdx])
			c.Writer.Write([]byte(boundary))
			c.Writer.Flush()

			imageIdx++
			imageIdx = imageIdx % len(images)

			time.Sleep(30 * time.Millisecond)
		}
	})

	router.Run(":3333")
}

func loadImages() (output [][]byte) {
	baseDir := "images"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	imageFileNames := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".jpg") {
			imageFileNames = append(imageFileNames, file.Name())
		}
	}

	sort.Strings(imageFileNames)

	for _, fileName := range imageFileNames {
		bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", baseDir, fileName))
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, bytes)
	}

	return
}

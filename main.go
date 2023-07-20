package main

import (
	"fmt"
	"go-cloud-storage/cloudbucket"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting...")

	r := gin.Default()

	// Ping endpoint for testing the server status
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello! There",
		})
	})

	// Endpoint for uploading a file to the cloud storage bucket
	r.POST("/upload-to-bucket", cloudbucket.HandleFileUploadToBucket)
	// Endpoint for retrieving a file from the cloud storage bucket
	r.GET("/get-file-from-bucket/:filename", cloudbucket.GetFileFromBucket)
	// Endpoint for listing all objects in the cloud storage bucket
	r.GET("/list-objects-in-bucket", cloudbucket.ListObjectsInBucket)
	// Endpoint for generating a signed URL for accessing a file in the cloud storage bucket
	r.GET("/get-signed-url/:filename", cloudbucket.GetAuthenticatedURL)
	// Endpoint for retrieving and storing all objects from the cloud storage bucket
	r.GET("/get-all-objects", cloudbucket.GetAllObjectsFromBucket)
	// Endpoint for deleting an object from the cloud storage bucket
	r.DELETE("/delete-object/:filename", cloudbucket.DeleteObjectFromBucket)
	// Endpoint for updating an object from the cloud storage bucket
	r.PUT("/update-object/:filename", cloudbucket.UpdateObjectInBucket)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

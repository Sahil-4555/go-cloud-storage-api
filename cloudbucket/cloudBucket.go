package cloudbucket

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

// HandleFileUploadToBucket uploads file to bucket
func HandleFileUploadToBucket(c *gin.Context) {
	// Function to handle file upload to the bucket
	bucket := "bucket_golang" //your bucket name

	// Error variable
	var err error

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Get the uploaded file from the request
	f, uploadedFile, err := c.Request.FormFile("file")
	if err != nil {
		// Return error response if file retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Close the file after the function ends
	defer f.Close()

	// Create a new storage writer to write the file to the bucket
	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)

	// Copy the file content to the storage writer
	if _, err := io.Copy(sw, f); err != nil {
		// Return error response if file copy fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Close the storage writer
	if err := sw.Close(); err != nil {
		// Return error response if storage writer closing fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Generate the URL for the uploaded file
	u, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		// Return error response if URL generation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	// Return success response with the URL
	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}

func GetFileFromBucket(c *gin.Context) {
	// Function to get a file from the bucket
	bucket := "bucket_golang"
	filename := c.Param("filename")
	folder := "File_From_Bucket"

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Get the storage object for the specified file
	obj := storageClient.Bucket(bucket).Object(filename)

	// Create a new reader for the object
	r, err := obj.NewReader(ctx)
	if err != nil {
		// Return error response if reader creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Close the reader after the function ends
	defer r.Close()

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		// Return error response if working directory retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Define the new file path with the folder and filename
	newFilePath := filepath.Join(wd, folder, filename)

	// Create a file in the desired storage location
	newFile, err := os.Create(newFilePath)
	if err != nil {
		// Return error response if file creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	defer newFile.Close()

	// Copy the file content from the reader to the new file
	_, err = io.Copy(newFile, r)
	if err != nil {
		// Return error response if file copy fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Return success response with the new file path
	c.JSON(http.StatusOK, gin.H{
		"message": "file stored in the new folder successfully",
		"path":    newFilePath,
	})
}

func ListObjectsInBucket(c *gin.Context) {
	// Function to list all objects in the bucket
	bucket := "bucket_golang" // your bucket name

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Get the bucket object
	bkt := storageClient.Bucket(bucket)

	// Channel to receive object names
	objectNamesChan := make(chan string)

	// Goroutine to list objects and send their names to the channel
	go func() {
		it := bkt.Objects(ctx, nil)
		for {
			objAttrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// Return error response if object retrieval fails
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				return
			}
			objectNamesChan <- objAttrs.Name
		}
		close(objectNamesChan)
	}()

	// Receive object names from the channel and collect them
	var objectNames []string
	for objName := range objectNamesChan {
		objectNames = append(objectNames, objName)
	}

	c.JSON(http.StatusOK, gin.H{
		"objects": objectNames,
	})
}

func GetAuthenticatedURL(c *gin.Context) {
	// Function to get an authenticated URL for a file in the bucket
	bucket := "bucket_golang"       // your bucket name
	filename := c.Param("filename") // retrieve the filename from the request URL parameter

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	defer storageClient.Close()

	// Set options for the signed URL
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	// Generate the signed URL for the file
	u, err := storageClient.Bucket(bucket).SignedURL(filename, opts)
	if err != nil {
		// Return error response if URL generation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Return success response with the signed URL
	c.JSON(http.StatusOK, gin.H{
		"url": u,
	})
}

func GetAllObjectsFromBucket(c *gin.Context) {
	// Function to retrieve and store all objects from the bucket
	bucket := "bucket_golang"
	folder := "All_Files_From_Bucket"

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		// Return error response if working directory retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Define the full path of the new folder
	fullFolderPath := filepath.Join(wd, folder)

	// Create the new folder if it doesn't exist
	if _, err := os.Stat(fullFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullFolderPath, 0755)
		if err != nil {
			// Return error response if folder creation fails
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
	}

	// List objects in the bucket
	it := storageClient.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Return error response if object retrieval fails
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		// Retrieve the object's name
		filename := attrs.Name

		// Get the object's content
		r, err := storageClient.Bucket(bucket).Object(filename).NewReader(ctx)
		if err != nil {
			// Return error response if reader creation fails
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		defer r.Close()

		// Define the new file path with the folder and filename
		newFilePath := filepath.Join(fullFolderPath, filename)

		// Read the object's content
		data, err := ioutil.ReadAll(r)
		if err != nil {
			// Return error response if content reading fails
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		// Write the object's content to the new file
		err = ioutil.WriteFile(newFilePath, data, 0644)
		if err != nil {
			// Return error response if file writing fails
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
	}

	// Return success response with the folder path
	c.JSON(http.StatusOK, gin.H{
		"message": "all objects stored in the new folder successfully",
		"path":    folder,
	})
}

func DeleteObjectFromBucket(c *gin.Context) {
	// Function to delete an object from the bucket
	bucket := "bucket_golang"         // your bucket name
	objectName := c.Param("filename") // retrieve the object name from the request URL parameter

	// Create a new App Engine context
	ctx := appengine.NewContext(c.Request)

	// Create a new storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		// Return error response if client creation fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	defer storageClient.Close()

	// Get the storage object for the specified object

	// Delete the object from the bucket
	obj := storageClient.Bucket(bucket).Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		// Return error response if object deletion fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Object %s deleted successfully", objectName),
	})
}

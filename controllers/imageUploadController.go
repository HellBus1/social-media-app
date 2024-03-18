package controllers

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUploadImage(ctx *gin.Context) {
	form, fileErr := ctx.MultipartForm()

	if fileErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process multipart form"})
		return
	}


	files := form.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	file := files[0]
	if file.Size <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is empty"})
		return
	}

	// Check file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file must be in *.jpg or *.jpeg format"})
		return
	}

	// Check file size
	if file.Size > 2*1024*1024 || file.Size < 10*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file size must be between 10KB and 2MB"})
		return
	}

	// Generate filename
	filename := uuid.New().String() + ext

	uploadError := uploadToS3(file)
	if uploadError != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to AWS S3"})
		return
	}

	// Construct file URL
	fileURL := "https://s3.amazonaws.com/" + os.Getenv("S3_BUCKET_NAME") + "/" + filename

	// Respond with file URL
	ctx.JSON(http.StatusOK, gin.H{"imageUrl": fileURL})
}

// Uploads file to AWS S3
func uploadToS3(file *multipart.FileHeader) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), ""),
	})
	if err != nil {
		return err
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Upload file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(file.Filename),
		ACL:    aws.String("public-read"),
		Body:   fileContent,
	})

	if err != nil {
		return err
	}

	return nil
}

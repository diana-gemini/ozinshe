package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func ImageUpload(c *gin.Context, images []*multipart.FileHeader) ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	}))

	uploader := s3manager.NewUploader(sess)

	var imageURLs []string

	for _, image := range images {
		src, err := image.Open()
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "failed to open file")
			return nil, err
		}
		defer src.Close()

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("ozinwe-diana"),
			Key:    aws.String(image.Filename),
			Body:   src,
		})
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "failed to upload file to S3")
			return nil, err
		}
		imageURL := fmt.Sprintf("https://ozinwe-diana.s3.eu-north-1.amazonaws.com/%s", image.Filename)
		imageURLs = append(imageURLs, imageURL)

	}

	return imageURLs, nil
}

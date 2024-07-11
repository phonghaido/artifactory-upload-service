package handlers

import (
	"context"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/phonghaido/artifactory-upload-service/configs"
	"github.com/phonghaido/artifactory-upload-service/helpers"
	l "github.com/sirupsen/logrus"
)

func newS3Client(cfg configs.Config) (*s3.Client, error) {
	s3Cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(s3Cfg), nil
}

func uploadFile(ctx context.Context, s3Client *s3.Client, src multipart.File, fileName string, bucket string) error {
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   src,
	})
	return err
}

func HandlePostUploadFile(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return helpers.MethodNotAllowed()
	}

	config, err := configs.GetConfig()
	if err != nil {
		return err
	}

	s3Client, err := newS3Client(*config)
	if err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return helpers.InvalidFile(err)
	}

	if file.Size > config.MaxSize {
		return helpers.InvalidFileSize()
	}

	src, err := file.Open()
	if err != nil {
		return helpers.InvalidFile(err)
	}
	defer src.Close()

	err = uploadFile(c.Request().Context(), s3Client, src, file.Filename, config.S3Bucket)
	if err != nil {
		return err
	}

	return nil
}

func HandlePostUploadFiles(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return helpers.MethodNotAllowed()
	}

	config, err := configs.GetConfig()
	if err != nil {
		return err
	}

	s3Client, err := newS3Client(*config)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return helpers.InvalidFile(err)
	}
	files := form.File["files"]

	var wg sync.WaitGroup
	for _, file := range files {
		if file.Size > config.MaxSize {
			return helpers.InvalidFileSize()
		}

		src, err := file.Open()
		if err != nil {
			return helpers.InvalidFile(err)
		}
		defer src.Close()

		wg.Add(1)
		go func(ctx context.Context, client *s3.Client, src multipart.File, fileName string, bucket string) {
			defer wg.Done()
			err := uploadFile(ctx, client, src, fileName, bucket)
			if err != nil {
				l.Errorf("uploaded file %s failed: %v", fileName, err)
			} else {
				l.Infof("successfully uploaded file %s", fileName)
			}
		}(c.Request().Context(), s3Client, src, file.Filename, config.S3Bucket)
	}

	wg.Wait()
	l.Infoln("Successfully uploaded all files")

	return nil
}

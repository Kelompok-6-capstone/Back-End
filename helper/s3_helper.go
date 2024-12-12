package helper

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	Client     *s3.Client
	BucketName string
}

func NewS3Uploader() (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		Client:     client,
		BucketName: os.Getenv("AWS_S3_BUCKET"),
	}, nil
}

func (u *S3Uploader) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	defer file.Close()

	// Baca file ke buffer
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := fileBuffer.ReadFrom(file); err != nil {
		return "", err
	}

	// Buat nama file unik
	fileName := fmt.Sprintf("%s/%s", folder, filepath.Base(fileHeader.Filename))

	// Upload file ke S3
	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileBuffer.Bytes()),
		ACL:    "public-read", // Pastikan file bisa diakses secara publik
	})
	if err != nil {
		return "", err
	}

	// URL file yang diunggah
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.BucketName, os.Getenv("AWS_REGION"), fileName)
	return fileURL, nil
}

func (u *S3Uploader) DeleteFile(fileKey string) error {
	_, err := u.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(fileKey),
	})
	return err
}

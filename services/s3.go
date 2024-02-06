package services

import (
	"app-service-com/config"
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadS3(file multipart.File, filename string, size int64) error {
	aws_access_key_id := config.Get(`aws.access_key`)
	aws_secret_access_key := config.Get(`aws.secret_key`)
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	cfg := aws.NewConfig().WithRegion(config.Get(`aws.region`)).WithCredentials(creds)
	svc := s3.New(session.Must(session.NewSession()), cfg)

	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	// "/name_image.jpg"
	path := "/" + filename
	params := &s3.PutObjectInput{
		Bucket:        aws.String(config.Get(`aws.bucket`)),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, err := svc.PutObject(params)

	if err != nil {
		WriteLog(err)
		return err
	}
	return nil
}

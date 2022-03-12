package repository

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/aelpxy/xoniaapp/model"
	"github.com/aelpxy/xoniaapp/model/apperrors"
	"github.com/aelpxy/xoniaapp/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"mime/multipart"
)

type s3FileRepository struct {
	S3Session  *session.Session
	BucketName string
}

func NewFileRepository(session *session.Session, bucketName string) model.FileRepository {
	return &s3FileRepository{
		S3Session:  session,
		BucketName: bucketName,
	}
}

func (s *s3FileRepository) UploadAvatar(header *multipart.FileHeader, directory string) (string, error) {
	uploader := s3manager.NewUploader(s.S3Session)

	id, _ := service.GenerateId()
	key := fmt.Sprintf("files/%s/%s.jpeg", directory, id)

	file, err := header.Open()

	if err != nil {
		log.Printf("Failed to open header: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	src, _, err := image.Decode(file)

	if err != nil {
		log.Printf("Failed to decode image: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	img := imaging.Resize(src, 150, 0, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})

	if err != nil {
		log.Printf("Failed to encode image: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	up, err := uploader.Upload(&s3manager.UploadInput{
		Body:        buf,
		Bucket:      aws.String(s.BucketName),
		ContentType: aws.String("image/jpeg"),
		Key:         aws.String(key),
	})

	if err != nil {
		log.Printf("Failed to upload file: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = file.Close(); err != nil {
		log.Printf("Failed to close file: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return up.Location, nil
}

func (s *s3FileRepository) UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, error) {
	uploader := s3manager.NewUploader(s.S3Session)

	key := fmt.Sprintf("files/%s/%s", directory, filename)

	file, err := header.Open()

	if err != nil {
		log.Printf("Failed to open header: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	up, err := uploader.Upload(&s3manager.UploadInput{
		Body:        file,
		Bucket:      aws.String(s.BucketName),
		ContentType: aws.String(mimetype),
		Key:         aws.String(key),
	})

	if err != nil {
		log.Printf("Failed to upload file: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = file.Close(); err != nil {
		log.Printf("Failed to close file: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return up.Location, nil
}

func (s *s3FileRepository) DeleteImage(key string) error {
	srv := s3.New(s.S3Session)
	_, err := srv.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Printf("Failed to delete image: %v\n", err.Error())
		return apperrors.NewInternal()
	}

	return nil
}

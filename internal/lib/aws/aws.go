package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	awslib "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3lib "github.com/aws/aws-sdk-go/service/s3"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
)

type S3Service struct{}

func NewS3Service() *S3Service {
	return &S3Service{}
}

func (s S3Service) CreatePresignedUploadURL(key string) (string, error) {
	c := config.Conf.AWS
	s3, err := s3Client()
	if err != nil {
		return "", err
	}

	r, _ := s3.PutObjectRequest(&s3lib.PutObjectInput{
		Bucket: awslib.String(c.S3Bucket),
		Key:    awslib.String(key),
	})

	url, err := r.Presign(c.S3PresignedURLExpirySeconds)
	if err != nil {
		logger.Error("Unable to create presigned upload url: %v", err)
		return "", err
	}

	return url, nil
}

func s3Client() (*s3lib.S3, error) {
	c := config.Conf.AWS

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})
	if err != nil {
		logger.Error("Unable to create AWS session: %v", err)
		return nil, err
	}

	return s3.New(sess), nil
}

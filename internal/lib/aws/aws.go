package aws

import (
	"time"

	awslib "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

func NewSession() (*session.Session, error) {
	c := config.GetS3()

	s, err := session.NewSession(&awslib.Config{
		Region:      awslib.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})

	if err != nil {
		log.Error("Unable to create aws session: %v", err)

		return nil, err
	}

	return s, nil
}

func CreatePresignedUploadUrl(key string) (string, error) {
	c := config.GetS3()
	sess, err := NewSession()

	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: awslib.String(c.Bucket),
		Key:    awslib.String(key),
	})

	url, err := req.Presign(time.Duration(time.Duration(c.PresignedUrlExpiry) * time.Minute))

	if err != nil {
		log.Error("Unable to create presigned upload url: %v", err)

		return "", err
	}

	return url, err
}

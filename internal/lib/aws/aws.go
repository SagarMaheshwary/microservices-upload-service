package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

func NewS3Session() (*s3.S3, error) {
	c := config.GetS3()

	sess, err := session.NewSession(&awssdk.Config{
		Region:      awssdk.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})

	if err != nil {
		log.Error("Unable to create s3 session: %v", err)

		return nil, err
	}

	svc := s3.New(sess)

	return svc, nil
}

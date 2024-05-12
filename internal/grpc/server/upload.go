package server

import (
	"context"
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	cons "github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type uploadServer struct {
	pb.UploadServer
}

func (u *uploadServer) CreatePresignedUrl(ctx context.Context, data *pb.CreatePresignedUrlRequest) (*pb.CreatePresignedUrlResponse, error) {
	c := config.GetS3()

	svc, err := aws.NewS3Session()

	if err != nil {
		log.Error("Unable to create s3 session: %v", err)

		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	key := fmt.Sprintf("video-%s", uuid.New().String())

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: awssdk.String(c.Bucket),
		Key:    awssdk.String(key),
	})

	url, err := req.Presign(time.Duration(15 * time.Minute))

	if err != nil {
		log.Error("Unable to create presigned url from request: %v", err)

		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	response := &pb.CreatePresignedUrlResponse{Url: url}

	return response, nil
}

func (u *uploadServer) UploadedWebhook(ctx context.Context, data *pb.UploadedWebhookRequest) (*pb.UploadedWebhookResponse, error) {
	//

	return nil, nil
}

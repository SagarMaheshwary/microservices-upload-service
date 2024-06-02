package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	cons "github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
	pb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type uploadServer struct {
	pb.UploadServer
}

func (u *uploadServer) CreatePresignedUrl(ctx context.Context, data *pb.CreatePresignedUrlRequest) (*pb.CreatePresignedUrlResponse, error) {
	uploadId := uuid.New().String()

	url, err := aws.CreatePresignedUploadUrl(fmt.Sprintf("%s/%s", cons.RawVideosDirectory, uploadId))

	if err != nil {
		log.Error("gRPC.CreatePresignedUrl unable to create presigned url from request: %v", err)

		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	response := &pb.CreatePresignedUrlResponse{
		Message: cons.MessageOK,
		Data: &pb.CreatePresignedUrlResponseData{
			UploadId: uploadId,
			Url:      url,
		},
	}

	return response, nil
}

func (u *uploadServer) UploadedWebhook(ctx context.Context, data *pb.UploadedWebhookRequest) (*pb.UploadedWebhookResponse, error) {
	type EncodeUploadedVideo struct {
		UploadId    string `json:"upload_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PublishedAt string `json:"published_at"`
		UserId      int    `json:"user_id"`
	}

	err := publisher.P.Publish(cons.QueueEncodeService, &broker.MessageType{
		Key: cons.MessageTypeEncodeUploadedVideo,
		Data: &EncodeUploadedVideo{
			UploadId:    data.UploadId,
			Title:       data.Title,
			Description: data.Description,
			PublishedAt: time.Now().Format(time.RFC3339),
			UserId:      1, //@TODO: send current user id
		},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	response := &pb.UploadedWebhookResponse{
		Message: cons.MessageOK,
		Data:    &pb.UploadedWebhookResponseData{},
	}

	return response, nil
}

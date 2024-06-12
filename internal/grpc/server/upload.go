package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	cons "github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
	pb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type uploadServer struct {
	pb.UploadServiceServer
}

func (u *uploadServer) CreatePresignedUrl(ctx context.Context, data *pb.CreatePresignedUrlRequest) (*pb.CreatePresignedUrlResponse, error) {
	videoId := uuid.New().String()
	thumbnailId := uuid.New().String()

	videoUrl, err := aws.CreatePresignedUploadUrl(fmt.Sprintf("%s/%s", cons.S3RawVideosDirectory, videoId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	thumbnailUrl, err := aws.CreatePresignedUploadUrl(fmt.Sprintf("%s/%s", cons.S3ThumbnailsDirectory, thumbnailId))

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	response := &pb.CreatePresignedUrlResponse{
		Message: cons.MessageOK,
		Data: &pb.CreatePresignedUrlResponseData{
			VideoId:      videoId,
			VideoUrl:     videoUrl,
			ThumbnailId:  thumbnailId,
			ThumbnailUrl: thumbnailUrl,
		},
	}

	return response, nil
}

func (u *uploadServer) UploadedWebhook(ctx context.Context, data *pb.UploadedWebhookRequest) (*pb.UploadedWebhookResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	id, exists := helper.GetFromMetadata(md, cons.HeaderUserId)

	if !exists {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	userId, _ := strconv.Atoi(id)

	type EncodeUploadedVideo struct {
		VideoId     string `json:"video_id"`
		ThumbnailId string `json:"thumbnail_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PublishedAt string `json:"published_at"`
		UserId      int    `json:"user_id"`
	}

	err := publisher.P.Publish(cons.QueueEncodeService, &broker.MessageType{
		Key: cons.MessageTypeEncodeUploadedVideo,
		Data: &EncodeUploadedVideo{
			VideoId:     data.VideoId,
			ThumbnailId: data.ThumbnailId,
			Title:       data.Title,
			Description: data.Description,
			PublishedAt: time.Now().Format(time.RFC3339),
			UserId:      userId,
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

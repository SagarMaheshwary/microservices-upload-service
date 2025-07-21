package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/publisher"
	uploadpb "github.com/sagarmaheshwary/microservices-upload-service/internal/proto/upload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type uploadServer struct {
	uploadpb.UploadServiceServer
	s3Storage aws.S3Storage
}

type EncodeUploadedVideo struct {
	VideoId     string `json:"video_id"`
	ThumbnailId string `json:"thumbnail_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
	UserId      int    `json:"user_id"`
}

func (u *uploadServer) CreatePresignedUrl(ctx context.Context, data *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error) {
	videoID := uuid.New().String()
	thumbnailID := uuid.New().String()

	videoURL, err := u.s3Storage.CreatePresignedUploadURL(fmt.Sprintf("%s/%s", constant.S3RawVideosDirectory, videoID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
	}

	thumbnailURL, err := u.s3Storage.CreatePresignedUploadURL(fmt.Sprintf("%s/%s", constant.S3ThumbnailsDirectory, thumbnailID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
	}

	response := &uploadpb.CreatePresignedUrlResponse{
		Message: constant.MessageOK,
		Data: &uploadpb.CreatePresignedUrlResponseData{
			VideoId:      videoID,
			VideoUrl:     videoURL,
			ThumbnailId:  thumbnailID,
			ThumbnailUrl: thumbnailURL,
		},
	}
	return response, nil
}

func (u *uploadServer) UploadedWebhook(ctx context.Context, data *uploadpb.UploadedWebhookRequest) (*uploadpb.UploadedWebhookResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	id, exists := helper.GetGRPCMetadataValue(md, constant.HeaderUserID)
	if !exists {
		return nil, status.Errorf(codes.Unauthenticated, constant.MessageUnauthorized)
	}

	userId, _ := strconv.Atoi(id)
	err := publisher.P.Publish(ctx, constant.QueueEncodeService, &publisher.MessageType{
		Key: constant.MessageTypeEncodeUploadedVideo,
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
		return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
	}

	response := &uploadpb.UploadedWebhookResponse{
		Message: constant.MessageOK,
		Data:    &uploadpb.UploadedWebhookResponseData{},
	}
	return response, nil
}

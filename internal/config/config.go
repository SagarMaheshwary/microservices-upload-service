package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	S3         *s3
	AMQP       *amqp
}

type grpcServer struct {
	Host string
	Port int
}

type s3 struct {
	Bucket             string
	Region             string
	AccessKey          string
	SecretKey          string
	PresignedUrlExpiry int
}

type amqp struct {
	Host           string
	Port           int
	Username       string
	Password       string
	PublishTimeout time.Duration
}

func Init() {
	envPath := path.Join(helper.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load %q: %v", envPath, err)
	}

	log.Info("Loaded %q", envPath)

	port, err := strconv.Atoi(Getenv("GRPC_PORT", "5002"))

	if err != nil {
		log.Error("Invalid GRPC_PORT value %v", err)
	}

	amqpPort, err := strconv.Atoi(Getenv("AMQP_PORT", "5672"))

	if err != nil {
		log.Error("Invalid AMQP_PORT value %v", err)
	}

	amqpPublishTimeout, err := strconv.Atoi(Getenv("AMQP_PUBLISH_TIMEOUT_SECONDS", "5"))

	if err != nil {
		log.Error("Invalid AMQP_PORT value %v", err)
	}

	s3UrlExpiry, err := strconv.Atoi(Getenv("AWS_S3_PRESIGNED_URL_EXPIRY", "0"))

	if err != nil {
		log.Error("Invalid AWS_S3_PRESIGNED_URL_EXPIRY value %v", err)
	}

	conf = &Config{
		GRPCServer: &grpcServer{
			Host: Getenv("GRPC_HOST", "localhost"),
			Port: port,
		},
		S3: &s3{
			Bucket:             Getenv("AWS_S3_BUCKET", ""),
			Region:             Getenv("AWS_S3_REGION", ""),
			AccessKey:          Getenv("AWS_S3_ACCESS_KEY", ""),
			SecretKey:          Getenv("AWS_S3_SECRET_KEY", ""),
			PresignedUrlExpiry: s3UrlExpiry,
		},
		AMQP: &amqp{
			Host:           Getenv("AMQP_HOST", "localhost"),
			Port:           amqpPort,
			Username:       Getenv("AMQP_USERNAME", "guest"),
			Password:       Getenv("AMQP_PASSWORD", "guest"),
			PublishTimeout: time.Duration(amqpPublishTimeout) * time.Second,
		},
	}
}

func GetgrpcServer() *grpcServer {
	return conf.GRPCServer
}

func GetS3() *s3 {
	return conf.S3
}

func Getamqp() *amqp {
	return conf.AMQP
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

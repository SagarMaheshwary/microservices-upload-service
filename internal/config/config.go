package config

import (
	"os"
	"path"
	"strconv"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	S3         *s3
}

type grpcServer struct {
	Host string
	Port int
}

type s3 struct {
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

func Init() {
	envPath := path.Join(helper.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load .env %q: %v", envPath, err)
	}

	log.Info("Loaded %q", envPath)

	port, err := strconv.Atoi(Getenv("GRPC_PORT", "5002"))

	if err != nil {
		log.Error("Invalid GRPC_PORT value %v", err)
	}

	conf = &Config{
		GRPCServer: &grpcServer{
			Host: Getenv("GRPC_HOST", "localhost"),
			Port: port,
		},
		S3: &s3{
			Bucket:    Getenv("AWS_S3_BUCKET", ""),
			Region:    Getenv("AWS_S3_REGION", ""),
			AccessKey: Getenv("AWS_S3_ACCESS_KEY", ""),
			SecretKey: Getenv("AWS_S3_SECRET_KEY", ""),
		},
	}
}

func GetgrpcServer() *grpcServer {
	return conf.GRPCServer
}

func GetS3() *s3 {
	return conf.S3
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

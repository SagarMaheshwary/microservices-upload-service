package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
)

var Conf *Config

type Config struct {
	GRPCServer *GRPCServer
	AWS        *AWS
	AMQP       *AMQP
	GRPCClient *GRPCClient
	Prometheus *Prometheus
	Jaeger     *Jaeger
}

type GRPCServer struct {
	Host string
	Port int
}

type AWS struct {
	Region               string
	S3Bucket             string
	AccessKey            string
	SecretKey            string
	S3PresignedURLExpiry int
}

type AMQP struct {
	Host           string
	Port           int
	Username       string
	Password       string
	PublishTimeout time.Duration
}

type GRPCClient struct {
	EncodeServiceURL string
	Timeout          time.Duration
}

type Prometheus struct {
	URL string
}

type Jaeger struct {
	URL string
}

func Init() {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := env.Load(envPath); err != nil {
			logger.Fatal("Failed to load .env %q: %v", envPath, err)
		}

		logger.Info("Loaded environment variables from %q", envPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	Conf = &Config{
		GRPCServer: &GRPCServer{
			Host: getEnv("GRPC_HOST", "localhost"),
			Port: getEnvInt("GRPC_PORT", 5002),
		},
		AWS: &AWS{
			Region:               getEnv("AWS_REGION", ""),
			AccessKey:            getEnv("AWS_ACCESS_KEY", ""),
			SecretKey:            getEnv("AWS_SECRET_KEY", ""),
			S3Bucket:             getEnv("AWS_S3_BUCKET", ""),
			S3PresignedURLExpiry: getEnvInt("AWS_S3_PRESIGNED_URL_EXPIRY", 15),
		},
		AMQP: &AMQP{
			Host:           getEnv("AMQP_HOST", "localhost"),
			Port:           getEnvInt("AMQP_PORT", 5672),
			Username:       getEnv("AMQP_USERNAME", "guest"),
			Password:       getEnv("AMQP_PASSWORD", "guest"),
			PublishTimeout: getEnvDuration("AMQP_PUBLISH_TIMEOUT_SECONDS", 5),
		},
		GRPCClient: &GRPCClient{
			EncodeServiceURL: getEnv("GRPC_ENCODE_SERVICE_URL", "localhost:5004"),
			Timeout:          getEnvDuration("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
		},
		Prometheus: &Prometheus{
			URL: getEnv("PROMETHEUS_URL", "localhost:5012"),
		},
		Jaeger: &Jaeger{
			URL: getEnv("JAEGER_URL", "localhost:4318"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}

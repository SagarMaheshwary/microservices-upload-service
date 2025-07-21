package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	encoderpc "github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/client/encode"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/prometheus"
	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	config.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx)

	promServer := prometheus.NewServer()
	go func() {
		if err := prometheus.Serve(promServer); err != nil && err != http.ErrServerClosed {
			stop()
		}
	}()

	go broker.MaintainConnection(ctx)

	encoderpc.NewClient(ctx)

	grpcServer := server.NewServer()
	go func() {
		if err := server.Serve(grpcServer); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	if err := shutdownJaeger(shutdownCtx); err != nil {
		logger.Warn("jaeger server shutdown error: %v", err)
	}

	shutdownCtx, cancel = context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	if err := promServer.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Prometheus server shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	logger.Info("Shutdown complete")
}

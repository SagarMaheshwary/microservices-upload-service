package prometheus

import (
	"net/http"

	prometheuslib "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/config"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/logger"
)

var (
	GRPCRequestCounter = prometheuslib.NewCounterVec(
		prometheuslib.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCRequestLatency = prometheuslib.NewHistogramVec(
		prometheuslib.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) of gRPC requests",
			Buckets: prometheuslib.DefBuckets,
		},
		[]string{"method"},
	)

	ServiceHealth = prometheuslib.NewGauge(prometheuslib.GaugeOpts{
		Name: "service_health_status",
		Help: "Health status of the service: 1=Healthy, 0=Unhealthy",
	})
)

func NewServer() *http.Server {
	prometheuslib.MustRegister(GRPCRequestCounter, GRPCRequestLatency, ServiceHealth)

	url := config.Conf.Prometheus.URL

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: url, Handler: mux}

	logger.Info("Starting prometheus metrics server %s", url)

	return server
}

func Serve(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Prometheus server error: %v", err)
	}
}

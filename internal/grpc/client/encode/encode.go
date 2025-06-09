package encode

import (
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var Encode *encodeClient

type encodeClient struct {
	health healthpb.HealthClient
}

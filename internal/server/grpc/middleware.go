package internalgrpc

import (
	"context"
	"time"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc"
)

func LoggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		response, err := handler(ctx, req)

		code := grpcLogging.DefaultErrorToCode(err)
		latency := time.Since(startTime)

		logger.Info("[GRPC] method=%s code=%s latency=%s", info.FullMethod, code, latency)
		return response, err
	}
}

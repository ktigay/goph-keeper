package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// WithLogging логирует запрос.
func WithLogging(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		l.Info(
			"request",
			"request", req,
		)

		resp, err = handler(ctx, req)

		l.Info(
			"response",
			"response", resp,
			"duration", time.Since(start),
			"size", proto.Size(resp.(proto.Message)),
			"error", err,
		)
		return resp, err
	}
}

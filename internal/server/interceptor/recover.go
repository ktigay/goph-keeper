package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithRecover перехват panic.
func WithRecover(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if e := recover(); e != nil {
				logger.Error("Recovering from", "error", e, "stack", string(debug.Stack()))

				resp = nil
				err = status.Errorf(codes.Internal, "panic: %v", e)
			}
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

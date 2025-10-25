package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor интерцептор для установки таймаута.
func TimeoutInterceptor(timeout int64) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

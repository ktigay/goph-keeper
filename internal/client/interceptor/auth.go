package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthService сервис авторизации.
type AuthService interface {
	GetJWT(ctx context.Context) (string, error)
}

// AuthInterceptor интерцептор аутентификации.
func AuthInterceptor(s AuthService) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		jwt, _ := s.GetJWT(ctx)
		if jwt != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

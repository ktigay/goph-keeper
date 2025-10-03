package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ktigay/goph-keeper/internal/entity"
	c "github.com/ktigay/goph-keeper/internal/server/context"
)

const (
	authorizationHeader = "authorization"
)

// JWTWrapper JWT обработчик.
type JWTWrapper interface {
	GenerateToken(payload entity.Identity) (string, error)
	ParseToken(s string) (*entity.Identity, error)
}

// Auth структура для авторизации.
type Auth struct {
	jwt        JWTWrapper
	accessList map[string]bool
}

// WithAuthorization интерцептор для работы с авторизацией.
func (i *Auth) WithAuthorization() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if ctx, err = i.authorization(ctx, info.FullMethod); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func (i *Auth) authorization(ctx context.Context, method string) (context.Context, error) {
	var (
		md       metadata.MD
		values   []string
		identity *entity.Identity
		err      error
		check    bool
		ok       bool
	)

	if check, ok = i.accessList[method]; !ok {
		return nil, fmt.Errorf("permission denied for method: %s", method)
	}

	if !check {
		return ctx, nil
	}

	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		return nil, fmt.Errorf("failed to extract authorization header")
	}

	if values, ok = md[authorizationHeader]; !ok || len(values) < 1 {
		return ctx, nil
	}

	if identity, err = i.jwt.ParseToken(values[0]); err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	return c.NewContextWithIdentity(ctx, *identity), nil
}

// NewAuth конструктор.
func NewAuth(jwt JWTWrapper, accessList map[string]bool) *Auth {
	return &Auth{
		jwt:        jwt,
		accessList: accessList,
	}
}

package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
	e "github.com/ktigay/goph-keeper/internal/entity"
	"github.com/ktigay/goph-keeper/internal/server/entity"
)

// AuthService сервис аутентификации.
//
//go:generate mockgen -destination=./mocks/mock_auth.go -package=mocks github.com/ktigay/goph-keeper/internal/server/handler/grpc AuthService
type AuthService interface {
	Register(ctx context.Context, login, password string) (*entity.User, error)
	Login(ctx context.Context, login, password string) (*entity.User, error)
}

// JWTWrapper обработчик JWT.
//
//go:generate mockgen -destination=./mocks/mock_jwt.go -package=mocks github.com/ktigay/goph-keeper/internal/server/handler/grpc JWTWrapper
type JWTWrapper interface {
	GenerateToken(payload e.Identity) (string, error)
	ParseToken(s string) (*e.Identity, error)
}

// AuthHandler обработчик аутентификации.
type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	srv AuthService
	jwt JWTWrapper
}

// Register регистрирует пользователя.
func (a *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	var (
		usr *entity.User
		err error
	)

	if usr, err = a.srv.Register(ctx, req.GetLogin(), req.GetPassword()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{
		UserUuid: usr.UUID,
	}, nil
}

// Login авторизует пользователя.
func (a *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var (
		usr   *entity.User
		token string
		err   error
	)

	if usr, err = a.srv.Login(ctx, req.GetLogin(), req.GetPassword()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, err = a.jwt.GenerateToken(e.Identity{
		UUID: usr.UUID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}

// NewAuthHandler конструктор.
func NewAuthHandler(s AuthService, j JWTWrapper) *AuthHandler {
	return &AuthHandler{
		srv: s,
		jwt: j,
	}
}

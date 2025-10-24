package auth

import (
	"context"
	"log/slog"

	"github.com/ktigay/goph-keeper/internal/client/entity"
	"github.com/ktigay/goph-keeper/internal/client/validator"
)

// Client клиент.
//
//go:generate mockgen -destination=./mocks/mock_client.go -package=mocks github.com/ktigay/goph-keeper/internal/client/service/auth Client
type Client interface {
	Login(ctx context.Context, data entity.Credentials) (string, error)
	Register(ctx context.Context, data entity.Credentials) (string, error)
}

// Repository репозиторий.
//
//go:generate mockgen -destination=./mocks/mock_auth.go -package=mocks github.com/ktigay/goph-keeper/internal/client/service/auth Repository
type Repository interface {
	SetJWT(ctx context.Context, jwt string) error
	GetJWT(ctx context.Context) (string, error)
}

// Service сервис.
type Service struct {
	client Client
	repo   Repository
	logger *slog.Logger
}

// Login авторизирует пользователя.
func (s *Service) Login(ctx context.Context, data entity.Credentials) error {
	if err := validator.ValidateCredentials(data); err != nil {
		return err
	}

	token, err := s.client.Login(ctx, data)
	if err != nil {
		s.logger.Debug("login failed", "error", err)
		return err
	}
	s.logger.Debug("login success", "token", token)
	return s.repo.SetJWT(ctx, token)
}

// Register регистрирует пользователя.
func (s *Service) Register(ctx context.Context, data entity.Credentials) error {
	if err := validator.ValidateCredentials(data); err != nil {
		return err
	}

	uuid, err := s.client.Register(ctx, data)
	if err != nil {
		return err
	}
	s.logger.Debug("register success", "uuid", uuid)
	return nil
}

// GetJWT возвращает JWT токен.
func (s *Service) GetJWT(ctx context.Context) (string, error) {
	jwt, err := s.repo.GetJWT(ctx)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

// IsAuthorized возвращает true, если пользователь авторизован.
func (s *Service) IsAuthorized(ctx context.Context) bool {
	jwt, _ := s.GetJWT(ctx)
	return jwt != ""
}

// New конструктор.
func New(c Client, r Repository, l *slog.Logger) *Service {
	return &Service{
		client: c,
		repo:   r,
		logger: l,
	}
}

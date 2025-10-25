package auth

import (
	"context"

	"github.com/ktigay/goph-keeper/internal/client/entity"
	e "github.com/ktigay/goph-keeper/internal/entity"
)

// Service сервис аутентификации.
type Service interface {
	Login(ctx context.Context, data entity.Credentials) error
	Register(ctx context.Context, data entity.Credentials) error
}

// SyncService сервис синхронизации данных.
type SyncService interface {
	Initialize(ctx context.Context) ([]e.UserData, error)
	SyncFromRemote(ctx context.Context) error
}

// Handler обработчик аутентификации.
type Handler struct {
	srv     Service
	syncSrv SyncService
}

// SignIn авторизует пользователя.
func (h *Handler) SignIn(ctx context.Context, l entity.Credentials) error {
	var err error

	if err = h.srv.Login(ctx, l); err != nil {
		return err
	}
	if _, err = h.syncSrv.Initialize(ctx); err != nil {
		return err
	}
	return nil
}

// SignUp регистрирует пользователя.
func (h *Handler) SignUp(ctx context.Context, l entity.Credentials) error {
	var err error

	if err = h.srv.Register(ctx, l); err != nil {
		return err
	}
	return nil
}

// New конструктор.
func New(srv Service, syncSrv SyncService) *Handler {
	return &Handler{
		srv:     srv,
		syncSrv: syncSrv,
	}
}

package userdata

import (
	"context"

	"github.com/ktigay/goph-keeper/internal/entity"
)

// Service сервис пользовательских данных.
type Service interface {
	Create(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Update(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Delete(ctx context.Context, uuids ...string) error
	Read(ctx context.Context, uuids ...string) ([]entity.UserData, error)
}

// Handler обработчик пользовательских данных.
type Handler struct {
	srv Service
}

// GetList возвращает список данных.
func (h *Handler) GetList(ctx context.Context) ([]entity.UserData, error) {
	return h.srv.Read(ctx)
}

// GetOne возвращает одну запись.
func (h *Handler) GetOne(ctx context.Context, uuid string) (*entity.UserData, error) {
	data, err := h.srv.Read(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

// ItemUpdate обновляет данные записи.
func (h *Handler) ItemUpdate(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	return h.srv.Update(ctx, data)
}

// ItemAdd добавляет запись.
func (h *Handler) ItemAdd(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	return h.srv.Create(ctx, data)
}

// ItemDelete удаляет запись.
func (h *Handler) ItemDelete(ctx context.Context, uuids ...string) error {
	return h.srv.Delete(ctx, uuids...)
}

// New конструктор.
func New(s Service) *Handler {
	return &Handler{
		srv: s,
	}
}

package userdata

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v "github.com/ktigay/goph-keeper/internal/validator"

	"github.com/ktigay/goph-keeper/internal/entity"
)

var (
	// ErrDataNotFound данные не найдены.
	ErrDataNotFound = errors.New("data not found")
	// ErrBadRequest неправильный запрос.
	ErrBadRequest = errors.New("bad request")
)

// Repository репозиторий.
//
//go:generate mockgen -destination=./mocks/mock_userdata.go -package=mocks github.com/ktigay/goph-keeper/internal/server/service/userdata Repository
type Repository interface {
	Create(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Update(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Delete(ctx context.Context, userUUID string, uuids ...string) error
	Read(ctx context.Context, userUUID string, uuids ...string) ([]entity.UserData, error)
}

// Service сервис.
type Service struct {
	repo Repository
}

// Create создаёт запись пользовательских данных.
func (s *Service) Create(ctx context.Context, userUUID string, data entity.UserData) (*entity.UserData, error) {
	if err := v.ValidateUserData(data); err != nil {
		return nil, ErrBadRequest
	}

	data.UserUUID = userUUID
	return s.repo.Create(ctx, data)
}

// Update обновляет запись пользовательских данных.
func (s *Service) Update(ctx context.Context, userUUID string, data entity.UserData) (*entity.UserData, error) {
	if err := v.ValidateUserData(data); err != nil {
		return nil, ErrBadRequest
	}

	data.UserUUID = userUUID
	return s.repo.Update(ctx, data)
}

// Delete удаляет записи пользовательских данных.
func (s *Service) Delete(ctx context.Context, userUUID string, uuids ...string) error {
	if len(uuids) == 0 {
		return ErrBadRequest
	}
	for _, u := range uuids {
		if err := uuid.Validate(u); err != nil {
			return ErrBadRequest
		}
	}
	return s.repo.Delete(ctx, userUUID, uuids...)
}

// Read возвращает записи пользовательских данных.
func (s *Service) Read(ctx context.Context, userUUID string, uuids ...string) ([]entity.UserData, error) {
	for _, u := range uuids {
		if err := uuid.Validate(u); err != nil {
			return nil, ErrBadRequest
		}
	}

	d, err := s.repo.Read(ctx, userUUID, uuids...)
	if err != nil {
		return nil, err
	}
	if len(d) == 0 {
		return nil, ErrDataNotFound
	}
	return d, nil
}

// New конструктор.
func New(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

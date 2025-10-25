package userdata

import (
	"context"

	"github.com/ktigay/goph-keeper/internal/entity"
	"github.com/ktigay/goph-keeper/internal/validator"
)

// Repository репозиторий.
//
//go:generate mockgen -destination=./mocks/mock_userdata.go -package=mocks github.com/ktigay/goph-keeper/internal/client/service/userdata Repository
type Repository interface {
	Sync(ctx context.Context, data []entity.UserData) error
	Create(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Update(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Replace(ctx context.Context, data entity.UserData) (*entity.UserData, error)
	Delete(ctx context.Context, uuids ...string) error
	Read(ctx context.Context, uuids ...string) ([]entity.UserData, error)
	ReadUnsynced(ctx context.Context) ([]entity.UserData, error)
}

// Service сервис пользовательских данных.
type Service struct {
	repo Repository
}

// Create создаёт запись пользовательских данных.
func (s *Service) Create(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	if err := validator.ValidateUserData(data); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, data)
}

// Update обновляет запись пользовательских данных.
func (s *Service) Update(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	if err := validator.ValidateUserData(data); err != nil {
		return nil, err
	}
	data.IsSynced = false
	return s.repo.Update(ctx, data)
}

// Delete удаляет запись пользовательских данных.
func (s *Service) Delete(ctx context.Context, uuids ...string) error {
	if len(uuids) == 0 {
		return nil
	}
	return s.repo.Delete(ctx, uuids...)
}

// Read читает записи пользовательских данных.
func (s *Service) Read(ctx context.Context, uuids ...string) ([]entity.UserData, error) {
	return s.repo.Read(ctx, uuids...)
}

// ReadOne читает запись пользовательских данных.
func (s *Service) ReadOne(ctx context.Context, uuid string) (*entity.UserData, error) {
	data, err := s.repo.Read(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

// New конструктор.
func New(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

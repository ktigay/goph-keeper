package sync

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ktigay/goph-keeper/internal/client/service/userdata"
	"github.com/ktigay/goph-keeper/internal/entity"
)

// Client клиент.
//
//go:generate mockgen -destination=./mocks/mock_client.go -package=mocks github.com/ktigay/goph-keeper/internal/client/service/sync Client
type Client interface {
	Create(ctx context.Context, d entity.UserData) (*entity.UserData, error)
	Update(ctx context.Context, d entity.UserData) (*entity.UserData, error)
	Read(ctx context.Context, uuid ...string) ([]entity.UserData, error)
	Delete(ctx context.Context, uuids ...string) error
}

// Service сервис синхронизации данных.
type Service struct {
	client Client
	repo   userdata.Repository
	logger *slog.Logger
}

// Initialize инициализирует пользовательские данные.
func (s *Service) Initialize(ctx context.Context) ([]entity.UserData, error) {
	data, err := s.readAllRemote(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Sync(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

// SyncToRemote синхронизирует локальные данные на сервер.
func (s *Service) SyncToRemote(ctx context.Context) ([]entity.UserData, error) {
	var (
		data []entity.UserData
		err  error
	)
	data, err = s.repo.ReadUnsynced(ctx)
	if err != nil {
		return nil, err
	}
	updated := make([]entity.UserData, len(data))
	for i := range data {
		var d *entity.UserData
		if data[i].IsNew {
			d, err = s.client.Create(ctx, data[i])
		} else {
			d, err = s.client.Update(ctx, data[i])
		}
		if err != nil {
			s.logger.Debug("error updating userdata", slog.String("uuid", data[i].UUID))
			return nil, err
		}

		d.IsSynced = true
		d.IsNew = false

		_, err = s.repo.Replace(ctx, *d)
		if err != nil {
			s.logger.Debug("error replacing userdata", slog.String("uuid", data[i].UUID))
			return nil, err
		}
		updated[i] = *d
	}
	return updated, nil
}

// SyncFromRemote синхронизирует данные с сервера.
func (s *Service) SyncFromRemote(ctx context.Context) error {
	data, err := s.readAllRemote(ctx)
	if err != nil {
		s.logger.Debug("error reading remote data", "err", err)
		return err
	}

	for _, d := range data {
		_, err = s.updateLocal(ctx, d)
		if err != nil {
			s.logger.Debug("error updating local data", "err", err)
			return err
		}
	}
	return nil
}

func (s *Service) updateLocal(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	old, err := s.readOneLocal(ctx, data.UUID)
	if err != nil {
		return nil, err
	}

	if old != nil && old.UpdatedAt.Compare(data.UpdatedAt) != 0 {
		if data.IsSynced {
			return nil, errors.New("data has been modified remotely")
		}
		return nil, errors.New("data has been modified since last sync")
	}

	data.IsSynced = true
	return s.repo.Replace(ctx, data)
}

func (s *Service) readAllRemote(ctx context.Context) ([]entity.UserData, error) {
	return s.client.Read(ctx)
}

func (s *Service) readOneLocal(ctx context.Context, uuid string) (*entity.UserData, error) {
	data, err := s.repo.Read(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("data not found")
	}
	return &data[0], nil
}

// New конструктор.
func New(c Client, r userdata.Repository, l *slog.Logger) *Service {
	return &Service{
		repo:   r,
		client: c,
		logger: l,
	}
}

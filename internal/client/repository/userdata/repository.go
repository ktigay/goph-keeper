package userdata

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ktigay/goph-keeper/internal/entity"
)

// Repository репозиторий.
type Repository struct {
	m            sync.Mutex
	data         map[string]entity.UserData
	initRequired bool
}

// Sync синхронизирует данные.
func (r *Repository) Sync(_ context.Context, data []entity.UserData) error {
	if !r.initRequired {
		return nil
	}

	r.data = make(map[string]entity.UserData, len(data))
	for _, d := range data {
		d.IsSynced = true
		r.data[d.UUID] = d
	}
	r.initRequired = false
	return nil
}

// Create создает данные.
func (r *Repository) Create(_ context.Context, data entity.UserData) (*entity.UserData, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if data.UUID == "" {
		data.UUID = uuid.New().String()
	}
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	r.data[data.UUID] = data
	return &data, nil
}

// Update обновляет данные.
func (r *Repository) Update(_ context.Context, data entity.UserData) (*entity.UserData, error) {
	r.m.Lock()
	defer r.m.Unlock()

	data.UpdatedAt = time.Now()
	r.data[data.UUID] = data
	return &data, nil
}

// Replace заменяет данные.
func (r *Repository) Replace(_ context.Context, data entity.UserData) (*entity.UserData, error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[data.UUID] = data
	return &data, nil
}

// Delete удаляет данные.
func (r *Repository) Delete(_ context.Context, uuids ...string) error {
	r.m.Lock()
	defer r.m.Unlock()

	for _, uid := range uuids {
		delete(r.data, uid)
	}
	return nil
}

// Read читает данные.
func (r *Repository) Read(_ context.Context, uuids ...string) ([]entity.UserData, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if len(uuids) == 0 {
		return sortByUpdated(slices.Collect(maps.Values(r.data)))
	}

	var data []entity.UserData
	for _, d := range r.data {
		if slices.Contains(uuids, d.UUID) {
			data = append(data, d)
		}
	}
	return sortByUpdated(data)
}

// ReadUnsynced читает не синхронизированные данные.
func (r *Repository) ReadUnsynced(_ context.Context) ([]entity.UserData, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var data []entity.UserData
	for _, d := range r.data {
		if !d.IsSynced {
			data = append(data, d)
		}
	}
	return sortByUpdated(data)
}

func sortByUpdated(data []entity.UserData) ([]entity.UserData, error) {
	slices.SortFunc(data, func(a, b entity.UserData) int {
		return b.UpdatedAt.Compare(a.UpdatedAt)
	})
	return data, nil
}

// New конструктор.
func New() *Repository {
	return &Repository{
		data:         make(map[string]entity.UserData),
		initRequired: true,
	}
}

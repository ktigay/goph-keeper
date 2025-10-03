package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/ktigay/goph-keeper/internal/server/db"
	"github.com/ktigay/goph-keeper/internal/server/entity"
)

var (
	insertQuery = `
		INSERT INTO "user" ("login", "password")
			VALUES ($1, $2) 
		RETURNING "uuid", "login", "password", "created_at", "updated_at"`

	selectByLoginQuery = `
		SELECT "uuid", "login", "password", "created_at", "updated_at"
		FROM "user"
		WHERE "login" = $1
	`
)

// Repository репозиторий.
type Repository struct {
	db     db.ConnWrapper
	logger *slog.Logger
}

// Create создаёт пользователя.
func (r *Repository) Create(ctx context.Context, login, password string) (*entity.User, error) {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	return r.queryRow(c, insertQuery, login, password)
}

// Read возвращает пользователя.
func (r *Repository) Read(ctx context.Context, login string) (*entity.User, error) {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	e, err := r.queryRow(c, selectByLoginQuery, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

func (r *Repository) queryRow(ctx context.Context, query string, args ...any) (*entity.User, error) {
	var (
		u   entity.User
		err error
	)

	if err = r.fullScan(
		r.db.Connection(ctx).QueryRow(ctx, query, args...),
		&u,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) fullScan(row pgx.Row, ud *entity.User) error {
	return row.Scan(
		&ud.UUID,
		&ud.Login,
		&ud.Password,
		&ud.UpdatedAt,
		&ud.CreatedAt,
	)
}

// New Конструктор.
func New(db db.ConnWrapper, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

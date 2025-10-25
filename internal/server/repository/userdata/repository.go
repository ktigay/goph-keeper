package userdata

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/ktigay/goph-keeper/internal/entity"
	"github.com/ktigay/goph-keeper/internal/server/db"
)

var (
	insertQuery = `
		INSERT INTO "user_data" ("user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at")
			VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING "uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at"`

	insertQueryWithUUID = `
		INSERT INTO "user_data" ("uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING "uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at"`

	updateQuery = `
		UPDATE "user_data" 
		SET 
		    "title" = $1, "type" = $2, "data" = $3, "metadata" = $4, "updated_at" = $5
		WHERE "uuid" = $6 AND "user_uuid" = $7
		RETURNING "uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at"`

	deleteQuery = `
		DELETE FROM "user_data"
		WHERE "user_uuid" = $1 AND "uuid" = ANY($2::uuid[])
	`

	selectByUserQuery = `
		SELECT "uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at"
		FROM "user_data"
		WHERE "user_uuid" = $1
	`
	selectByUserAndIDQuery = `
		SELECT "uuid", "user_uuid", "title", "type", "data", "metadata", "created_at", "updated_at"
		FROM "user_data"
		WHERE "user_uuid" = $1 AND "uuid" = ANY($2::uuid[])
	`
)

// Repository репозиторий.
type Repository struct {
	db     db.ConnWrapper
	logger *slog.Logger
}

// Create создаёт запись пользовательских данных.
func (r *Repository) Create(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	if data.UUID != "" {
		return r.queryRow(c, insertQueryWithUUID, data.UUID, data.UserUUID, data.Title, data.Type, data.Data, data.MetaData, data.CreatedAt, data.UpdatedAt)
	}
	return r.queryRow(c, insertQuery, data.UserUUID, data.Title, data.Type, data.Data, data.MetaData, data.CreatedAt, data.UpdatedAt)
}

// Update обновляет запись пользовательских данных.
func (r *Repository) Update(ctx context.Context, data entity.UserData) (*entity.UserData, error) {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	return r.queryRow(c, updateQuery, data.Title, data.Type, data.Data, data.MetaData, data.UpdatedAt, data.UUID, data.UserUUID)
}

// Delete удаляет запись пользовательских данных.
func (r *Repository) Delete(ctx context.Context, userUUID string, uuids ...string) error {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	_, err := r.db.Connection(ctx).Exec(c, deleteQuery, userUUID, "{"+strings.Join(uuids, ",")+"}")
	return err
}

// Read читает записи пользовательских данных.
func (r *Repository) Read(ctx context.Context, userUUID string, uuids ...string) ([]entity.UserData, error) {
	c, cancel := context.WithTimeout(ctx, db.RequestTimeout)
	defer cancel()

	var (
		rows pgx.Rows
		err  error
	)
	if len(uuids) == 0 {
		rows, err = r.db.Connection(ctx).Query(c, selectByUserQuery, userUUID)
	} else {
		rows, err = r.db.Connection(ctx).Query(c, selectByUserAndIDQuery, userUUID, "{"+strings.Join(uuids, ",")+"}")
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	d := make([]entity.UserData, 0)
	for rows.Next() {
		var order entity.UserData
		err = r.fullScan(
			rows,
			&order,
		)
		if err != nil {
			return nil, err
		}

		d = append(d, order)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return d, nil
}

func (r *Repository) queryRow(ctx context.Context, query string, args ...any) (*entity.UserData, error) {
	var (
		ud  entity.UserData
		err error
	)

	if err = r.fullScan(
		r.db.Connection(ctx).QueryRow(ctx, query, args...),
		&ud,
	); err != nil {
		return nil, err
	}
	return &ud, nil
}

func (r *Repository) fullScan(row pgx.Row, ud *entity.UserData) error {
	return row.Scan(
		&ud.UUID,
		&ud.UserUUID,
		&ud.Title,
		&ud.Type,
		&ud.Data,
		&ud.MetaData,
		&ud.CreatedAt,
		&ud.UpdatedAt,
	)
}

// New Конструктор.
func New(db db.ConnWrapper, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

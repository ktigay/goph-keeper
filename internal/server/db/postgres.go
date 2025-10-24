package db

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	// RequestTimeout таймаут запроса.
	RequestTimeout = 1 * time.Second
	connTimeout    = 1 * time.Second
	schemaPath     = "./internal/server/db/schema.sql"
)

// NewPgxPool Новый пул соединений к БД.
func NewPgxPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, connTimeout)
	defer cancel()
	if err = pool.Ping(c); err != nil {
		return nil, err
	}

	return pool, nil
}

// CreateSchema создает структуру БД.
func CreateSchema(ctx context.Context, pool *pgxpool.Pool) error {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	file, err := os.OpenFile(schemaPath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Print("File close error:", "err", err)
		}
	}()

	var buff []byte
	if buff, err = io.ReadAll(file); err != nil {
		return err
	}

	_, err = pool.Exec(c, string(buff))
	if err != nil {
		return err
	}

	return nil
}

// ConnWrapper Интерфейс обертки для коннекта БД.
type ConnWrapper interface {
	Connection(ctx context.Context) Conn
}

// Conn Интерфейс для работы с БД.
type Conn interface {
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

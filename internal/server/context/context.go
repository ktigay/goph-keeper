package context

import (
	"context"
	"errors"

	"github.com/ktigay/goph-keeper/internal/entity"
)

type contextKey int

const (
	ctxKey contextKey = iota
)

// NewContextWithIdentity возвращает контекст с идентификационными данными.
func NewContextWithIdentity(ctx context.Context, identity entity.Identity) context.Context {
	return context.WithValue(ctx, ctxKey, identity)
}

// IdentityFromContext извлекает идентификационные данные из контекста.
func IdentityFromContext(ctx context.Context) (*entity.Identity, error) {
	e, ok := ctx.Value(ctxKey).(entity.Identity)
	if !ok {
		return nil, errors.New("no identity found in context")
	}

	return &e, nil
}

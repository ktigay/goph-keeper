package auth

import "context"

// Repository репозиторий.
type Repository struct {
	jwt string
}

// SetJWT сохраняет JWT токен в репозиторий.
func (r *Repository) SetJWT(_ context.Context, jwt string) error {
	r.jwt = jwt
	return nil
}

// GetJWT возвращает JWT токен.
func (r *Repository) GetJWT(_ context.Context) (string, error) {
	return r.jwt, nil
}

// New конструктор.
func New() *Repository {
	return &Repository{}
}

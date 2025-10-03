package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/ktigay/goph-keeper/internal/server/entity"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

var (
	// ErrLoginOrPwdEmpty логин или пароль пустые.
	ErrLoginOrPwdEmpty = errors.New("login or password is empty")
	// ErrUserNotFound пользователь не найден.
	ErrUserNotFound = errors.New("user not found")
	// ErrWrongPassword неправильный пароль.
	ErrWrongPassword = errors.New("wrong password")
)

// Repository репозиторий.
//
//go:generate mockgen -destination=./mocks/mock_auth.go -package=mocks github.com/ktigay/goph-keeper/internal/server/service/auth Repository
type Repository interface {
	Create(ctx context.Context, login, password string) (*entity.User, error)
	Read(ctx context.Context, login string) (*entity.User, error)
}

// Service сервис.
type Service struct {
	repo Repository
}

// Register регистрирует пользователя.
func (s *Service) Register(ctx context.Context, login, password string) (*entity.User, error) {
	var (
		newUsr       *entity.User
		hashedPasswd []byte
		err          error
	)

	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)
	if login == "" || password == "" {
		return nil, ErrLoginOrPwdEmpty
	}

	hashedPasswd, err = bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, err
	}

	if newUsr, err = s.repo.Create(ctx, login, string(hashedPasswd)); err != nil {
		return nil, err
	}

	return newUsr, err
}

// Login авторизует пользователя.
func (s *Service) Login(ctx context.Context, login, password string) (*entity.User, error) {
	var (
		usr *entity.User
		err error
	)

	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)
	if login == "" || password == "" {
		return nil, ErrLoginOrPwdEmpty
	}

	if usr, err = s.repo.Read(ctx, login); err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrWrongPassword
		}

		return nil, err
	}

	return usr, nil
}

// New конструктор.
func New(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

package security

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	subjectKey = "sub"
)

// ErrMissingIdentity ошибка - нет идентификатора пользователя.
var ErrMissingIdentity = fmt.Errorf("missing identity")

// JWTWrapper Обертка для JWT.
type JWTWrapper[T any] struct {
	secret []byte
}

// GenerateToken генерирует токен JWT.
func (j *JWTWrapper[T]) GenerateToken(payload T) (string, error) {
	v, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		subjectKey: string(v),
	})
	return token.SignedString(j.secret)
}

// ParseToken парсит токен JWT.
func (j *JWTWrapper[T]) ParseToken(s string) (*T, error) {
	s = strings.TrimPrefix(s, "Bearer ")
	if len(s) == 0 {
		return nil, ErrMissingIdentity
	}

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	var (
		claims jwt.MapClaims
		subj   string
		ok     bool
	)
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, ErrMissingIdentity
	}

	if subj, err = claims.GetSubject(); err != nil {
		return nil, ErrMissingIdentity
	}

	var resp T
	if err = json.Unmarshal([]byte(subj), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// NewJWTWrapper Конструктор.
func NewJWTWrapper[T any](secret string) *JWTWrapper[T] {
	return &JWTWrapper[T]{
		secret: []byte(secret),
	}
}

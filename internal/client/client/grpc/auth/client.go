package auth

import (
	"context"

	"github.com/ktigay/goph-keeper/internal/client/entity"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
)

// Client grpc клиент.
type Client struct {
	conn auth.AuthServiceClient
}

// Login авторизует пользователя.
func (c *Client) Login(ctx context.Context, data entity.Credentials) (string, error) {
	req := &auth.LoginRequest{
		Login:    data.Login,
		Password: data.Password,
	}

	resp, err := c.conn.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

// Register регистрирует пользователя.
func (c *Client) Register(ctx context.Context, data entity.Credentials) (string, error) {
	req := &auth.RegisterRequest{
		Login:    data.Login,
		Password: data.Password,
	}

	resp, err := c.conn.Register(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.UserUuid, nil
}

// New конструктор.
func New(c auth.AuthServiceClient) *Client {
	return &Client{c}
}

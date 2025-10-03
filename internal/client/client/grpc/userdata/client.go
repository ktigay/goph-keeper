package userdata

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/data/mapper"
	"github.com/ktigay/goph-keeper/internal/entity"
)

// Client клиент.
type Client struct {
	conn data.UserDataServiceClient
}

// Create создает запись пользовательских данных.
func (c *Client) Create(ctx context.Context, d entity.UserData) (*entity.UserData, error) {
	req := data.CreateUserDataItemRequest{
		Item: mapper.MapEntityToItem(d),
	}
	resp, err := c.conn.CreateUserDataItem(ctx, &req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Item == nil {
		return nil, fmt.Errorf("response is nil")
	}
	d = mapper.MapItemToEntity(resp.GetItem(), "")
	d.IsSynced = true
	d.IsNew = false
	return &d, nil
}

// Read читает записи пользовательских данных.
func (c *Client) Read(ctx context.Context, uuid ...string) ([]entity.UserData, error) {
	req := data.GetUserDataItemsRequest{ItemUuids: uuid}
	resp, err := c.conn.GetUserDataItems(ctx, &req)

	code := status.Code(err)
	if code != codes.OK && code != codes.NotFound {
		return nil, err
	}
	if resp == nil {
		return []entity.UserData{}, nil
	}

	d := make([]entity.UserData, len(resp.Items))
	for i := range resp.Items {
		d[i] = mapper.MapItemToEntity(resp.Items[i], "")
	}
	return d, nil
}

// Update обновляет запись пользовательских данных.
func (c *Client) Update(ctx context.Context, d entity.UserData) (*entity.UserData, error) {
	req := data.UpdateUserDataItemRequest{
		Item: mapper.MapEntityToItem(d),
	}
	resp, err := c.conn.UpdateUserDataItem(ctx, &req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Item == nil {
		return nil, fmt.Errorf("response is nil")
	}
	d = mapper.MapItemToEntity(resp.GetItem(), "")
	return &d, nil
}

// Delete удаляет запись пользовательских данных.
func (c *Client) Delete(ctx context.Context, uuid ...string) error {
	req := data.DeleteUserDataItemsRequest{ItemUuids: uuid}
	_, err := c.conn.DeleteUserDataItems(ctx, &req)
	return err
}

// New конструктор.
func New(conn data.UserDataServiceClient) *Client {
	return &Client{
		conn: conn,
	}
}

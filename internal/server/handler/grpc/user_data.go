package grpc

import (
	"context"
	"errors"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/data/mapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/entity"
	c "github.com/ktigay/goph-keeper/internal/server/context"
	"github.com/ktigay/goph-keeper/internal/server/service/userdata"
)

// UserDataService сервис пользовательских данных.
//
//go:generate mockgen -destination=./mocks/mock_userdata.go -package=mocks github.com/ktigay/goph-keeper/internal/server/handler/grpc UserDataService
type UserDataService interface {
	Create(ctx context.Context, userUUID string, data entity.UserData) (*entity.UserData, error)
	Update(ctx context.Context, userUUID string, data entity.UserData) (*entity.UserData, error)
	Delete(ctx context.Context, userUUID string, uuids ...string) error
	Read(ctx context.Context, userUUID string, uuids ...string) ([]entity.UserData, error)
}

// UserDataHandler обработчик пользовательских данных.
type UserDataHandler struct {
	data.UnimplementedUserDataServiceServer
	srv UserDataService
}

// CreateUserDataItem создает запись пользовательских данных.
func (u *UserDataHandler) CreateUserDataItem(ctx context.Context, request *data.CreateUserDataItemRequest) (*data.CreateUserDataItemResponse, error) {
	var (
		identity *entity.Identity
		d        *entity.UserData
		err      error
	)

	if identity, err = c.IdentityFromContext(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization required: %v", err)
	}

	d, err = u.srv.Create(ctx, identity.UUID, mapper.MapItemToEntity(request.Item, identity.UUID))
	if err != nil {
		return nil, status.Errorf(mapErrorToCode(err), "%v", err)
	}

	return &data.CreateUserDataItemResponse{
		Item: mapper.MapEntityToItem(*d),
	}, nil
}

// UpdateUserDataItem обновляет запись пользовательских данных.
func (u *UserDataHandler) UpdateUserDataItem(ctx context.Context, request *data.UpdateUserDataItemRequest) (*data.UpdateUserDataItemResponse, error) {
	var (
		identity *entity.Identity
		d        *entity.UserData
		err      error
	)

	if identity, err = c.IdentityFromContext(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization required: %v", err)
	}

	d, err = u.srv.Update(ctx, identity.UUID, mapper.MapItemToEntity(request.Item, identity.UUID))
	if err != nil {
		return nil, status.Errorf(mapErrorToCode(err), "%v", err)
	}

	return &data.UpdateUserDataItemResponse{
		Item: mapper.MapEntityToItem(*d),
	}, nil
}

// GetUserDataItem возвращает запись пользовательских данных.
func (u *UserDataHandler) GetUserDataItem(ctx context.Context, request *data.GetUserDataItemRequest) (*data.GetUserDataItemResponse, error) {
	var (
		identity *entity.Identity
		d        []entity.UserData
		err      error
	)

	if identity, err = c.IdentityFromContext(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization required: %v", err)
	}

	d, err = u.srv.Read(ctx, identity.UUID, request.ItemUuid)
	if err != nil {
		return nil, status.Errorf(mapErrorToCode(err), "%v", err)
	}

	if len(d) == 0 {
		return nil, status.Errorf(codes.NotFound, "item uuid %s not found", request.ItemUuid)
	}

	return &data.GetUserDataItemResponse{
		Item: mapper.MapEntityToItem(d[0]),
	}, nil
}

// GetUserDataItems возвращает записи пользовательских данных.
func (u *UserDataHandler) GetUserDataItems(ctx context.Context, request *data.GetUserDataItemsRequest) (*data.GetUserDataItemsResponse, error) {
	var (
		identity *entity.Identity
		d        []entity.UserData
		err      error
	)

	if identity, err = c.IdentityFromContext(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization required: %v", err)
	}

	d, err = u.srv.Read(ctx, identity.UUID, request.ItemUuids...)
	if err != nil {
		return nil, status.Errorf(mapErrorToCode(err), "%v", err)
	}

	items := make([]*data.UserDataItem, 0, len(d))
	for _, v := range d {
		items = append(items, mapper.MapEntityToItem(v))
	}

	return &data.GetUserDataItemsResponse{
		Items: items,
	}, nil
}

// DeleteUserDataItems удаляет записи пользовательских данных.
func (u *UserDataHandler) DeleteUserDataItems(ctx context.Context, request *data.DeleteUserDataItemsRequest) (*emptypb.Empty, error) {
	var (
		identity *entity.Identity
		err      error
	)

	if identity, err = c.IdentityFromContext(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization required: %v", err)
	}

	if err = u.srv.Delete(ctx, identity.UUID, request.ItemUuids...); err != nil {
		return nil, status.Errorf(mapErrorToCode(err), "%v", err)
	}
	return &emptypb.Empty{}, nil
}

// NewUserDataHandler конструктор.
func NewUserDataHandler(s UserDataService) *UserDataHandler {
	return &UserDataHandler{
		srv: s,
	}
}

func mapErrorToCode(err error) codes.Code {
	switch true {
	case errors.Is(err, userdata.ErrDataNotFound):
		return codes.NotFound
	case errors.Is(err, userdata.ErrBadRequest):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

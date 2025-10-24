package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/data/mapper"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/entity"
	c "github.com/ktigay/goph-keeper/internal/server/context"
	"github.com/ktigay/goph-keeper/internal/server/handler/grpc/mocks"
)

func TestUserDataHandler_CreateUserDataItem(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) UserDataService
	}
	type args struct {
		ctx     context.Context
		request *data.CreateUserDataItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.CreateUserDataItemResponse
		wantErr bool
	}{
		{
			name: "CreateUserDataItem_Authorization_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return srv
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &data.CreateUserDataItemRequest{},
			},
			wantErr: true,
		},
		{
			name: "CreateUserDataItem_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&entity.UserData{
						UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
					}, nil)
					return srv
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return c.NewContextWithIdentity(ctx, entity.Identity{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					})
				}(),
				request: &data.CreateUserDataItemRequest{
					Item: &data.UserDataItem{},
				},
			},
			wantErr: false,
			want: func() *data.CreateUserDataItemResponse {
				return &data.CreateUserDataItemResponse{
					Item: mapper.MapEntityToItem(entity.UserData{
						UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
					}),
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := &UserDataHandler{
				srv: tt.fields.srv(ctrl),
			}
			got, err := u.CreateUserDataItem(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserDataItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUserDataItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataHandler_DeleteUserDataItems(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) UserDataService
	}
	type args struct {
		ctx     context.Context
		request *data.DeleteUserDataItemsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "DeleteUserDataItems_Authorization_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return srv
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &data.DeleteUserDataItemsRequest{},
			},
			wantErr: true,
		},
		{
			name: "DeleteUserDataItems_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
					return srv
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return c.NewContextWithIdentity(ctx, entity.Identity{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					})
				}(),
				request: &data.DeleteUserDataItemsRequest{
					ItemUuids: []string{
						"9b89b845-164b-498d-bc0e-f197fec9008a",
						"5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
					},
				},
			},
			wantErr: false,
			want:    &emptypb.Empty{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := &UserDataHandler{
				srv: tt.fields.srv(ctrl),
			}
			got, err := u.DeleteUserDataItems(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserDataItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUserDataItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataHandler_GetUserDataItem(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) UserDataService
	}
	type args struct {
		ctx     context.Context
		request *data.GetUserDataItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.GetUserDataItemResponse
		wantErr bool
	}{
		{
			name: "GetUserDataItem_Authorization_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return srv
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &data.GetUserDataItemRequest{},
			},
			wantErr: true,
		},
		{
			name: "GetUserDataItem_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return([]entity.UserData{
						{
							UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
						},
					}, nil)
					return srv
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return c.NewContextWithIdentity(ctx, entity.Identity{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					})
				}(),
				request: &data.GetUserDataItemRequest{
					ItemUuid: "10c33409-d8cc-4673-9bfc-3182a894acd4",
				},
			},
			wantErr: false,
			want: func() *data.GetUserDataItemResponse {
				return &data.GetUserDataItemResponse{
					Item: mapper.MapEntityToItem(entity.UserData{
						UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
					}),
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := &UserDataHandler{
				srv: tt.fields.srv(ctrl),
			}
			got, err := u.GetUserDataItem(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDataItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserDataItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataHandler_GetUserDataItems(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) UserDataService
	}
	type args struct {
		ctx     context.Context
		request *data.GetUserDataItemsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.GetUserDataItemsResponse
		wantErr bool
	}{
		{
			name: "GetUserDataItems_Authorization_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return srv
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &data.GetUserDataItemsRequest{},
			},
			wantErr: true,
		},
		{
			name: "GetUserDataItems_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return([]entity.UserData{
						{
							UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
						},
					}, nil)
					return srv
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return c.NewContextWithIdentity(ctx, entity.Identity{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					})
				}(),
				request: &data.GetUserDataItemsRequest{
					ItemUuids: []string{
						"10c33409-d8cc-4673-9bfc-3182a894acd4",
					},
				},
			},
			wantErr: false,
			want: func() *data.GetUserDataItemsResponse {
				return &data.GetUserDataItemsResponse{
					Items: []*data.UserDataItem{
						mapper.MapEntityToItem(entity.UserData{
							UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
						}),
					},
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := &UserDataHandler{
				srv: tt.fields.srv(ctrl),
			}
			got, err := u.GetUserDataItems(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDataItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserDataItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataHandler_UpdateUserDataItem(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) UserDataService
	}
	type args struct {
		ctx     context.Context
		request *data.UpdateUserDataItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.UpdateUserDataItemResponse
		wantErr bool
	}{
		{
			name: "UpdateUserDataItem_Authorization_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return srv
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &data.UpdateUserDataItemRequest{},
			},
			wantErr: true,
		},
		{
			name: "UpdateUserDataItem_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) UserDataService {
					srv := mocks.NewMockUserDataService(ctrl)
					srv.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&entity.UserData{
						UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
					}, nil)
					return srv
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return c.NewContextWithIdentity(ctx, entity.Identity{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					})
				}(),
				request: &data.UpdateUserDataItemRequest{
					Item: &data.UserDataItem{},
				},
			},
			wantErr: false,
			want: func() *data.UpdateUserDataItemResponse {
				return &data.UpdateUserDataItemResponse{
					Item: mapper.MapEntityToItem(entity.UserData{
						UUID: "10c33409-d8cc-4673-9bfc-3182a894acd4",
					}),
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := &UserDataHandler{
				srv: tt.fields.srv(ctrl),
			}
			got, err := u.UpdateUserDataItem(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserDataItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUserDataItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

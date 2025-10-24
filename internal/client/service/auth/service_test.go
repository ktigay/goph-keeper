package auth

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ktigay/goph-keeper/internal/client/entity"
	"github.com/ktigay/goph-keeper/internal/client/service/auth/mocks"
	"github.com/ktigay/goph-keeper/internal/log"
)

func TestService_GetJWT(t *testing.T) {
	type fields struct {
		client func(controller *gomock.Controller) Client
		repo   func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().GetJWT(gomock.Any()).Times(1).Return("jwt-token", nil)
					return repo
				},
				client: func(controller *gomock.Controller) Client {
					return mocks.NewMockClient(controller)
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: "jwt-token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			got, err := s.GetJWT(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsAuthorized(t *testing.T) {
	type fields struct {
		client func(controller *gomock.Controller) Client
		repo   func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Authorized",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().GetJWT(gomock.Any()).Times(1).Return("jwt-token", nil)
					return repo
				},
				client: func(controller *gomock.Controller) Client {
					return mocks.NewMockClient(controller)
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
		{
			name: "Not_Authorized",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().GetJWT(gomock.Any()).Times(1).Return("", nil)
					return repo
				},
				client: func(controller *gomock.Controller) Client {
					return mocks.NewMockClient(controller)
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			if got := s.IsAuthorized(tt.args.ctx); got != tt.want {
				t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	type fields struct {
		client func(controller *gomock.Controller) Client
		repo   func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx  context.Context
		data entity.Credentials
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Login_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().SetJWT(gomock.Eq(context.Background()), gomock.Eq("jwt-token")).Times(1).Return(nil)
					return repo
				},
				client: func(controller *gomock.Controller) Client {
					c := mocks.NewMockClient(controller)
					c.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return("jwt-token", nil)
					return c
				},
			},
			args: args{
				ctx: context.Background(),
				data: entity.Credentials{
					Login:    "login",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "Login_Failed_Empty_Login",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().SetJWT(gomock.Eq(context.Background()), gomock.Eq("jwt-token")).Times(0)
					return repo
				},
				client: func(controller *gomock.Controller) Client {
					c := mocks.NewMockClient(controller)
					c.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
					return c
				},
			},
			args: args{
				ctx: context.Background(),
				data: entity.Credentials{
					Login:    " ",
					Password: "password",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			if err := s.Login(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	type fields struct {
		client func(controller *gomock.Controller) Client
		repo   func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx  context.Context
		data entity.Credentials
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Register_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					return mocks.NewMockRepository(ctrl)
				},
				client: func(controller *gomock.Controller) Client {
					c := mocks.NewMockClient(controller)
					c.EXPECT().Register(gomock.Eq(context.Background()), gomock.Eq(entity.Credentials{
						Login:    "login",
						Password: "password",
					})).Times(1).Return("9b89b845-164b-498d-bc0e-f197fec9008a", nil)
					return c
				},
			},
			args: args{
				ctx: context.Background(),
				data: entity.Credentials{
					Login:    "login",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "Register_Failed_Empty_Password",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					return mocks.NewMockRepository(ctrl)
				},
				client: func(controller *gomock.Controller) Client {
					c := mocks.NewMockClient(controller)
					c.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
					return c
				},
			},
			args: args{
				ctx: context.Background(),
				data: entity.Credentials{
					Login:    "login",
					Password: " ",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			if err := s.Register(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

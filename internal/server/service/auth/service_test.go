package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ktigay/goph-keeper/internal/server/entity"
	"github.com/ktigay/goph-keeper/internal/server/service/auth/mocks"
)

func TestService_Login(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "Login_success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any()).Times(1).Return(
						&entity.User{
							Login:    "test",
							Password: "$2a$10$dvLRmGJ8HMdgLTHeklRXNelmcoHb82y5lGIX3JJl4tawa9/mNEkba",
						}, nil)
					return repo
				},
			},
			args: args{
				login:    "test",
				password: "test",
			},
			want: &entity.User{
				Login:    "test",
				Password: "$2a$10$dvLRmGJ8HMdgLTHeklRXNelmcoHb82y5lGIX3JJl4tawa9/mNEkba",
			},
			wantErr: false,
		},
		{
			name: "Login_ErrUserNotFound",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any()).Times(1).Return(
						nil, nil)
					return repo
				},
			},
			args: args{
				login:    "test",
				password: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Login_ErrLoginOrPwdEmpty",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				login:    "",
				password: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				repo: tt.fields.repo(ctrl),
			}
			got, err := s.Login(context.Background(), tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "Register_success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(
						&entity.User{
							Login: "test",
						}, nil)
					return repo
				},
			},
			args: args{
				login:    "test",
				password: "test",
			},
			want: &entity.User{
				Login: "test",
			},
			wantErr: false,
		},
		{
			name: "Register_ErrLoginOrPwdEmpty",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				login:    "test",
				password: " ",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Register_ErrLoginOrPwdEmpty_#2",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				login:    "",
				password: " test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				repo: tt.fields.repo(ctrl),
			}
			got, err := s.Register(context.Background(), tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package grpc

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
	"github.com/ktigay/goph-keeper/internal/server/entity"
	"github.com/ktigay/goph-keeper/internal/server/handler/grpc/mocks"
)

func TestAuthHandler_Login(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) AuthService
		jwt func(ctrl *gomock.Controller) JWTWrapper
	}
	type args struct {
		ctx context.Context
		req *auth.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *auth.LoginResponse
		wantErr bool
	}{
		{
			name: "Login_Success",
			fields: fields{
				srv: func(ctrl *gomock.Controller) AuthService {
					srv := mocks.NewMockAuthService(ctrl)
					srv.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&entity.User{
						UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
					}, nil)
					return srv
				},
				jwt: func(ctrl *gomock.Controller) JWTWrapper {
					w := mocks.NewMockJWTWrapper(ctrl)
					w.EXPECT().GenerateToken(gomock.Any()).Times(1).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA", nil)
					return w
				},
			},
			want: &auth.LoginResponse{
				Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA",
			},
			wantErr: false,
		},
		{
			name: "Login_Failed",
			fields: fields{
				srv: func(ctrl *gomock.Controller) AuthService {
					srv := mocks.NewMockAuthService(ctrl)
					srv.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, fmt.Errorf("some error"))
					return srv
				},
				jwt: func(ctrl *gomock.Controller) JWTWrapper {
					w := mocks.NewMockJWTWrapper(ctrl)
					w.EXPECT().GenerateToken(gomock.Any()).Times(0)
					return w
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			a := &AuthHandler{
				srv: tt.fields.srv(ctrl),
				jwt: tt.fields.jwt(ctrl),
			}
			got, err := a.Login(tt.args.ctx, tt.args.req)
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

func TestAuthHandler_Register(t *testing.T) {
	type fields struct {
		srv func(ctrl *gomock.Controller) AuthService
		jwt JWTWrapper
	}
	type args struct {
		ctx context.Context
		req *auth.RegisterRequest
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
				srv: func(ctrl *gomock.Controller) AuthService {
					srv := mocks.NewMockAuthService(ctrl)
					srv.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&entity.User{}, nil)
					return srv
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			a := &AuthHandler{
				srv: tt.fields.srv(ctrl),
				jwt: tt.fields.jwt,
			}
			_, err := a.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

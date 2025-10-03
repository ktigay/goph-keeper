package interceptor

import (
	"context"
	"reflect"
	"testing"

	"github.com/ktigay/goph-keeper/internal/entity"
	c "github.com/ktigay/goph-keeper/internal/server/context"
	"github.com/ktigay/goph-keeper/internal/server/security"
	"google.golang.org/grpc/metadata"
)

func TestAuth_authorization(t *testing.T) {
	type fields struct {
		jwt        JWTWrapper
		accessList map[string]bool
	}
	type args struct {
		ctx    context.Context
		method string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    context.Context
		wantErr bool
	}{
		{
			name: "Access_Method_Not_In_List_Error",
			fields: fields{
				jwt: security.NewJWTWrapper[entity.Identity]("secret"),
				accessList: map[string]bool{
					"method1": true,
					"method2": true,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: "method3",
			},
			wantErr: true,
		},
		{
			name: "Access_Method_Without_Authorization_Success",
			fields: fields{
				jwt: security.NewJWTWrapper[entity.Identity]("secret"),
				accessList: map[string]bool{
					"method1": false,
					"method2": true,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: "method1",
			},
			wantErr: false,
			want:    context.Background(),
		},
		{
			name: "Access_Method_With_Authorization_Empty_Metadata_error",
			fields: fields{
				jwt: security.NewJWTWrapper[entity.Identity]("secret"),
				accessList: map[string]bool{
					"method1": true,
					"method2": true,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: "method1",
			},
			wantErr: true,
		},
		{
			name: "Access_Method_With_Authorization_Parse_Token_Success",
			fields: fields{
				jwt: security.NewJWTWrapper[entity.Identity]("secret"),
				accessList: map[string]bool{
					"method1": true,
					"method2": true,
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					m := map[string]string{
						"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA",
					}
					return metadata.NewIncomingContext(ctx, metadata.New(m))
				}(),
				method: "method1",
			},
			wantErr: false,
			want: func() context.Context {
				ctx := context.Background()
				m := map[string]string{
					"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA",
				}
				ctx = metadata.NewIncomingContext(ctx, metadata.New(m))
				return c.NewContextWithIdentity(ctx, entity.Identity{
					UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
				})
			}(),
		},
		{
			name: "Access_Method_With_Authorization_Parse_Token_Failed",
			fields: fields{
				jwt: security.NewJWTWrapper[entity.Identity]("secret"),
				accessList: map[string]bool{
					"method1": true,
					"method2": true,
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					m := map[string]string{
						"authorization": "Bearer eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA",
					}
					return metadata.NewIncomingContext(ctx, metadata.New(m))
				}(),
				method: "method1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Auth{
				jwt:        tt.fields.jwt,
				accessList: tt.fields.accessList,
			}
			got, err := i.authorization(tt.args.ctx, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("authorization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authorization() got = %v, want %v", got, tt.want)
			}
		})
	}
}

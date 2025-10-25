package security

import (
	"reflect"
	"testing"

	"github.com/ktigay/goph-keeper/internal/entity"
)

func TestJWTWrapper_GenerateToken(t *testing.T) {
	type args[T any] struct {
		payload T
	}
	type testCase[T any] struct {
		name    string
		j       JWTWrapper[T]
		args    args[T]
		want    string
		wantErr bool
	}
	tests := []testCase[entity.Identity]{
		{
			name: "Generate_Token_success",
			j:    *NewJWTWrapper[entity.Identity]("secret"),
			args: args[entity.Identity]{
				payload: entity.Identity{
					UUID: "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
				},
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCI0ZDhkZTlkYy1iM2IzLTRjNDUtYjcxYi0xODlmYjQxODM3ZWFcIn0ifQ.UJoVLAx7dzjXFL5cEI6RTZwRyj9GZl6la5T3oOQVxw4",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GenerateToken(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTWrapper_ParseToken(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[T any] struct {
		name    string
		j       JWTWrapper[T]
		args    args
		want    *T
		wantErr bool
	}
	tests := []testCase[entity.Identity]{
		{
			name: "Parse_Token_success",
			j:    *NewJWTWrapper[entity.Identity]("secret"),
			args: args{
				s: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJVVUlEXCI6XCIzM2IwNjYxOS0xZWU3LTNkYjUtODI3ZC0wZGM4NWRmMWY3NTlcIn0ifQ.oaNdHD11OLp25Y0LIvMiAmzJsiNLiS-kQGRUXiQk-zA",
			},
			want: &entity.Identity{
				UUID: "33b06619-1ee7-3db5-827d-0dc85df1f759",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.ParseToken(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

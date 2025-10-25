package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		envs map[string]string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Check_Priority_Flags_Redefine_Envs",
			args: args{
				envs: map[string]string{
					"GRPC_ADDRESS": ":18090",
					"LOG_LEVEL":    "error",
					"DATABASE_URI": "postgres://postgres:postgres@localhost:15429/db?sslmode=disable",
					"JWT_SECRET":   "secret_secret",
					"CONFIG":       "",
				},
				args: []string{
					"-a=:18080",
					"-l=error",
					"-d=postgres://postgres:postgres@227.0.0.1:15429/db?sslmode=disable",
					"-s=secret_secret_priority",
				},
			},
			want: &Config{
				ServerGRPCHost: ":18080",
				LogLevel:       "error",
				DatabaseDSN:    "postgres://postgres:postgres@227.0.0.1:15429/db?sslmode=disable",
				AuthSecret:     "secret_secret_priority",
			},
			wantErr: false,
		},
		{
			name: "Check_Envs_Set",
			args: args{
				envs: map[string]string{
					"GRPC_ADDRESS": ":18090",
					"LOG_LEVEL":    "error",
					"DATABASE_URI": "postgres://postgres:postgres@localhost:15429/db?sslmode=disable",
					"JWT_SECRET":   "secret_secret",
					"CONFIG":       "",
				},
				args: []string{},
			},
			want: &Config{
				ServerGRPCHost: ":18090",
				LogLevel:       "error",
				DatabaseDSN:    "postgres://postgres:postgres@localhost:15429/db?sslmode=disable",
				AuthSecret:     "secret_secret",
			},
			wantErr: false,
		},
		{
			name: "Check_Flags_Set",
			args: args{
				envs: map[string]string{},
				args: []string{
					"-a=:28090",
					"-l=fatal",
					"-d=postgres://postgres:postgres@127.0.0.1:15429/db?sslmode=disable",
					"-s=secret_secret_flags",
				},
			},
			want: &Config{
				ServerGRPCHost: ":28090",
				LogLevel:       "fatal",
				DatabaseDSN:    "postgres://postgres:postgres@127.0.0.1:15429/db?sslmode=disable",
				AuthSecret:     "secret_secret_flags",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envs != nil {
				for k, v := range tt.args.envs {
					if err := os.Setenv(k, v); err != nil {
						t.Fatal(err)
					}
				}
			}

			got, err := New(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

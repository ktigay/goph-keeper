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
					"GRPC_ADDRESS":        ":18090",
					"LOG_LEVEL":           "error",
					"SRV_SYNC_INTERVAL":   "5000",
					"SRV_REQUEST_TIMEOUT": "500",
				},
				args: []string{
					"-a=:18080",
					"-l=error",
					"-i=4000",
				},
			},
			want: &Config{
				ServerGRPCHost:    ":18080",
				LogFile:           defaultLogFile,
				LogLevel:          "error",
				SrvSyncToInterval: 4000,
				SrvRequestTimeout: 500,
			},
			wantErr: false,
		},
		{
			name: "Check_Envs_Set",
			args: args{
				envs: map[string]string{
					"GRPC_ADDRESS":        ":18090",
					"LOG_LEVEL":           "error",
					"SRV_SYNC_INTERVAL":   "5000",
					"SRV_REQUEST_TIMEOUT": "500",
				},
				args: []string{},
			},
			want: &Config{
				ServerGRPCHost:    ":18090",
				LogFile:           defaultLogFile,
				LogLevel:          "error",
				SrvSyncToInterval: 5000,
				SrvRequestTimeout: 500,
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
					"-i=4000",
					"-t=550",
				},
			},
			want: &Config{
				ServerGRPCHost:    ":28090",
				LogFile:           defaultLogFile,
				LogLevel:          "fatal",
				SrvSyncToInterval: 4000,
				SrvRequestTimeout: 550,
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

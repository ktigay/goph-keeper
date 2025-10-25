package validator

import (
	"testing"

	"github.com/ktigay/goph-keeper/internal/client/entity"
)

func TestValidateCredentials(t *testing.T) {
	type args struct {
		data entity.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Validation_with_error_#1",
			args: args{
				data: entity.Credentials{
					Login: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Validation_with_error_#2",
			args: args{
				data: entity.Credentials{
					Login:    "",
					Password: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Validation_success",
			args: args{
				data: entity.Credentials{
					Login:    "test",
					Password: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCredentials(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

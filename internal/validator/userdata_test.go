package validator

import (
	"testing"

	"github.com/ktigay/goph-keeper/internal/entity"
)

func TestValidateUserData(t *testing.T) {
	type args struct {
		data entity.UserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty_UUID_NewRecord",
			args: args{
				data: entity.UserData{
					Title: "title",
					Type:  entity.DataTypeText,
					Data:  []byte("data"),
					IsNew: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Wrong_UUID_error",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "f72fab39-4640-49b9-9672",
					Type:  entity.DataTypeText,
					Data:  []byte("data"),
				},
			},
			wantErr: true,
		},
		{
			name: "Valid_UUID",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeText,
					Data:  []byte("data"),
				},
			},
			wantErr: false,
		},
		{
			name: "Validate_Card_with_error_#1",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeCard,
					Data:  []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":null}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Validate_Card_with_error_#2",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeCard,
					Data:  []byte(`{"number":"111111","exp_month":"11","exp_year":null,"cvc":"112"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Validate_Card_with_error_#3",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeCard,
					Data:  []byte(`{"number":"111111","exp_month":null,"exp_year":"11","cvc":"112"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Validate_Card_with_error_#4",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeCard,
					Data:  []byte(`{"number":"","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Validate_Card_success",
			args: args{
				data: entity.UserData{
					Title: "title",
					UUID:  "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:  entity.DataTypeCard,
					Data:  []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateUserData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

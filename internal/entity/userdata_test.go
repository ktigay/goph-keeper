package entity

import (
	"reflect"
	"testing"
)

func TestUserData_GetData(t *testing.T) {
	type fields struct {
		Type UserDataType
		Data []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "Get_Text_Data",
			fields: fields{
				Type: DataTypeText,
				Data: []byte("test"),
			},
			want: "test",
		},
		{
			name: "Get_Base64_Data",
			fields: fields{
				Type: DataTypeBinary,
				Data: []byte("test"),
			},
			want: "dGVzdA==",
		},
		{
			name: "Get_Card_Data",
			fields: fields{
				Type: DataTypeCard,
				Data: []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
			},
			want: UserDataCard{
				Number:   "111111",
				ExpMonth: "11",
				ExpYear:  "11",
				CVC:      "112",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserData{
				Type: tt.fields.Type,
				Data: tt.fields.Data,
			}
			if got := u.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserData_SetData(t *testing.T) {
	type fields struct {
		Type UserDataType
		Data []byte
	}
	type args struct {
		d any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Set_Base64_Data_Success",
			args: args{
				d: "dGVzdA==",
			},
			fields: fields{
				Type: DataTypeBinary,
			},
			wantErr: false,
		},
		{
			name: "Set_Base64_Data_Failed",
			args: args{
				d: 122222,
			},
			fields: fields{
				Type: DataTypeBinary,
			},
			wantErr: true,
		},
		{
			name: "Set_Base64_Data_Decode_Failed",
			args: args{
				d: "2",
			},
			fields: fields{
				Type: DataTypeBinary,
			},
			wantErr: true,
		},
		{
			name: "Set_Card_Data_Success",
			args: args{
				d: UserDataCard{
					Number:   "111111",
					ExpMonth: "11",
					ExpYear:  "11",
					CVC:      "112",
				},
			},
			fields: fields{
				Type: DataTypeCard,
			},
			wantErr: false,
		},
		{
			name: "Set_Card_Data_Failed",
			args: args{
				d: "test",
			},
			fields: fields{
				Type: DataTypeCard,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserData{
				Type: tt.fields.Type,
				Data: tt.fields.Data,
			}
			if err := u.SetData(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("SetData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

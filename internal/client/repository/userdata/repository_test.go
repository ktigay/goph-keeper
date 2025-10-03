package userdata

import (
	"reflect"
	"testing"
	"time"

	"github.com/ktigay/goph-keeper/internal/entity"
)

func Test_sortByUpdated(t *testing.T) {
	type args struct {
		data []entity.UserData
	}
	tests := []struct {
		name    string
		args    args
		want    []entity.UserData
		wantErr bool
	}{
		{
			name: "Sort_by_Updated",
			args: args{
				data: []entity.UserData{
					{
						Title:     "title1",
						UpdatedAt: time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
					},
					{
						Title:     "title2",
						UpdatedAt: time.Date(2025, 1, 11, 12, 22, 0, 0, time.UTC),
					},
					{
						Title:     "title3",
						UpdatedAt: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
					},
					{
						Title:     "title4",
						UpdatedAt: time.Date(2025, 1, 11, 12, 23, 0, 0, time.UTC),
					},
				},
			},
			want: []entity.UserData{
				{
					Title:     "title3",
					UpdatedAt: time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
				},
				{
					Title:     "title1",
					UpdatedAt: time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
				},
				{
					Title:     "title4",
					UpdatedAt: time.Date(2025, 1, 11, 12, 23, 0, 0, time.UTC),
				},
				{
					Title:     "title2",
					UpdatedAt: time.Date(2025, 1, 11, 12, 22, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortByUpdated(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortByUpdated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByUpdated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

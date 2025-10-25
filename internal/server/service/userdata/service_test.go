package userdata

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktigay/goph-keeper/internal/server/service/userdata/mocks"

	"github.com/ktigay/goph-keeper/internal/entity"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx      context.Context
		userUUID string
		data     entity.UserData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserData
		wantErr bool
	}{
		{
			name: "Create_UserData_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any()).Times(1).
						DoAndReturn(func(_ context.Context, data entity.UserData) (*entity.UserData, error) {
							return &data, nil
						})
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				data: entity.UserData{
					Title:    "title",
					UUID:     "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					UserUUID: "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:     entity.DataTypeCard,
					Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
			},
			want: &entity.UserData{
				Title:    "title",
				UUID:     "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
				UserUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				Type:     entity.DataTypeCard,
				Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
			},
			wantErr: false,
		},
		{
			name: "Create_UserData_BadRequest",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				data: entity.UserData{
					Title:    "title",
					UserUUID: "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:     entity.DataTypeCard,
					Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
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
			got, err := s.Create(tt.args.ctx, tt.args.userUUID, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx      context.Context
		userUUID string
		uuids    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete_UserData_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(), gomock.All()).Times(1).Return(nil)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				uuids: []string{
					"9b89b845-164b-498d-bc0e-f197fec9008a",
					"5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
				},
			},
			wantErr: false,
		},
		{
			name: "Delete_UserData_BadRequest",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				repo: tt.fields.repo(ctrl),
			}
			if err := s.Delete(tt.args.ctx, tt.args.userUUID, tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Read(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx      context.Context
		userUUID string
		uuids    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.UserData
		wantErr bool
	}{
		{
			name: "Read_UserData_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any(), gomock.All()).
						Times(1).
						Return([]entity.UserData{
							{
								Title: "title",
							},
						}, nil)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				uuids: []string{
					"9b89b845-164b-498d-bc0e-f197fec9008a",
					"5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
				},
			},
			want: []entity.UserData{
				{
					Title: "title",
				},
			},
			wantErr: false,
		},
		{
			name: "Read_UserData_Success_All",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any(), gomock.All()).
						Times(1).
						Return([]entity.UserData{
							{
								Title: "title",
							},
						}, nil)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
			},
			want: []entity.UserData{
				{
					Title: "title",
				},
			},
			wantErr: false,
		},
		{
			name: "Read_UserData_ErrDataNotFound",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any(), gomock.All()).Times(1).Return(nil, nil)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				uuids: []string{
					"9b89b845-164b-498d-bc0e-f197fec9008a",
					"5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Read_UserData_ReadAll_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any()).Times(1).Return([]entity.UserData{{
						Title: "title",
					}}, nil)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
			},
			want: []entity.UserData{
				{
					Title: "title",
				},
			},
			wantErr: false,
		},
		{
			name: "Read_UserData_Wrong_UUID_BadRequest",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Read(gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				uuids: []string{
					"1234-1233-1222",
				},
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
			got, err := s.Read(tt.args.ctx, tt.args.userUUID, tt.args.uuids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type fields struct {
		repo func(controller *gomock.Controller) Repository
	}
	type args struct {
		ctx      context.Context
		userUUID string
		data     entity.UserData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserData
		wantErr bool
	}{
		{
			name: "Update_UserData_Success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Update(gomock.Any(), gomock.Any()).Times(1).
						DoAndReturn(func(_ context.Context, data entity.UserData) (*entity.UserData, error) {
							return &data, nil
						})
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				data: entity.UserData{
					Title:    "title",
					UUID:     "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					UserUUID: "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Type:     entity.DataTypeCard,
					Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
			},
			want: &entity.UserData{
				Title:    "title",
				UUID:     "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
				UserUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				Type:     entity.DataTypeCard,
				Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
			},
			wantErr: false,
		},
		{
			name: "Update_UserData_BadRequest",
			fields: fields{
				repo: func(ctrl *gomock.Controller) Repository {
					repo := mocks.NewMockRepository(ctrl)
					repo.EXPECT().
						Update(gomock.Any(), gomock.Any()).Times(0)
					return repo
				},
			},
			args: args{
				ctx:      context.Background(),
				userUUID: "513bf07c-2148-43a5-8e18-d42d1548ae48",
				data: entity.UserData{
					Title:    "title",
					UUID:     "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					UserUUID: "4d8de9dc-b3b3-4c45-b71b-189fb41837ea",
					Data:     []byte(`{"number":"111111","exp_month":"11","exp_year":"11","cvc":"112"}`),
				},
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
			got, err := s.Update(tt.args.ctx, tt.args.userUUID, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package sync

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ktigay/goph-keeper/internal/client/service/sync/mocks"
	"github.com/ktigay/goph-keeper/internal/client/service/userdata"
	m "github.com/ktigay/goph-keeper/internal/client/service/userdata/mocks"
	"github.com/ktigay/goph-keeper/internal/entity"
	"github.com/ktigay/goph-keeper/internal/log"
)

func TestService_SyncFromRemote(t *testing.T) {
	type fields struct {
		client func(*gomock.Controller) Client
		repo   func(*gomock.Controller) userdata.Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "SyncFromRemote_Success",
			fields: fields{
				client: func(ctrl *gomock.Controller) Client {
					cl := mocks.NewMockClient(ctrl)
					cl.EXPECT().Read(gomock.Any()).Times(1).Return(
						[]entity.UserData{
							{
								UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
								Title: "Test",
								Type:  entity.DataTypeText,
								Data:  []byte("Test"),
								IsNew: true,
							},
						}, nil)
					return cl
				},
				repo: func(ctrl *gomock.Controller) userdata.Repository {
					e := entity.UserData{
						UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: true,
					}
					repo := m.NewMockRepository(ctrl)
					repo.EXPECT().Read(gomock.Any(), gomock.Any()).Times(1).Return([]entity.UserData{e}, nil)
					n := e
					n.IsSynced = true
					repo.EXPECT().Replace(gomock.Any(), gomock.Eq(n)).Times(1).Return(&entity.UserData{}, nil)
					return repo
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			if err := s.SyncFromRemote(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("SyncFromRemote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Initialize(t *testing.T) {
	type fields struct {
		client func(*gomock.Controller) Client
		repo   func(*gomock.Controller) userdata.Repository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []entity.UserData
		wantErr bool
	}{
		{
			name: "Initialize_Success",
			fields: fields{
				client: func(ctrl *gomock.Controller) Client {
					cl := mocks.NewMockClient(ctrl)
					cl.EXPECT().Read(gomock.Any()).Times(1).Return(
						[]entity.UserData{
							{
								UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
								Title: "Test",
								Type:  entity.DataTypeText,
								Data:  []byte("Test"),
							},
						}, nil)
					return cl
				},
				repo: func(ctrl *gomock.Controller) userdata.Repository {
					d := []entity.UserData{
						{
							UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
							Title: "Test",
							Type:  entity.DataTypeText,
							Data:  []byte("Test"),
						},
					}
					repo := m.NewMockRepository(ctrl)
					repo.EXPECT().Sync(gomock.Any(), gomock.Eq(d)).Times(1).Return(nil)
					return repo
				},
			},
			want: []entity.UserData{
				{
					UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
					Title: "Test",
					Type:  entity.DataTypeText,
					Data:  []byte("Test"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := &Service{
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			got, err := s.Initialize(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SyncToRemote(t *testing.T) {
	type fields struct {
		client func(*gomock.Controller) Client
		repo   func(*gomock.Controller) userdata.Repository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []entity.UserData
		wantErr bool
	}{
		{
			name: "SyncToRemote_Created_Success",
			fields: fields{
				client: func(ctrl *gomock.Controller) Client {
					e := entity.UserData{
						UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: true,
					}
					cl := mocks.NewMockClient(ctrl)
					cl.EXPECT().Create(gomock.Any(), gomock.Eq(e)).Times(1).Return(&e, nil)
					return cl
				},
				repo: func(ctrl *gomock.Controller) userdata.Repository {
					e := entity.UserData{
						UUID:  "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: true,
					}
					repo := m.NewMockRepository(ctrl)
					repo.EXPECT().ReadUnsynced(gomock.Any()).Return([]entity.UserData{e}, nil)
					n := e
					n.IsNew = false
					n.IsSynced = true
					repo.EXPECT().Replace(gomock.Any(), gomock.Eq(n)).Times(1).Return(&entity.UserData{}, nil)
					return repo
				},
			},
			want: []entity.UserData{
				{
					UUID:     "5d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
					Title:    "Test",
					Type:     entity.DataTypeText,
					Data:     []byte("Test"),
					IsNew:    false,
					IsSynced: true,
				},
			},
			wantErr: false,
		},
		{
			name: "SyncToRemote_Updated_Success",
			fields: fields{
				client: func(ctrl *gomock.Controller) Client {
					e := entity.UserData{
						UUID:  "1d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: false,
					}
					cl := mocks.NewMockClient(ctrl)
					cl.EXPECT().Update(gomock.Any(), gomock.Eq(e)).Times(1).Return(&e, nil)
					return cl
				},
				repo: func(ctrl *gomock.Controller) userdata.Repository {
					e := entity.UserData{
						UUID:  "1d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: false,
					}
					repo := m.NewMockRepository(ctrl)
					repo.EXPECT().ReadUnsynced(gomock.Any()).Return([]entity.UserData{e}, nil)
					n := e
					n.IsNew = false
					n.IsSynced = true
					repo.EXPECT().Replace(gomock.Any(), gomock.Eq(n)).Times(1).Return(&entity.UserData{}, nil)
					return repo
				},
			},
			want: []entity.UserData{
				{
					UUID:     "1d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
					Title:    "Test",
					Type:     entity.DataTypeText,
					Data:     []byte("Test"),
					IsNew:    false,
					IsSynced: true,
				},
			},
			wantErr: false,
		},
		{
			name: "SyncToRemote_Update_Failed",
			fields: fields{
				client: func(ctrl *gomock.Controller) Client {
					e := entity.UserData{
						UUID:  "2d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: false,
					}
					cl := mocks.NewMockClient(ctrl)
					cl.EXPECT().Update(gomock.Any(), gomock.Eq(e)).Times(1).Return(nil, fmt.Errorf("some error"))
					return cl
				},
				repo: func(ctrl *gomock.Controller) userdata.Repository {
					repo := m.NewMockRepository(ctrl)
					repo.EXPECT().ReadUnsynced(gomock.Any()).Return([]entity.UserData{{
						UUID:  "2d33d26d-47b5-4f6e-ac39-cf63f5a1c1cf",
						Title: "Test",
						Type:  entity.DataTypeText,
						Data:  []byte("Test"),
						IsNew: false,
					}}, nil)
					repo.EXPECT().Replace(gomock.Any(), gomock.Any()).Times(0)
					return repo
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
				client: tt.fields.client(ctrl),
				repo:   tt.fields.repo(ctrl),
				logger: log.MockLogger,
			}
			got, err := s.SyncToRemote(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncToRemote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncToRemote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

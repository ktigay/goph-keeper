package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
	"github.com/ktigay/goph-keeper/internal/entity"
)

// MapItemToEntity мапит [data.UserDataItem] в [entity.UserData].
func MapItemToEntity(item *data.UserDataItem, userUUID string) entity.UserData {
	return entity.UserData{
		UUID:     item.Uuid,
		UserUUID: userUUID,
		Title:    item.Title,
		Type:     entity.UserDataType(data.UserDataItem_DataType_name[int32(item.Type)]),
		Data:     item.Data,
		MetaData: func() []entity.MetaData {
			d := make([]entity.MetaData, 0, len(item.Metadata))
			for _, m := range item.Metadata {
				d = append(d, entity.MetaData{
					Title: m.Title,
					Value: m.Value,
				})
			}
			return d
		}(),
		CreatedAt: item.CreatedAt.AsTime(),
		UpdatedAt: item.UpdatedAt.AsTime(),
	}
}

// MapEntityToItem мапит [entity.UserData] в [data.UserDataItem].
func MapEntityToItem(e entity.UserData) *data.UserDataItem {
	return &data.UserDataItem{
		Uuid:  e.UUID,
		Title: e.Title,
		Type: data.UserDataItem_DataType(
			data.UserDataItem_DataType_value[string(e.Type)],
		),
		Data: e.Data,
		Metadata: func() []*data.MetaData {
			d := make([]*data.MetaData, 0, len(e.MetaData))
			for _, m := range e.MetaData {
				d = append(d, &data.MetaData{
					Title: m.Title,
					Value: m.Value,
				})
			}
			return d
		}(),
		CreatedAt: timestamppb.New(e.CreatedAt),
		UpdatedAt: timestamppb.New(e.UpdatedAt),
	}
}

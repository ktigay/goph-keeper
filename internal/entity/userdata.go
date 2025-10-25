package entity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

// UserDataType тип данных.
type UserDataType string

var (
	// DataTypeText текстовый.
	DataTypeText UserDataType = "TEXT"
	// DataTypeBinary бинарный.
	DataTypeBinary UserDataType = "BINARY"
	// DataTypeCard кредитная карта.
	DataTypeCard UserDataType = "CARD"
	// DataTypes типы данных.
	DataTypes = []UserDataType{DataTypeText, DataTypeBinary, DataTypeCard}
)

// UserData сущность пользовательских данных.
type UserData struct {
	UUID      string       `validate:"required_if=IsNew false,omitempty,uuid"`
	UserUUID  string       `validate:"omitempty,uuid"`
	Title     string       `validate:"required"`
	Type      UserDataType `validate:"required"`
	Data      []byte       `validate:"required"`
	MetaData  []MetaData   `validate:"omitempty"`
	IsSynced  bool
	IsNew     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SetData устанавливает значение для [UserData.Data].
func (u *UserData) SetData(d any) error {
	switch u.Type {
	case DataTypeText:
		u.Data = []byte(d.(string))
	case DataTypeBinary:
		t, ok := d.(string)
		if !ok {
			return errors.New("type assertion fail")
		}
		b, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return err
		}
		u.Data = b
	case DataTypeCard:
		var (
			dv  []byte
			err error
		)
		s, ok := d.(UserDataCard)
		if !ok {
			return errors.New("invalid user data type")
		}
		if dv, err = json.Marshal(s); err != nil {
			return err
		}
		u.Data = dv
	}
	return nil
}

// GetData возвращает смапленное значение для [UserData.Data].
func (u *UserData) GetData() any {
	switch u.Type {
	case DataTypeText:
		return string(u.Data)
	case DataTypeBinary:
		return base64.StdEncoding.EncodeToString(u.Data)
	case DataTypeCard:
		c := UserDataCard{}
		err := json.Unmarshal(u.Data, &c)
		if err != nil {
			return UserDataCard{}
		}
		return c
	}
	return nil
}

// UserDataCard стркутура типа [DataTypeCard].
type UserDataCard struct {
	Number   string `json:"number" validate:"required"`
	ExpMonth string `json:"exp_month" validate:"required"`
	ExpYear  string `json:"exp_year" validate:"required"`
	CVC      string `json:"cvc" validate:"required"`
}

// MetaData метаданные.
type MetaData struct {
	Title string
	Value string
}

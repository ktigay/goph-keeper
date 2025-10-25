package validator

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/ktigay/goph-keeper/internal/entity"
)

// ValidateUserData валидирует [entity.UserData].
func ValidateUserData(data entity.UserData) error {
	vd := validator.New()
	if err := vd.Struct(&data); err != nil {
		return err
	}

	if data.Type == entity.DataTypeCard {
		card := entity.UserDataCard{}
		if err := json.Unmarshal(data.Data, &card); err != nil {
			return err
		}
		if err := vd.Struct(card); err != nil {
			return err
		}
	}
	return nil
}

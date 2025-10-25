package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/ktigay/goph-keeper/internal/client/entity"
)

// ValidateCredentials валидирует аутентификационные данные.
func ValidateCredentials(data entity.Credentials) error {
	data.Login = strings.TrimSpace(data.Login)
	data.Password = strings.TrimSpace(data.Password)
	vd := validator.New()
	if err := vd.Struct(data); err != nil {
		return err
	}
	return nil
}

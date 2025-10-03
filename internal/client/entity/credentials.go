package entity

// Credentials пользовательские данные.
type Credentials struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

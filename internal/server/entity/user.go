package entity

import "time"

// User структура пользователя.
type User struct {
	UUID      string
	Login     string `validate:"required"`
	Password  string `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

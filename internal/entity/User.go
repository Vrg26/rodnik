package entity

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Phone    string    `json:"phone" db:"phone"`
	Password string    `json:"password" db:"password"`
	Leaves   float64   `json:"leaves" db:"balance"`
}

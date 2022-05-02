package entity

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" format:"uuid"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Leaves   float64   `json:"leaves"`
}

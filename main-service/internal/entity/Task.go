package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id            uuid.UUID `json:"id" db:"id"`
	Title         string    `json:"title,omitempty" db:"title"`
	Description   string    `json:"description,omitempty" db:"description"`
	Status        string    `json:"status,omitempty" db:"status_id"`
	CreatorId     uuid.UUID `json:"creator_id,omitempty" db:"creator_id"`
	HelperId      uuid.UUID `json:"helper_id,omitempty" db:"helper_id"`
	Cost          float64   `json:"cost,omitempty" db:"cost"`
	DateRelevance time.Time `json:"date_relevance,omitempty" db:"date_relevance"`
	CreatedOn     time.Time `json:"created_on" db:"created_on"`
	Overdue       bool      `json:"overdue,omitempty" db:"overdue"`
}

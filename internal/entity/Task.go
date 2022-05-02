package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title,omitempty"`
	Description   string    `json:"description,omitempty"`
	Status        string    `json:"status,omitempty"`
	CreatorId     uuid.UUID `json:"creator_id,omitempty"`
	HelperId      uuid.UUID `json:"helper_id,omitempty"`
	Cost          float64   `json:"cost,omitempty"`
	DateRelevance time.Time `json:"date_relevance,omitempty"`
	CreatedOn     time.Time `json:"created_on"`
	Overdue       bool      `json:"overdue,omitempty"`
}

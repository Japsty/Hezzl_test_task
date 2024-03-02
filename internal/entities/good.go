package entities

import (
	"time"
)

// Good - основная структура, содержит все те же поля, что и в таблице Goods
// @Good
type Good struct {
	ID          int       `json:"id,omitempty"`
	ProjectID   int       `json:"projectId" validate:"required,numeric,min=0"`
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description string    `json:"description,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Removed     bool      `json:"removed,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

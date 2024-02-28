package entities

import (
	"time"
)

type Good struct {
	ID          int       `json:"id,omitempty"`
	ProjectID   int       `json:"projectId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Removed     bool      `json:"removed,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type AddGoodRequest struct {
	Name string `json:"name"`
}

type UpdateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RemoveGoodResponse struct {
	ID         int  `json:"id"`
	CampaignID int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}

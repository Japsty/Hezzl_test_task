package storage

import "Hezzl_test_task/internal/entities"

type RemoveGoodResponse struct {
	ID         int  `json:"id"`
	CampaignID int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type ListGoodsResponse struct {
	Meta  Meta            `json:"meta"`
	Goods []entities.Good `json:"goods"`
}

type PriorityObj struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

type ReprioritiizeResponse struct {
	Priorities []PriorityObj `json:"priorities"`
}

type ClickhouseLog struct {
	ID          int    `db:"id"`
	ProjectID   int    `db:"project_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Priority    int    `db:"priority"`
	Removed     bool   `db:"removed"`
	CreatedAt   string `db:"created_at"`
}

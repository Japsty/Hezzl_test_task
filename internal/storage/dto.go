package storage

import (
	"Hezzl_test_task/internal/entities"
	"time"
)

// RemoveGoodResponse - форматированный вывод метода RemoveGood
type RemoveGoodResponse struct {
	ID         int  `json:"id"`
	CampaignID int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}

// Meta часть ответа метода ListGoods
type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

// ListGoodsResponse - форматированный вывод метода ListGoods, содержит в себе структуру Meta
type ListGoodsResponse struct {
	Meta  Meta            `json:"meta"`
	Goods []entities.Good `json:"goods"`
}

// PriorityObj составляющая структуры ReprioritiizeResponse
type PriorityObj struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

// ReprioritiizeResponse - вывод метода ReprioritiizeGood
type ReprioritiizeResponse struct {
	Priorities []PriorityObj `json:"priorities"`
}

type ClickhouseLog struct {
	ID          int       `db:"id"`
	ProjectID   int       `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     uint8     `json:"removed"`
	EventTime   time.Time `db:"created_at"`
}

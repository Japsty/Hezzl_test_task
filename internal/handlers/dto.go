package handlers

type addGoodRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type updateGoodRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

type idsRequest struct {
	Id        int `form:"id" validate:"required,numeric,gt=0"`
	ProjectId int `form:"projectId" validate:"required,numeric,gt=0"`
}

type goodsRequest struct {
	Limit  int `form:"limit" validate:"omitempty,gt=0" default:"10"`
	Offset int `form:"offset" validate:"omitempty,gte=0" default:"1"`
}

type patchGoodReprioritiizeRequest struct {
	NewPriority int `json:"newPriority" validate:"required,numeric,min=0"`
}

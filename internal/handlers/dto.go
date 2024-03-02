package handlers

// addGoodRequest - форматированный запрос к методу AddGood
type addGoodRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"omitempty,min=3,max=255"`
}

// updateGoodRequest - форматированный запрос к методу PatchGoodUpdate
type updateGoodRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

// idsRequest - форматированный запрос с id для удобной валидации
type idsRequest struct {
	Id        int `form:"id" validate:"required,numeric,gt=0"`
	ProjectId int `form:"projectId" validate:"required,numeric,gt=0"`
}

// goodsRequest - структура запроса для метода ListGoods
// сделана для удобства валидирования
type goodsRequest struct {
	Limit  int `form:"limit" validate:"omitempty,gt=0" default:"10"`
	Offset int `form:"offset" validate:"omitempty,gte=0" default:"1"`
}

// patchGoodReprioritiizeRequest - структура для валидирования запроса к patchGoodReprioritiize
type patchGoodReprioritiizeRequest struct {
	NewPriority int `json:"newPriority" validate:"required,numeric,min=0"`
}

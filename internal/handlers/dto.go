package handlers

type AddGoodRequest struct {
	Name string `json:"name"`
}

type UpdateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PatchGoodReprioritiizeRequest struct {
	NewPriority int `json:"newPriority"`
}

package entities

type priority struct {
	ID       int
	Priority int
}

type priorities []priority

type ReprioritiizeResponse struct {
	priorities
}

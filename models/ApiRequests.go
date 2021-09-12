package models

type LoginRequest struct {
	Username string `json:"username"`
}

type RefreshRequest struct {
	Username string `json:"username"`
	PlanetId string `json:"planet_id"`
}

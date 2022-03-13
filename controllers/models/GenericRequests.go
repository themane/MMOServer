package models

type LoginRequest struct {
	Username string `json:"username" example:"devashish"`
}

type DeployRequest struct {
	Username string `json:"username" example:"devashish"`
	PlanetId string `json:"planet_id" example:"001:002:03"`
}

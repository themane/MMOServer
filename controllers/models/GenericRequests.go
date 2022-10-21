package models

type LoginRequest struct {
	Username string `json:"username" example:"devashish"`
	IdToken  string `json:"id_token" example:"asdbf1412b"`
}

type DeployRequest struct {
	Username string `json:"username" example:"devashish"`
	PlanetId string `json:"planet_id" example:"001:002:03"`
}

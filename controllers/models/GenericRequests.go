package models

type RegistrationRequest struct {
	Username string `json:"username" example:"asjfdkj13"`
	Species  string `json:"species" example:"ARYAN"`
	Birthday string `json:"birthday" example:"2001-01-12"`
	Location string `json:"location" example:"UTC"`
}

type DeployRequest struct {
	PlanetId string `json:"planet_id" example:"001:002:03"`
}

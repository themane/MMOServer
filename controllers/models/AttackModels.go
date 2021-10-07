package models

import "github.com/themane/MMOServer/models"

type AttackResponse struct {
	AttackTime string `json:"attack_time"`
	ReturnTime string `json:"return_time"`
}

type SpyRequest struct {
	Attacker     string             `json:"attacker" example:"devashish"`
	FromPlanetId string             `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string             `json:"to_planet_id" example:"001:002:05"`
	Scouts       []models.Formation `json:"scouts"`
}

type AttackRequest struct {
	Attacker     string                                   `json:"attacker" example:"devashish"`
	FromPlanetId string                                   `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string                                   `json:"to_planet_id" example:"001:002:05"`
	Formation    map[string]map[string][]models.Formation `json:"formation"`
}

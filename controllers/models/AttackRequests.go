package models

type SpyRequest struct {
	Attacker     string      `json:"attacker" example:"devashish"`
	FromPlanetId string      `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string      `json:"to_planet_id" example:"001:002:05"`
	Scouts       []Formation `json:"scouts"`
}

type AttackRequest struct {
	Attacker     string                            `json:"attacker" example:"devashish"`
	FromPlanetId string                            `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string                            `json:"to_planet_id" example:"001:002:05"`
	Formation    map[string]map[string][]Formation `json:"formation"`
}

type Formation struct {
	Name     string `json:"name" example:"ANUJ"`
	Quantity int    `json:"quantity" example:"15"`
}

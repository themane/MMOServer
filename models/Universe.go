package models

type Universe struct {
	NumSystems int                `json:"num_systems"`
	Systems    map[string]SystemU `json:"systems"`
}

type SystemU struct {
	System     int                `json:"system"`
	NumSystems int                `json:"num_systems"`
	Sectors    map[string]SectorU `json:"sectors"`
}

type SectorU struct {
	Sector     int                `json:"sector"`
	NumPlanets int                `json:"num_planets"`
	Planets    map[string]PlanetU `json:"planets"`
}

type PlanetU struct {
	Planet       int     `json:"planet"`
	Mines        []MineU `json:"mines"`
	PlanetConfig string  `json:"planet_config"`
}

type MineU struct {
	Id           string `json:"_id"`
	Type         string `json:"type"`
	MaxLimit     int    `json:"max_limit"`
	IncreaseRate int    `json:"increase_rate"`
}

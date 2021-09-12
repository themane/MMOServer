package models

type Universe struct {
	NumSystems int                  `json:"num_systems"`
	Systems    map[string]SystemUni `json:"systems"`
}

type SystemUni struct {
	System     int                  `json:"system"`
	NumSystems int                  `json:"num_systems"`
	Sectors    map[string]SectorUni `json:"sectors"`
}

type SectorUni struct {
	Sector     int                  `json:"sector"`
	NumPlanets int                  `json:"num_planets"`
	Planets    map[string]PlanetUni `json:"planets"`
}

type PlanetUni struct {
	Planet       int                `json:"planet"`
	Mines        map[string]MineUni `json:"mines"`
	PlanetConfig string             `json:"planet_config"`
}

type MineUni struct {
	Id           string       `json:"_id"`
	Type         ResourceType `json:"type"`
	MaxLimit     int          `json:"max_limit"`
	IncreaseRate int          `json:"increase_rate"`
}

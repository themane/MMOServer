package models

type PlanetUni struct {
	Id           string             `json:"_id"`
	Position     PlanetPosition     `json:"position"`
	Mines        map[string]MineUni `json:"mines"`
	PlanetConfig string             `json:"planet_config"`
	Occupied     bool               `json:"occupied"`
	Distance     int                `json:"distance"`
}

type MineUni struct {
	Id           string `json:"_id"`
	Type         string `json:"type"`
	MaxLimit     int    `json:"max_limit"`
	IncreaseRate int    `json:"increase_rate"`
}

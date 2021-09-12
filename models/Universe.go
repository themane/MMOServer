package models

type Universe struct {
	Planets map[string]struct {
		Mines []struct {
			Id           string `json:"_id"`
			Type         string `json:"type"`
			MaxLimit     int    `json:"max_limit"`
			IncreaseRate int    `json:"increase_rate"`
		} `json:"mines"`
		Position PlanetPosition `json:"position"`
	} `json:"planets"`
}

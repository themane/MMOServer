package models

type LoginResponse struct {
	Profile         Profile            `json:"profile"`
	HomeSector      Sector             `json:"home_sector"`
	OccupiedPlanets []StaticPlanetData `json:"occupied_planets"`
}

type StaticPlanetData struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
	Home         bool           `json:"home" example:"true"`
}

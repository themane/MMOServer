package models

type LoginResponse struct {
	Profile         Profile            `json:"profile"`
	HomeSector      Sector             `json:"home_sector"`
	HomePlanet      OccupiedPlanet     `json:"home_planet"`
	OccupiedPlanets []StaticPlanetData `json:"occupied_planets"`
}

type StaticPlanetData struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
	Home         bool           `json:"home" example:"true"`
}

func (p *StaticPlanetData) Init(planetUni PlanetUni, homeSectorId string) {
	p.PlanetConfig = planetUni.PlanetConfig
	p.Home = homeSectorId == planetUni.Position.SectorId()
	p.Position = planetUni.Position.Clone()
}

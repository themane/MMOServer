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

func (planet *StaticPlanetData) Init(universalStaticPlanetData PlanetU, position PlanetPosition, home bool) {
	planet.PlanetConfig = universalStaticPlanetData.PlanetConfig
	planet.Position = position.Clone()
	planet.Home = home
}

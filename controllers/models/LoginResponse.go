package models

import (
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type LoginResponse struct {
	Profile         Profile            `json:"profile"`
	HomeSector      Sector             `json:"home_sector"`
	HomePlanetId    string             `json:"home_planet_id"`
	OccupiedPlanets []StaticPlanetData `json:"occupied_planets"`
}

type StaticPlanetData struct {
	PlanetConfig string                `json:"planet_config" example:"Planet2.json"`
	Position     models.PlanetPosition `json:"position"`
	Home         bool                  `json:"home" example:"true"`
}

func (p *StaticPlanetData) Init(planetUni repoModels.PlanetUni, homeSectorId string) {
	p.PlanetConfig = planetUni.PlanetConfig
	p.Home = homeSectorId == planetUni.Position.SectorId()
	p.Position = planetUni.Position.Clone()
}

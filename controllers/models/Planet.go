package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type UnoccupiedPlanet struct {
	PlanetConfig string                `json:"planet_config" example:"Planet2.json"`
	Position     models.PlanetPosition `json:"position"`
}

func (u *UnoccupiedPlanet) Init(planet repoModels.PlanetUni) {
	u.Position = planet.Position.Clone()
	u.PlanetConfig = planet.PlanetConfig
}

type OccupiedPlanet struct {
	PlanetConfig string                `json:"planet_config" example:"Planet2.json"`
	Position     models.PlanetPosition `json:"position"`
	Resources    Resources             `json:"resources"`
	Population   Population            `json:"population"`
	Mines        []Mine                `json:"mines"`
	Home         bool                  `json:"home" example:"true"`
}

func (o *OccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser,
	waterConstants constants.ResourceConstants, grapheneConstants constants.ResourceConstants) {

	o.Position = planetUni.Position.Clone()
	o.PlanetConfig = planetUni.PlanetConfig
	o.Resources.Init(planetUser)
	o.Population.Init(planetUser)
	for mineId, mineUser := range planetUser.Mines {
		mine := Mine{}
		mine.Init(planetUni.Mines[mineId], mineUser, waterConstants, grapheneConstants)
		o.Mines = append(o.Mines, mine)
	}
	o.Home = planetUser.Home
}

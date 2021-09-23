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
	Shields      []Shield              `json:"shields"`
	Home         bool                  `json:"home" example:"true"`
}

func (o *OccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser,
	buildingConstants map[string]constants.BuildingConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants) {

	o.Position = planetUni.Position.Clone()
	o.PlanetConfig = planetUni.PlanetConfig
	o.Resources.Init(planetUser)
	o.Population.Init(planetUser)
	o.Shields = InitAllShields(planetUser, defenceConstants[constants.Shield], buildingConstants[constants.Shield])
	for mineId := range planetUser.Mines {
		mine := Mine{}
		mine.Init(planetUni.Mines[mineId], planetUser,
			buildingConstants[constants.WaterMiningPlant], buildingConstants[constants.GrapheneMiningPlant],
			waterConstants, grapheneConstants,
		)
		o.Mines = append(o.Mines, mine)
	}
	o.Home = planetUser.Home
}

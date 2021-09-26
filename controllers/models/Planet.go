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
	PlanetConfig            string                `json:"planet_config" example:"Planet2.json"`
	Position                models.PlanetPosition `json:"position"`
	Resources               Resources             `json:"resources"`
	Population              Population            `json:"population"`
	Mines                   []Mine                `json:"mines"`
	Shields                 []Shield              `json:"shields"`
	IdleDefences            []Defence             `json:"idle_defences" bson:"idle_defences"`
	IdleDefenceShipCarriers []DefenceShipCarrier  `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	AvailableShips          []Ship                `json:"available_ships" bson:"available_ships"`
	Home                    bool                  `json:"home" example:"true"`
}

func (o *OccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser,
	buildingConstants map[string]constants.BuildingConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants, shipConstants map[string]constants.ShipConstants) {

	o.PlanetConfig = planetUni.PlanetConfig
	o.Position = planetUni.Position.Clone()
	o.Resources.Init(planetUser)
	o.Population.Init(planetUser)
	for mineId := range planetUser.Mines {
		mine := Mine{}
		mine.Init(planetUni.Mines[mineId], planetUser,
			buildingConstants[constants.WaterMiningPlant], buildingConstants[constants.GrapheneMiningPlant],
			waterConstants, grapheneConstants,
		)
		o.Mines = append(o.Mines, mine)
	}
	o.Shields = InitAllShields(planetUser, defenceConstants, buildingConstants[constants.Shield])
	o.IdleDefences = InitAllIdleDefences(planetUser.Defences, defenceConstants)
	o.IdleDefenceShipCarriers = InitAllIdleDefenceShipCarriers(planetUser, defenceConstants[constants.Vikram], shipConstants)
	for shipName, shipUser := range planetUser.Ships {
		deployedShips := 0
		for _, vikram := range planetUser.DefenceShipCarriers {
			deployedShips += vikram.HostingShips[shipName]
		}
		availableShips := shipUser.Quantity - deployedShips
		if availableShips <= 0 {
			continue
		}
		s := Ship{}
		s.Init(shipName, availableShips, shipUser, shipConstants[shipName])
		o.AvailableShips = append(o.AvailableShips, s)
	}
	o.Home = planetUser.Home
}

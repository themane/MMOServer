package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type UserResponse struct {
	Profile         Profile               `json:"profile"`
	HomeSector      Sector                `json:"home_sector"`
	OccupiedPlanets []StaticPlanetData    `json:"occupied_planets"`
	Notifications   []models.Notification `json:"notifications"`
}

type PlanetResponse struct {
	OccupiedPlanet OccupiedPlanet        `json:"occupied_planet"`
	Notifications  []models.Notification `json:"notifications"`
}

type UserPlanetResponse struct {
	BasePlanet              bool                  `json:"base_planet" example:"true"`
	Resources               Resources             `json:"resources"`
	Population              Population            `json:"population"`
	Shields                 []Shield              `json:"shields"`
	IdleDefences            []Defence             `json:"idle_defences" bson:"idle_defences"`
	IdleDefenceShipCarriers []DefenceShipCarrier  `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	HomePlanet              bool                  `json:"home_planet" example:"true"`
	Notifications           []models.Notification `json:"notifications"`
}

func (p *UserPlanetResponse) Init(planetUser repoModels.PlanetUser,
	buildingConstants map[string]constants.BuildingConstants,
	defenceConstants map[string]constants.DefenceConstants, shipConstants map[string]constants.ShipConstants,
	notifications []models.Notification) {

	if planetUser.BasePlanet {
		p.BasePlanet = true
		return
	}
	p.Resources.Init(planetUser)
	p.Population.Init(planetUser)
	p.Shields = InitAllShields(planetUser, defenceConstants, buildingConstants[constants.Shield])
	p.IdleDefences = InitAllIdleDefences(planetUser.Defences, defenceConstants)
	p.IdleDefenceShipCarriers = InitAllIdleDefenceShipCarriers(planetUser, defenceConstants[constants.Vikram], shipConstants)
	p.HomePlanet = planetUser.HomePlanet
	p.Notifications = notifications
}

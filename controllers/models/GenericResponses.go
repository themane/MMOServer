package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/buildings"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type UserResponse struct {
	Profile         Profile               `json:"profile"`
	HomeSector      Sector                `json:"home_sector"`
	OccupiedPlanets []StaticPlanetData    `json:"occupied_planets"`
	Notifications   []models.Notification `json:"notifications"`
}

type SectorResponse struct {
	Sector          Sector                `json:"sector"`
	OccupiedPlanets []StaticPlanetData    `json:"occupied_planets"`
	Notifications   []models.Notification `json:"notifications"`
}

type PlanetResponse struct {
	OccupiedPlanet OccupiedPlanet        `json:"occupied_planet"`
	Notifications  []models.Notification `json:"notifications"`
}

type UserPlanetResponse struct {
	BasePlanet              bool                          `json:"base_planet" example:"true"`
	Resources               *Resources                    `json:"resources"`
	Population              *Population                   `json:"population"`
	Shields                 []buildings.Shield            `json:"shields"`
	IdleDefences            []military.Defence            `json:"idle_defences" bson:"idle_defences"`
	IdleDefenceShipCarriers []military.DefenceShipCarrier `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	HomePlanet              bool                          `json:"home_planet" example:"true"`
	Notifications           []models.Notification         `json:"notifications"`
}

func (p *UserPlanetResponse) Init(planetUser repoModels.PlanetUser,
	upgradeConstants map[string]constants.UpgradeConstants,
	defenceConstants map[string]constants.DefenceConstants, shipConstants map[string]constants.ShipConstants,
	notifications []models.Notification) {

	if planetUser.BasePlanet {
		p.BasePlanet = true
		return
	}
	p.Resources = InitResources(planetUser)
	p.Population = InitPopulation(planetUser)
	p.Shields = buildings.InitAllShields(planetUser, defenceConstants, upgradeConstants[constants.Shield])
	p.IdleDefences = military.InitAllIdleDefences(planetUser.Defences, defenceConstants)
	p.IdleDefenceShipCarriers = military.InitAllIdleDefenceShipCarriers(planetUser, defenceConstants[constants.Vikram], shipConstants)
	p.HomePlanet = planetUser.HomePlanet
	p.Notifications = notifications
}

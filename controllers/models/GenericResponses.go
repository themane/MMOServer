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
	BasePlanet          bool                          `json:"base_planet" example:"true"`
	Resources           *Resources                    `json:"resources"`
	Population          *Population                   `json:"population"`
	Shields             []buildings.Shield            `json:"shields"`
	Defences            []military.Defence            `json:"idle_defences" bson:"idle_defences"`
	DefenceShipCarriers []military.DefenceShipCarrier `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	HomePlanet          bool                          `json:"home_planet" example:"true"`
	Notifications       []models.Notification         `json:"notifications"`
}

func (p *UserPlanetResponse) Init(planetUser repoModels.PlanetUser,
	upgradeConstants map[string]constants.UpgradeConstants, sheildConstants map[string]map[string]interface{},
	militaryConstants map[string]constants.MilitaryConstants, speciesConstants constants.SpeciesConstants,
	notifications []models.Notification) {

	if planetUser.BasePlanet {
		p.BasePlanet = true
		return
	}
	p.Resources = InitResources(planetUser)
	p.Population = InitPopulation(planetUser, militaryConstants)
	p.Shields = buildings.InitAllShields(planetUser, sheildConstants, upgradeConstants[constants.Shield])
	for _, unitName := range speciesConstants.AvailableUnits {
		if militaryConstant, ok := militaryConstants[unitName]; ok {
			if militaryConstant.Type == constants.Defender {
				d := military.Defence{}
				d.Init(unitName, planetUser.Defences[unitName], militaryConstant.Levels)
				p.Defences = append(p.Defences, d)
			}
		}
	}
	p.DefenceShipCarriers = military.InitAllDefenceShipCarriers(planetUser, militaryConstants)
	p.HomePlanet = planetUser.HomePlanet
	p.Notifications = notifications
}

type ErrorResponse struct {
	Message  string `json:"message"`
	HttpCode int    `json:"http_code"`
}

package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/buildings"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

//  Planet Types
const (
	ResourcePlanet  string = "RESOURCE_PLANET"
	AbandonedPlanet string = "ABANDONED_PLANET"
	HomePlanet      string = "HOME_PLANET"
	BasePlanet      string = "BASE_PLANET"
	EnemyPlanet     string = "ENEMY_PLANET"
)

type UnoccupiedPlanet struct {
	Position     models.PlanetPosition     `json:"position"`
	PlanetConfig string                    `json:"planet_config" example:"Planet2.json"`
	Distance     int                       `json:"distance" example:"14"`
	Shields      []UnoccupiedPlanetShield  `json:"shields"`
	Water        int                       `json:"water" example:"150"`
	Graphene     int                       `json:"graphene" example:"140"`
	Defences     []UnoccupiedPlanetDefence `json:"defences"`
	Occupied     string                    `json:"occupied" example:"devashish"`
	Invulnerable bool                      `json:"invulnerable" example:"true"`
	PlanetType   string                    `json:"planet_type" example:"ENEMY_PLANET"`
}

type UnoccupiedPlanetShield struct {
	Id   string `json:"_id" example:"SHLD101"`
	Type string `json:"type" example:"INVULNERABLE"`
}

type UnoccupiedPlanetDefence struct {
	Type     string `json:"type" example:"BOMBER"`
	Level    int    `json:"level" example:"1"`
	Quantity int    `json:"quantity" example:"5"`
}

func (u *UnoccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser, occupiedUser string) {
	u.Position = planetUni.Position.Clone()
	u.PlanetConfig = planetUni.PlanetConfig
	u.Distance = planetUni.Distance
	u.Occupied = occupiedUser
	u.Invulnerable = false
	shieldIds := constants.GetShieldIds()
	planetType := planetUni.GetPlanetType()
	if planetType == constants.Primitive || planetType == constants.Resource {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Unavailable})
		}
		u.PlanetType = ResourcePlanet
	} else if planetType == constants.Base {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Invulnerable})
		}
		u.PlanetType = BasePlanet
		u.Invulnerable = true
	} else if planetType == constants.Bot {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Broken})
		}
		u.PlanetType = AbandonedPlanet
		u.Water = planetUser.Water.Amount
		u.Graphene = planetUser.Graphene.Amount
		for defenceType, defence := range planetUser.Defences {
			u.Defences = append(u.Defences, UnoccupiedPlanetDefence{defenceType, defence.Level, defence.Quantity})
		}
	} else {
		for _, shieldId := range shieldIds {
			if planetUser.Buildings[shieldId].BuildingLevel > 0 {
				u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Active})
			} else {
				u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Disabled})
			}
		}
		u.PlanetType = EnemyPlanet
		u.Water = planetUser.Water.Amount
		u.Graphene = planetUser.Graphene.Amount
		for defenceType, defence := range planetUser.Defences {
			u.Defences = append(u.Defences, UnoccupiedPlanetDefence{defenceType, 0, defence.Quantity})
		}
	}
}

type OccupiedPlanet struct {
	PlanetConfig            string                             `json:"planet_config" example:"Planet2.json"`
	Position                models.PlanetPosition              `json:"position"`
	Distance                int                                `json:"distance" example:"14"`
	BasePlanet              bool                               `json:"base_planet" example:"true"`
	Resources               *Resources                         `json:"resources,omitempty"`
	Population              *Population                        `json:"population,omitempty"`
	Mines                   []buildings.Mine                   `json:"mines,omitempty"`
	Shields                 []buildings.Shield                 `json:"shields,omitempty"`
	PopulationControlCenter *buildings.PopulationControlCenter `json:"population_control_center,omitempty"`
	Defences                []military.Defence                 `json:"defences,omitempty"`
	DefenceShipCarriers     []military.DefenceShipCarrier      `json:"defence_ship_carriers,omitempty"`
	Ships                   []military.Ship                    `json:"ships,omitempty"`
	Scouts                  []military.Ship                    `json:"scouts,omitempty"`
	HomePlanet              bool                               `json:"home_planet" example:"true"`
	AttackMissions          []ActiveMission                    `json:"attack_missions,omitempty"`
	SpyMissions             []ActiveMission                    `json:"spy_missions,omitempty"`
	PlanetType              string                             `json:"planet_type" example:"BASE_PLANET"`
}

func (o *OccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser, customHomePlanetId string,
	attackMissions []repoModels.AttackMission, spyMissions []repoModels.SpyMission,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]constants.BuildingConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants, shipConstants map[string]constants.ShipConstants,
	speciesConstants constants.SpeciesConstants,
) {

	o.PlanetConfig = planetUni.PlanetConfig
	o.Position = planetUni.Position.Clone()
	o.Distance = planetUni.Distance
	if planetUser.BasePlanet {
		o.BasePlanet = true
		o.PlanetType = BasePlanet
		return
	}
	o.PlanetType = HomePlanet
	o.Resources = InitResources(planetUser)
	o.Population = InitPopulation(planetUser)
	o.Mines = buildings.InitAllMines(planetUni, planetUser,
		upgradeConstants[constants.WaterMiningPlant], upgradeConstants[constants.GrapheneMiningPlant],
		waterConstants, grapheneConstants)
	o.PopulationControlCenter = buildings.InitPopulationControlCenter(planetUser,
		upgradeConstants[constants.PopulationControlCenter], buildingConstants[constants.PopulationControlCenter])
	o.Shields = buildings.InitAllShields(planetUser, defenceConstants, upgradeConstants[constants.Shield])

	for _, unitName := range speciesConstants.AvailableUnits {
		if defenceConstant, ok := defenceConstants[unitName]; ok {
			if defenceConstant.Type == constants.Defender {
				d := military.Defence{}
				d.Init(unitName, planetUser.Defences[unitName], defenceConstant)
				o.Defences = append(o.Defences, d)
			}
		}
		if shipConstant, ok := shipConstants[unitName]; ok {
			if shipConstant.Type == constants.Scout {
				s := military.Ship{}
				s.InitScout(unitName, planetUser.Ships[unitName], spyMissions, shipConstants[unitName])
				o.Scouts = append(o.Scouts, s)
			} else {
				s := military.Ship{}
				s.Init(unitName, planetUser.Ships[unitName], attackMissions, shipConstants[unitName])
				o.Ships = append(o.Ships, s)
			}
		}
	}
	o.DefenceShipCarriers = military.InitAllDefenceShipCarriers(planetUser, defenceConstants)

	o.HomePlanet = planetUser.HomePlanet || planetUni.Position.Id == customHomePlanetId

	for _, attackMission := range attackMissions {
		activeMission := ActiveMission{}
		activeMission.InitAttackMission(attackMission)
		o.AttackMissions = append(o.AttackMissions, activeMission)
	}
	for _, spyMission := range spyMissions {
		activeMission := ActiveMission{}
		activeMission.InitSpyMission(spyMission)
		o.SpyMissions = append(o.SpyMissions, activeMission)
	}
}

type StaticPlanetData struct {
	PlanetConfig string                `json:"planet_config" example:"Planet2.json"`
	Position     models.PlanetPosition `json:"position"`
	SameSector   bool                  `json:"same_sector" example:"true"`
}

func (p *StaticPlanetData) Init(planetUni repoModels.PlanetUni, homeSectorId string) {
	p.PlanetConfig = planetUni.PlanetConfig
	p.SameSector = homeSectorId == planetUni.Position.SectorId()
	p.Position = planetUni.Position.Clone()
}

package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/buildings"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/controllers/models/researches"
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
		for shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Unavailable})
		}
		u.PlanetType = ResourcePlanet
	} else if planetType == constants.Base {
		for shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Invulnerable})
		}
		u.PlanetType = BasePlanet
		u.Invulnerable = true
	} else if planetType == constants.Bot {
		for shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Broken})
		}
		u.PlanetType = AbandonedPlanet
		u.Water = planetUser.Water.Amount
		u.Graphene = planetUser.Graphene.Amount
		for defenceType, defence := range planetUser.Defences {
			u.Defences = append(u.Defences, UnoccupiedPlanetDefence{defenceType, defence.Level, defence.Quantity})
		}
	} else {
		for shieldId := range shieldIds {
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
	Researches              []researches.Research              `json:"researches,omitempty"`
	PopulationControlCenter *buildings.PopulationControlCenter `json:"population_control_center,omitempty"`
	AttackProductionCenter  *buildings.UnitProductionCenter    `json:"attack_production_center,omitempty"`
	DefenceProductionCenter *buildings.UnitProductionCenter    `json:"defence_production_center,omitempty"`
	DiamondStorage          *buildings.ResourceStorage         `json:"diamond_storage,omitempty"`
	WaterPressureTank       *buildings.ResourceStorage         `json:"water_pressure_tank,omitempty"`
	ResearchLab             *buildings.ResearchLab             `json:"research_lab,omitempty"`
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
	buildingConstants map[string]map[string]map[string]interface{},
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	researchConstants map[string]constants.ResearchConstants,
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
	o.Population = InitPopulation(planetUser, militaryConstants)
	o.Mines = buildings.InitAllMines(planetUni, planetUser,
		upgradeConstants[constants.WaterMiningPlant], upgradeConstants[constants.GrapheneMiningPlant],
		waterConstants, grapheneConstants)
	o.PopulationControlCenter = buildings.InitPopulationControlCenter(planetUser,
		upgradeConstants[constants.PopulationControlCenter], buildingConstants[constants.PopulationControlCenter])
	o.AttackProductionCenter = buildings.InitAttackProductionCenter(planetUser,
		upgradeConstants[constants.AttackProductionCenter], buildingConstants[constants.AttackProductionCenter])
	o.DefenceProductionCenter = buildings.InitDefenceProductionCenter(planetUser,
		upgradeConstants[constants.DefenceProductionCenter], buildingConstants[constants.DefenceProductionCenter])
	o.DiamondStorage = buildings.InitDiamondStorage(planetUser,
		upgradeConstants[constants.DiamondStorage], buildingConstants[constants.DiamondStorage])
	o.WaterPressureTank = buildings.InitWaterPressureTank(planetUser,
		upgradeConstants[constants.WaterPressureTank], buildingConstants[constants.WaterPressureTank])
	o.Shields = buildings.InitAllShields(planetUser, buildingConstants[constants.Shield], upgradeConstants[constants.Shield])
	o.Researches = researches.InitAllResearches(planetUser, researchConstants)
	o.ResearchLab = buildings.InitResearchLab(planetUser, upgradeConstants[constants.ResearchLab], buildingConstants[constants.ResearchLab])

	for _, unitName := range speciesConstants.AvailableUnits {
		if planetUser.Defences[unitName].Level > 0 {
			if defenceConstant, ok := militaryConstants[unitName]; ok {
				if defenceConstant.Type == constants.Defender {
					d := military.Defence{}
					d.Init(unitName, planetUser.Defences[unitName], defenceConstant.Levels)
					o.Defences = append(o.Defences, d)
				}
			}
		}
		if planetUser.Ships[unitName].Level > 0 {
			if shipConstant, ok := militaryConstants[unitName]; ok {
				if shipConstant.Type == constants.Scout {
					s := military.Ship{}
					s.InitScout(unitName, shipConstant.Type, planetUser.Ships[unitName], spyMissions, shipConstant.Levels)
					o.Scouts = append(o.Scouts, s)
				} else {
					s := military.Ship{}
					s.Init(unitName, shipConstant.Type, planetUser.Ships[unitName], attackMissions, planetUser.DefenceShipCarriers, shipConstant.Levels)
					o.Ships = append(o.Ships, s)
				}
			}
		}
	}
	o.DefenceShipCarriers = military.InitAllDefenceShipCarriers(planetUser, militaryConstants)

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

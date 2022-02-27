package models

import (
	"github.com/themane/MMOServer/constants"
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
	PlanetConfig            string                  `json:"planet_config" example:"Planet2.json"`
	Position                models.PlanetPosition   `json:"position"`
	Distance                int                     `json:"distance" example:"14"`
	BasePlanet              bool                    `json:"base_planet" example:"true"`
	Resources               Resources               `json:"resources"`
	Population              Population              `json:"population"`
	Mines                   []Mine                  `json:"mines"`
	Shields                 []Shield                `json:"shields"`
	PopulationControlCenter PopulationControlCenter `json:"population_control_center"`
	IdleDefences            []Defence               `json:"idle_defences" bson:"idle_defences"`
	IdleDefenceShipCarriers []DefenceShipCarrier    `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	AvailableAttackShips    []Ship                  `json:"available_attack_ships" bson:"available_attack_ships"`
	Scouts                  []Ship                  `json:"scouts" bson:"scouts"`
	HomePlanet              bool                    `json:"home_planet" example:"true"`
	AttackMissions          []ActiveMission         `json:"attack_missions"`
	SpyMissions             []ActiveMission         `json:"spy_missions"`
	PlanetType              string                  `json:"planet_type" example:"BASE_PLANET"`
}

func (o *OccupiedPlanet) Init(planetUni repoModels.PlanetUni, planetUser repoModels.PlanetUser, customHomePlanetId string,
	attackMissions []repoModels.AttackMission, spyMissions []repoModels.SpyMission,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]constants.BuildingConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants, shipConstants map[string]constants.ShipConstants) {

	o.PlanetConfig = planetUni.PlanetConfig
	o.Position = planetUni.Position.Clone()
	o.Distance = planetUni.Distance
	if planetUser.BasePlanet {
		o.BasePlanet = true
		o.PlanetType = BasePlanet
		return
	}
	o.PlanetType = HomePlanet
	o.Resources.Init(planetUser)
	o.Population.Init(planetUser)
	for mineId := range planetUser.Mines {
		mine := Mine{}
		mine.Init(planetUni.Mines[mineId], planetUser,
			upgradeConstants[constants.WaterMiningPlant], upgradeConstants[constants.GrapheneMiningPlant],
			waterConstants, grapheneConstants,
		)
		o.Mines = append(o.Mines, mine)
	}
	o.PopulationControlCenter.Init(planetUser,
		upgradeConstants[constants.PopulationControlCenter], buildingConstants[constants.PopulationControlCenter])
	o.Shields = InitAllShields(planetUser, defenceConstants, upgradeConstants[constants.Shield])
	o.IdleDefences = InitAllIdleDefences(planetUser.Defences, defenceConstants)
	o.IdleDefenceShipCarriers = InitAllIdleDefenceShipCarriers(planetUser, defenceConstants[constants.Vikram], shipConstants)
	for shipName, shipUser := range planetUser.Ships {
		availableShips := planetUser.GetAvailableShip(shipName)
		if availableShips <= 0 {
			continue
		}
		s := Ship{}
		s.Init(shipName, availableShips, shipUser, shipConstants[shipName])
		if s.Type == constants.Scout {
			o.Scouts = append(o.Scouts, s)
		} else {
			o.AvailableAttackShips = append(o.AvailableAttackShips, s)
		}

	}
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

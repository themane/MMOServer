package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strings"
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
	if planetUni.Occupied == "PRIMITIVE" || planetUni.Occupied == "" {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Unavailable})
		}
	} else if planetUni.BasePlanet == true {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Invulnerable})
		}
		u.Invulnerable = true
	} else if strings.HasPrefix(planetUni.Occupied, "BOT") {
		for _, shieldId := range shieldIds {
			u.Shields = append(u.Shields, UnoccupiedPlanetShield{shieldId, constants.Broken})
		}
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
		u.Water = planetUser.Water.Amount
		u.Graphene = planetUser.Graphene.Amount
		for defenceType, defence := range planetUser.Defences {
			u.Defences = append(u.Defences, UnoccupiedPlanetDefence{defenceType, 0, defence.Quantity})
		}
	}
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
	AvailableAttackShips    []Ship                `json:"available_attack_ships" bson:"available_attack_ships"`
	Scouts                  []Ship                `json:"scouts" bson:"scouts"`
	Home                    bool                  `json:"home" example:"true"`
	Distance                int                   `json:"distance" example:"14"`
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
		if s.Type == constants.Scout {
			o.Scouts = append(o.Scouts, s)
		} else {
			o.AvailableAttackShips = append(o.AvailableAttackShips, s)
		}

	}
	o.Home = planetUser.Home
}

package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Shield struct {
	Id                        string                    `json:"_id" example:"SHLD101"`
	Level                     int                       `json:"level" example:"3"`
	BuildingState             BuildingState             `json:"building_state"`
	Workers                   int                       `json:"workers" example:"12"`
	NextLevelShieldAttributes NextLevelShieldAttributes `json:"next_level_attributes"`
	NextLevelRequirements     NextLevelRequirements     `json:"next_level_requirements"`
	DeployedDefences          []Defence                 `json:"deployed_defences"`
}

type Defence struct {
	Type             string `json:"type" example:"BOMBER"`
	Level            int    `json:"level" example:"1"`
	Quantity         int    `json:"quantity" example:"5"`
	HitPoints        int    `json:"hit_points" example:"400"`
	Armor            int    `json:"armor" example:"2"`
	MinAttack        int    `json:"min_attack" example:"10"`
	MaxAttack        int    `json:"max_attack" example:"12"`
	Range            int    `json:"range" example:"2"`
	SingleHitTargets int    `json:"single_hit_targets" example:"1"`
}

type DefenceShipCarrier struct {
	Id            string `json:"_id" example:"DSC001"`
	Level         int    `json:"level" example:"1"`
	HitPoints     int    `json:"hit_points" example:"400"`
	Armor         int    `json:"armor" example:"5"`
	DeployedShips []Ship `json:"deployed_ships"`
}

type NextLevelShieldAttributes struct {
	CurrentHitPoints int `json:"current_hit_points" example:"1"`
	NextHitPoints    int `json:"next_hit_points" example:"1"`
	MaxHitPoints     int `json:"max_hit_points" example:"12"`
}

func InitAllShields(planetUser models.PlanetUser,
	defenceConstants map[string]constants.DefenceConstants, shieldBuildingConstants constants.BuildingConstants) []Shield {

	var shields []Shield
	shieldIds := []string{"SHLD01", "SHLD02", "SHLD03"}
	for _, shieldId := range shieldIds {
		s := Shield{}
		s.Id = shieldId
		s.Level = planetUser.Buildings[shieldId].BuildingLevel
		s.BuildingState.Init(planetUser.Buildings[shieldId], shieldBuildingConstants)
		s.Workers = planetUser.Buildings[shieldId].Workers
		s.NextLevelRequirements.Init(planetUser.Buildings[shieldId].BuildingLevel, shieldBuildingConstants)
		s.NextLevelShieldAttributes.Init(planetUser.Buildings[shieldId].BuildingLevel, defenceConstants[constants.Shield])
		for defenceType, defenceUser := range planetUser.Defences {
			if deployedDefences, ok := defenceUser.GuardingShield[shieldId]; ok {
				d := Defence{}
				d.Init(defenceType, deployedDefences, defenceUser, defenceConstants[defenceType])
				s.DeployedDefences = append(s.DeployedDefences, d)
			}
		}
		shields = append(shields, s)
	}
	return shields
}

func InitAllIdleDefences(defencesUser map[string]models.Defence,
	defenceConstants map[string]constants.DefenceConstants) []Defence {

	var defences []Defence
	for defenceType, defenceUser := range defencesUser {
		deployedOnShields := 0
		for _, deployed := range defenceUser.GuardingShield {
			deployedOnShields += deployed
		}
		idleDefences := defenceUser.Quantity - deployedOnShields
		if idleDefences <= 0 {
			continue
		}
		d := Defence{}
		d.Init(defenceType, idleDefences, defenceUser, defenceConstants[defenceType])
		defences = append(defences, d)
	}
	return defences
}

func (d *Defence) Init(defenceType string, quantity int, defenceUser models.Defence, defenceConstants constants.DefenceConstants) {
	d.Type = defenceType
	d.Level = defenceUser.Level
	d.Quantity = quantity
	currentLevelString := strconv.Itoa(defenceUser.Level)
	d.HitPoints = defenceConstants.Levels[currentLevelString].HitPoints
	d.Armor = defenceConstants.Levels[currentLevelString].Armor
	d.MinAttack = defenceConstants.Levels[currentLevelString].MinAttack
	d.MaxAttack = defenceConstants.Levels[currentLevelString].MaxAttack
	d.Range = defenceConstants.Levels[currentLevelString].Range
	d.SingleHitTargets = defenceConstants.Levels[currentLevelString].SingleHitTargets
}

func InitAllIdleDefenceShipCarriers(planetUser models.PlanetUser,
	defenceShipCarrierConstants constants.DefenceConstants, shipConstants map[string]constants.ShipConstants) []DefenceShipCarrier {

	var defenceShipCarriers []DefenceShipCarrier
	for id, defenceShipCarrierUser := range planetUser.DefenceShipCarriers {
		if defenceShipCarrierUser.GuardingShield == "" {
			continue
		}
		d := DefenceShipCarrier{}
		d.Id = id
		d.Level = defenceShipCarrierUser.Level
		currentLevelString := strconv.Itoa(defenceShipCarrierUser.Level)
		d.HitPoints = defenceShipCarrierConstants.Levels[currentLevelString].HitPoints
		d.Armor = defenceShipCarrierConstants.Levels[currentLevelString].Armor
		for shipName, shipQuantity := range defenceShipCarrierUser.HostingShips {
			s := Ship{}
			s.Init(shipName, shipQuantity, planetUser.Ships[shipName], shipConstants[shipName])
			d.DeployedShips = append(d.DeployedShips, s)
		}
		defenceShipCarriers = append(defenceShipCarriers, d)
	}
	return defenceShipCarriers
}

func (n *NextLevelShieldAttributes) Init(currentLevel int, shieldConstants constants.DefenceConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(shieldConstants.MaxLevel)
	n.CurrentHitPoints = shieldConstants.Levels[currentLevelString].HitPoints
	n.MaxHitPoints = shieldConstants.Levels[maxLevelString].HitPoints
	if currentLevel+1 < shieldConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.NextHitPoints = shieldConstants.Levels[nextLevelString].HitPoints
	}
}

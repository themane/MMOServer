package military

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

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

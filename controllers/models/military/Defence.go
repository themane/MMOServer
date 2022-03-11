package military

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type DeployedDefence struct {
	Name  string `json:"name" example:"BOMBER"`
	Level int    `json:"level" example:"1"`
	Units int    `json:"units" example:"5"`
}

type Defence struct {
	Name             string `json:"name" example:"BOMBER"`
	Level            int    `json:"level" example:"1"`
	TotalUnits       int    `json:"total_units" example:"5"`
	IdleUnits        int    `json:"idle_units" example:"5"`
	RequiredSoldiers int    `json:"required_soldiers" example:"40"`
	RequiredWorkers  int    `json:"required_workers" example:"20"`
	HitPoints        int    `json:"hit_points" example:"400"`
	Armor            int    `json:"armor" example:"2"`
	MinAttack        int    `json:"min_attack" example:"10"`
	MaxAttack        int    `json:"max_attack" example:"12"`
	Range            int    `json:"range" example:"2"`
	SingleHitTargets int    `json:"single_hit_targets" example:"1"`
}

func (d *Defence) Init(unitName string, defenceUser models.Defence, defenceConstants constants.DefenceConstants) {
	d.Name = unitName
	d.Level = defenceUser.Level
	d.TotalUnits = defenceUser.Quantity

	deployedUnits := 0
	for _, quantity := range defenceUser.GuardingShield {
		deployedUnits += quantity
	}
	d.IdleUnits = d.TotalUnits - deployedUnits

	currentLevelString := strconv.Itoa(defenceUser.Level)
	d.RequiredSoldiers = defenceConstants.Levels[currentLevelString].RequiredSoldiers
	d.RequiredWorkers = defenceConstants.Levels[currentLevelString].RequiredWorkers
	d.HitPoints = defenceConstants.Levels[currentLevelString].HitPoints
	d.Armor = defenceConstants.Levels[currentLevelString].Armor
	d.MinAttack = defenceConstants.Levels[currentLevelString].MinAttack
	d.MaxAttack = defenceConstants.Levels[currentLevelString].MaxAttack
	d.Range = defenceConstants.Levels[currentLevelString].Range
	d.SingleHitTargets = defenceConstants.Levels[currentLevelString].SingleHitTargets
}

type DeployedDefenceShipCarrier struct {
	Id            string         `json:"_id" example:"DSC001"`
	Name          string         `json:"name" example:"VIKRAM"`
	Level         int            `json:"level" example:"1"`
	DeployedShips []DeployedShip `json:"deployed_ships"`
}

type DefenceShipCarrier struct {
	Id               string         `json:"_id" example:"DSC001"`
	Name             string         `json:"name" example:"VIKRAM"`
	Level            int            `json:"level" example:"1"`
	RequiredSoldiers int            `json:"required_soldiers" example:"40"`
	RequiredWorkers  int            `json:"required_workers" example:"20"`
	HitPoints        int            `json:"hit_points" example:"400"`
	Armor            int            `json:"armor" example:"5"`
	DeployedShips    []DeployedShip `json:"deployed_ships"`
	Idle             bool           `json:"idle" example:"true"`
}

func InitAllDefenceShipCarriers(planetUser models.PlanetUser,
	defenceConstants map[string]constants.DefenceConstants) []DefenceShipCarrier {

	var defenceShipCarriers []DefenceShipCarrier
	for id, defenceShipCarrierUser := range planetUser.DefenceShipCarriers {
		currentLevelString := strconv.Itoa(defenceShipCarrierUser.Level)
		d := DefenceShipCarrier{
			Id:               id,
			Name:             defenceShipCarrierUser.Name,
			Level:            defenceShipCarrierUser.Level,
			RequiredSoldiers: defenceConstants[defenceShipCarrierUser.Name].Levels[currentLevelString].RequiredSoldiers,
			RequiredWorkers:  defenceConstants[defenceShipCarrierUser.Name].Levels[currentLevelString].RequiredWorkers,
			HitPoints:        defenceConstants[defenceShipCarrierUser.Name].Levels[currentLevelString].HitPoints,
			Armor:            defenceConstants[defenceShipCarrierUser.Name].Levels[currentLevelString].Armor,
			DeployedShips:    GetDeployedShips(planetUser.Ships, defenceShipCarrierUser.HostingShips),
			Idle:             defenceShipCarrierUser.GuardingShield != "",
		}
		defenceShipCarriers = append(defenceShipCarriers, d)
	}
	return defenceShipCarriers
}

func GetDeployedShips(ships map[string]models.Ship, hostingShips map[string]int) []DeployedShip {
	var deployedShips []DeployedShip
	for unitName, units := range hostingShips {
		s := DeployedShip{
			Name:  unitName,
			Level: ships[unitName].Level,
			Units: units,
		}
		deployedShips = append(deployedShips, s)
	}
	return deployedShips
}

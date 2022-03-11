package military

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type DeployedShip struct {
	Name  string `json:"name" example:"ANUJ"`
	Level int    `json:"level" example:"2"`
	Units int    `json:"units" example:"30"`
}

type Ship struct {
	Name             string `json:"name" example:"ANUJ"`
	Type             string `json:"type" example:"ATTACKER"`
	Level            int    `json:"level" example:"2"`
	TotalUnits       int    `json:"total_units" example:"30"`
	AvailableUnits   int    `json:"available_units" example:"15"`
	RequiredSoldiers int    `json:"required_soldiers" example:"40"`
	RequiredWorkers  int    `json:"required_workers" example:"20"`
	HitPoints        int    `json:"hit_points" example:"40"`
	Armor            int    `json:"armor" example:"2"`
	ResourceCapacity int    `json:"resource_capacity" example:"40"`
	WorkerCapacity   int    `json:"worker_capacity" example:"20"`
	MinAttack        int    `json:"min_attack" example:"5"`
	MaxAttack        int    `json:"max_attack" example:"7"`
	Range            int    `json:"range" example:"2"`
	Speed            int    `json:"speed" example:"600"`
}

func (s *Ship) Init(unitName string, shipUser models.Ship, attackMissions []models.AttackMission, shipConstants constants.ShipConstants) {
	s.Name = unitName
	s.Level = shipUser.Level
	s.TotalUnits = shipUser.Quantity

	deployedUnits := 0
	for _, mission := range attackMissions {
		for _, shieldFormation := range mission.Formation {
			for _, pointFormation := range shieldFormation {
				deployedUnits += pointFormation[unitName]
			}
		}
	}
	s.AvailableUnits = s.TotalUnits - deployedUnits

	currentLevelString := strconv.Itoa(shipUser.Level)
	s.RequiredSoldiers = shipConstants.Levels[currentLevelString].RequiredSoldiers
	s.RequiredWorkers = shipConstants.Levels[currentLevelString].RequiredWorkers
	s.HitPoints = shipConstants.Levels[currentLevelString].HitPoints
	s.Armor = shipConstants.Levels[currentLevelString].Armor
	s.ResourceCapacity = shipConstants.Levels[currentLevelString].ResourceCapacity
	s.WorkerCapacity = shipConstants.Levels[currentLevelString].WorkerCapacity
	s.MinAttack = shipConstants.Levels[currentLevelString].MinAttack
	s.MaxAttack = shipConstants.Levels[currentLevelString].MaxAttack
	s.Range = shipConstants.Levels[currentLevelString].Range
	s.Speed = shipConstants.Levels[currentLevelString].Speed
}

func (s *Ship) InitScout(unitName string, shipUser models.Ship, spyMissions []models.SpyMission, shipConstants constants.ShipConstants) {
	s.Name = unitName
	s.Level = shipUser.Level
	s.TotalUnits = shipUser.Quantity

	deployedUnits := 0
	for _, mission := range spyMissions {
		deployedUnits += mission.Scouts[unitName]
	}
	s.AvailableUnits = s.TotalUnits - deployedUnits

	currentLevelString := strconv.Itoa(shipUser.Level)
	s.RequiredSoldiers = shipConstants.Levels[currentLevelString].RequiredSoldiers
	s.RequiredWorkers = shipConstants.Levels[currentLevelString].RequiredWorkers
	s.HitPoints = shipConstants.Levels[currentLevelString].HitPoints
	s.Armor = shipConstants.Levels[currentLevelString].Armor
	s.ResourceCapacity = shipConstants.Levels[currentLevelString].ResourceCapacity
	s.WorkerCapacity = shipConstants.Levels[currentLevelString].WorkerCapacity
	s.MinAttack = shipConstants.Levels[currentLevelString].MinAttack
	s.MaxAttack = shipConstants.Levels[currentLevelString].MaxAttack
	s.Range = shipConstants.Levels[currentLevelString].Range
	s.Speed = shipConstants.Levels[currentLevelString].Speed
}

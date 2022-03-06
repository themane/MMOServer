package military

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Ship struct {
	Name             string `json:"name" example:"ANUJ"`
	Type             string `json:"type" example:"ATTACKER"`
	Level            int    `json:"level" example:"2"`
	Quantity         int    `json:"quantity" example:"15"`
	HitPoints        int    `json:"hit_points" example:"40"`
	Armor            int    `json:"armor" example:"2"`
	ResourceCapacity int    `json:"resource_capacity" example:"40"`
	WorkerCapacity   int    `json:"worker_capacity" example:"20"`
	MinAttack        int    `json:"min_attack" example:"5"`
	MaxAttack        int    `json:"max_attack" example:"7"`
	Range            int    `json:"range" example:"2"`
	Speed            int    `json:"speed" example:"600"`
}

func (s *Ship) Init(shipName string, shipQuantity int, shipUser repoModels.Ship, shipConstants constants.ShipConstants) {
	s.Name = shipName
	s.Type = shipConstants.Type
	s.Level = shipUser.Level
	s.Quantity = shipQuantity
	currentLevelString := strconv.Itoa(shipUser.Level)
	s.HitPoints = shipConstants.Levels[currentLevelString].HitPoints
	s.Armor = shipConstants.Levels[currentLevelString].Armor
	s.ResourceCapacity = shipConstants.Levels[currentLevelString].ResourceCapacity
	s.WorkerCapacity = shipConstants.Levels[currentLevelString].WorkerCapacity
	s.MinAttack = shipConstants.Levels[currentLevelString].MinAttack
	s.MaxAttack = shipConstants.Levels[currentLevelString].MaxAttack
	s.Range = shipConstants.Levels[currentLevelString].Range
	s.Speed = shipConstants.Levels[currentLevelString].Speed
}

package military

import (
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type DeployedShip struct {
	Name  string `json:"name" example:"ANUJ"`
	Level int    `json:"level" example:"2"`
	Units int    `json:"units" example:"30"`
}

type Ship struct {
	Name                 string                `json:"name" example:"ANUJ"`
	Type                 string                `json:"type" example:"ATTACKER"`
	Level                int                   `json:"level" example:"2"`
	TotalUnits           int                   `json:"total_units" example:"30"`
	AvailableUnits       int                   `json:"available_units" example:"15"`
	ShipAttributes       models.ShipAttributes `json:"ship_attributes"`
	CreationRequirements models.Requirements   `json:"creation_requirements"`
	DestructionReturns   models.Returns        `json:"destruction_returns"`
	UnderConstruction    *UnderConstruction    `json:"under_construction,omitempty"`
}

func (s *Ship) Init(unitName string, shipUser repoModels.Ship,
	attackMissions []repoModels.AttackMission, defenceShipCarriers map[string]repoModels.DefenceShipCarrier,
	shipConstants map[string]map[string]interface{}) {

	s.Name = unitName
	s.Level = shipUser.Level
	s.TotalUnits = shipUser.Quantity
	s.AvailableUnits = repoModels.GetAvailableShips(unitName, attackMissions, defenceShipCarriers, shipUser.Quantity)

	if shipUser.Level > 0 {
		currentLevelString := strconv.Itoa(shipUser.Level)
		s.ShipAttributes.Init(shipConstants[currentLevelString])
		s.CreationRequirements.Init(shipConstants[currentLevelString])
		s.DestructionReturns.InitDestructionReturns(shipConstants[currentLevelString])
		s.UnderConstruction = InitUnderConstruction(shipUser.UnderConstruction, shipConstants[currentLevelString])
	}
}

func (s *Ship) InitScout(unitName string, shipUser repoModels.Ship, spyMissions []repoModels.SpyMission, shipConstants map[string]map[string]interface{}) {
	s.Name = unitName
	s.Level = shipUser.Level
	s.TotalUnits = shipUser.Quantity
	s.AvailableUnits = repoModels.GetAvailableScouts(unitName, spyMissions, shipUser.Quantity)

	currentLevelString := strconv.Itoa(shipUser.Level)
	if shipUser.Level > 0 {
		s.ShipAttributes.Init(shipConstants[currentLevelString])
		s.CreationRequirements.Init(shipConstants[currentLevelString])
		s.DestructionReturns.InitDestructionReturns(shipConstants[currentLevelString])
		s.UnderConstruction = InitUnderConstruction(shipUser.UnderConstruction, shipConstants[currentLevelString])
	}
}

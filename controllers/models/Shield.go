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
}

type NextLevelShieldAttributes struct {
	CurrentHitPoints int `json:"current_hit_points" example:"1"`
	NextHitPoints    int `json:"next_hit_points" example:"1"`
	MaxHitPoints     int `json:"max_hit_points" example:"12"`
}

func initAllShields(planetUser models.PlanetUser,
	shieldConstants constants.DefenceConstants, shieldBuildingConstants constants.BuildingConstants) []Shield {

	var shields []Shield
	for shieldId, shieldLevel := range planetUser.Shields {
		s := Shield{}
		s.Id = shieldId
		s.Level = shieldLevel
		s.BuildingState.Init(planetUser.Buildings[shieldId], shieldBuildingConstants)
		s.Workers = planetUser.Buildings[shieldId].Workers
		s.NextLevelRequirements.Init(shieldLevel, shieldBuildingConstants)
		s.NextLevelShieldAttributes.Init(shieldLevel, shieldConstants)
		shields = append(shields, s)
	}
	return shields
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

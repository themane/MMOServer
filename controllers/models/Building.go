package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type BuildingState struct {
	State            string        `json:"state" example:"WORKING"`
	MinutesRemaining int           `json:"minutes_remaining_per_worker" example:"1440"`
	CancelReturns    CancelReturns `json:"cancel_returns"`
}

type CancelReturns struct {
	WaterReturned    int `json:"water_returned" example:"5"`
	GrapheneReturned int `json:"graphene_returned" example:"101"`
	ShelioReturned   int `json:"shelio_returned" example:"0"`
}

type NextLevelRequirements struct {
	GrapheneRequired         int `json:"graphene_required" example:"101"`
	WaterRequired            int `json:"water_required" example:"5"`
	ShelioRequired           int `json:"shelio_required" example:"0"`
	MinutesRequiredPerWorker int `json:"minutes_required_per_worker" example:"1440"`
}

func (b *BuildingState) Init(building models.Building, buildingConstants constants.BuildingConstants) {
	if building.BuildingMinutesPerWorker > 0 {
		b.State = constants.UpgradingState
	} else {
		b.State = constants.WorkingState
	}
	b.MinutesRemaining = building.BuildingMinutesPerWorker
	b.CancelReturns.Init(building.BuildingMinutesPerWorker, building.BuildingLevel, buildingConstants)
}

func (c *CancelReturns) Init(buildingMinutesPerWorker int, buildingLevel int, buildingConstants constants.BuildingConstants) {
	nextLevelString := strconv.Itoa(buildingLevel + 1)
	ratio := buildingMinutesPerWorker / buildingConstants.Levels[nextLevelString].MinutesRequired

	c.WaterReturned = buildingConstants.Levels[nextLevelString].WaterRequired * ratio
	c.GrapheneReturned = buildingConstants.Levels[nextLevelString].GrapheneRequired * ratio
	c.ShelioReturned = buildingConstants.Levels[nextLevelString].ShelioRequired * ratio
}

func (n *NextLevelRequirements) Init(currentLevel int, buildingConstants constants.BuildingConstants) {
	if currentLevel+1 < buildingConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.GrapheneRequired = buildingConstants.Levels[nextLevelString].GrapheneRequired
		n.WaterRequired = buildingConstants.Levels[nextLevelString].WaterRequired
		n.ShelioRequired = buildingConstants.Levels[nextLevelString].ShelioRequired
		n.MinutesRequiredPerWorker = buildingConstants.Levels[nextLevelString].MinutesRequired
	}
}

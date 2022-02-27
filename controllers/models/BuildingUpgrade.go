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

func (b *BuildingState) Init(building models.Building, upgradeConstants constants.UpgradeConstants) {
	if building.BuildingMinutesPerWorker > 0 {
		b.State = constants.UpgradingState
		b.MinutesRemaining = building.BuildingMinutesPerWorker
		b.CancelReturns.Init(building.BuildingMinutesPerWorker, building.BuildingLevel, upgradeConstants)
	} else {
		b.State = constants.WorkingState
	}
}

func (c *CancelReturns) Init(buildingMinutesPerWorker int, buildingLevel int, upgradeConstants constants.UpgradeConstants) {
	if buildingLevel < upgradeConstants.MaxLevel {
		nextLevelString := strconv.Itoa(buildingLevel + 1)
		ratio := buildingMinutesPerWorker / upgradeConstants.Levels[nextLevelString].MinutesRequired

		c.WaterReturned = upgradeConstants.Levels[nextLevelString].WaterRequired * ratio
		c.GrapheneReturned = upgradeConstants.Levels[nextLevelString].GrapheneRequired * ratio
		c.ShelioReturned = upgradeConstants.Levels[nextLevelString].ShelioRequired * ratio
	}
}

func (n *NextLevelRequirements) Init(currentLevel int, upgradeConstants constants.UpgradeConstants) {
	if currentLevel+1 < upgradeConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.GrapheneRequired = upgradeConstants.Levels[nextLevelString].GrapheneRequired
		n.WaterRequired = upgradeConstants.Levels[nextLevelString].WaterRequired
		n.ShelioRequired = upgradeConstants.Levels[nextLevelString].ShelioRequired
		n.MinutesRequiredPerWorker = upgradeConstants.Levels[nextLevelString].MinutesRequired
	}
}

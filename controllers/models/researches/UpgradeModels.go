package researches

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"math"
	"strconv"
)

type State struct {
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
	GrapheneRequired         float64              `json:"graphene_required" example:"101"`
	WaterRequired            float64              `json:"water_required" example:"5"`
	ShelioRequired           float64              `json:"shelio_required" example:"0"`
	SpecialRequirements      []SpecialRequirement `json:"special_requirements"`
	MinutesRequiredPerWorker float64              `json:"minutes_required_per_worker" example:"1440"`
}

type SpecialRequirement struct {
	Fulfilled   bool   `json:"fulfilled" example:"true"`
	Description string `json:"description"`
}

func (b *State) Init(research models.ResearchUser, researchConstants constants.ResearchConstants) {
	if research.ResearchMinutesPerWorker > 0 {
		b.State = constants.UpgradingState
		b.MinutesRemaining = research.ResearchMinutesPerWorker
		b.CancelReturns.Init(research.ResearchMinutesPerWorker, research.Level, researchConstants)
	} else {
		b.State = constants.WorkingState
	}
}

func (c *CancelReturns) Init(minutesRemaining int, currentLevel int, researchConstants constants.ResearchConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	ratio := float64(minutesRemaining) / researchConstants.Requirements[currentLevelString]["minutes_required"].(float64)

	c.WaterReturned = int(math.Floor(researchConstants.Requirements[currentLevelString]["water_required"].(float64) * ratio))
	c.GrapheneReturned = int(math.Floor(researchConstants.Requirements[currentLevelString]["graphene_required"].(float64) * ratio))
	c.ShelioReturned = int(math.Floor(researchConstants.Requirements[currentLevelString]["shelio_required"].(float64) * ratio))
}

func (n *NextLevelRequirements) Init(currentLevel int, researchConstants constants.ResearchConstants) {
	nextLevelString := strconv.Itoa(currentLevel + 1)
	for requirementName, requirement := range researchConstants.Requirements[nextLevelString] {
		if requirementName == "graphene_required" {
			n.GrapheneRequired = requirement.(float64)
		} else if requirementName == "water_required" {
			n.WaterRequired = requirement.(float64)
		} else if requirementName == "shelio_required" {
			n.ShelioRequired = requirement.(float64)
		} else if requirementName == "minutes_required" {
			n.MinutesRequiredPerWorker = requirement.(float64)
		} else {
			s := SpecialRequirement{
				Description: requirement.(map[string]interface{})["description"].(string),
				Fulfilled:   true,
			}
			n.SpecialRequirements = append(n.SpecialRequirements, s)
		}
	}
}

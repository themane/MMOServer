package researches

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Research struct {
	Name                  string                 `json:"name" example:"POPULATION_CONTROL_CENTER"`
	Description           string                 `json:"description"`
	Level                 int                    `json:"level" example:"3"`
	Bonus                 map[string]interface{} `json:"bonus"`
	NextLevelBonus        map[string]interface{} `json:"next_level_bonus"`
	ResearchState         State                  `json:"research_state"`
	NextLevelRequirements *NextLevelRequirements `json:"next_level_requirements"`
}

func InitAllResearches(planetUser repoModels.PlanetUser,
	researchConstants map[string]constants.ResearchConstants) []Research {

	var researches []Research
	for researchName, researchConstant := range researchConstants {
		var currentLevel int
		research := planetUser.GetResearch(researchName)
		if research == nil {
			currentLevel = 0
		} else {
			currentLevel = research.Level
		}
		r := Research{
			Name:        researchName,
			Description: researchConstant.Description,
			Level:       currentLevel,
		}
		if currentLevel > 0 {
			r.Bonus = researchConstant.Bonus[strconv.Itoa(currentLevel)]
		}
		if currentLevel < researchConstant.MaxLevel {
			r.NextLevelBonus = researchConstant.Bonus[strconv.Itoa(currentLevel+1)]
			r.NextLevelRequirements = &NextLevelRequirements{}
			r.NextLevelRequirements.Init(currentLevel, researchConstant)
		}
		if currentLevel == 0 {
			r.Bonus = map[string]interface{}{}
			for k := range r.NextLevelBonus {
				r.Bonus[k] = 0
			}
		}
		r.ResearchState.Init(research, researchConstant)
		researches = append(researches, r)
	}
	return researches
}

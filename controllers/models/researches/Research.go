package researches

import (
	"fmt"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Research struct {
	Name                  string                 `json:"name" example:"POPULATION_CONTROL_CENTER"`
	Description           string                 `json:"description"`
	Level                 int                    `json:"level" example:"3"`
	Bonus                 []ResearchBonus        `json:"bonus"`
	NextLevelBonus        []ResearchBonus        `json:"next_level_bonus"`
	ResearchState         State                  `json:"research_state"`
	NextLevelRequirements *NextLevelRequirements `json:"next_level_requirements"`
}

type ResearchBonus struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
			for k, v := range researchConstant.Bonus[strconv.Itoa(currentLevel)] {
				value := fmt.Sprintf("%#v", v)
				r.Bonus = append(r.Bonus, ResearchBonus{k, value})
			}
		}
		if currentLevel < researchConstant.MaxLevel {
			for k, v := range researchConstant.Bonus[strconv.Itoa(currentLevel+1)] {
				value := fmt.Sprintf("%#v", v)
				r.NextLevelBonus = append(r.NextLevelBonus, ResearchBonus{k, value})
				if currentLevel == 0 {
					r.Bonus = append(r.Bonus, ResearchBonus{k, "0"})
				}
			}
			r.NextLevelRequirements = &NextLevelRequirements{}
			r.NextLevelRequirements.Init(currentLevel, researchConstant)
		}
		r.ResearchState.Init(research, researchConstant)
		researches = append(researches, r)
	}
	return researches
}

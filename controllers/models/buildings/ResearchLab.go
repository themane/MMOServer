package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type ResearchLab struct {
	BuildingId            string                `json:"building_id" example:"RESEARCH_LAB"`
	Level                 int                   `json:"level" example:"3"`
	Workers               int                   `json:"workers" example:"12"`
	Soldiers              int                   `json:"soldiers" example:"15"`
	BuildingState         State                 `json:"building_state"`
	BuildingAttributes    ResearchLabAttributes `json:"building_attributes"`
	NextLevelRequirements NextLevelRequirements `json:"next_level_requirements"`
}

type ResearchLabAttributes struct {
	WorkersMaxLimit  FloatBuildingAttributes `json:"workers_max_limit"`
	SoldiersMaxLimit FloatBuildingAttributes `json:"soldiers_max_limit"`
}

func InitResearchLab(planetUser repoModels.PlanetUser,
	upgradeConstants constants.UpgradeConstants,
	buildingConstants map[string]map[string]interface{}) *ResearchLab {

	r := new(ResearchLab)
	r.BuildingId = constants.ResearchLab
	r.Level = planetUser.Buildings[constants.ResearchLab].BuildingLevel
	r.Workers = planetUser.Buildings[constants.ResearchLab].Workers
	r.BuildingState.Init(planetUser.Buildings[constants.ResearchLab], upgradeConstants)
	r.BuildingAttributes.Init(planetUser.Buildings[constants.ResearchLab].BuildingLevel, upgradeConstants.MaxLevel, buildingConstants)
	r.NextLevelRequirements.Init(planetUser.Buildings[constants.ResearchLab].BuildingLevel, upgradeConstants)
	return r
}

func (p *ResearchLabAttributes) Init(currentLevel int, maxLevel int,
	researchLabBuildingConstants map[string]map[string]interface{}) {

	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		p.WorkersMaxLimit.Current = researchLabBuildingConstants[currentLevelString]["workers_max_limit"].(float64)
		p.SoldiersMaxLimit.Current = researchLabBuildingConstants[currentLevelString]["soldiers_max_limit"].(float64)
	}
	maxLevelString := strconv.Itoa(maxLevel)
	p.WorkersMaxLimit.Max = researchLabBuildingConstants[maxLevelString]["workers_max_limit"].(float64)
	p.SoldiersMaxLimit.Max = researchLabBuildingConstants[maxLevelString]["soldiers_max_limit"].(float64)

	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		p.WorkersMaxLimit.Next = researchLabBuildingConstants[nextLevelString]["workers_max_limit"].(float64)
		p.SoldiersMaxLimit.Next = researchLabBuildingConstants[nextLevelString]["soldiers_max_limit"].(float64)
	}
}

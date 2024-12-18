package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type ResearchLab struct {
	BuildingId            string                            `json:"building_id" example:"RESEARCH_LAB"`
	Level                 int                               `json:"level" example:"3"`
	Workers               int                               `json:"workers" example:"12"`
	Soldiers              int                               `json:"soldiers" example:"15"`
	BuildingState         repoModels.State                  `json:"building_state"`
	BuildingAttributes    ResearchLabAttributes             `json:"building_attributes"`
	NextLevelRequirements *repoModels.NextLevelRequirements `json:"next_level_requirements"`
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
	if planetUser.GetBuilding(constants.ResearchLab) != nil {
		r.Level = planetUser.GetBuilding(constants.ResearchLab).BuildingLevel
		r.Workers = planetUser.GetBuilding(constants.ResearchLab).Workers
		r.Soldiers = planetUser.GetBuilding(constants.ResearchLab).Soldiers
	}
	r.BuildingState.Init(planetUser.GetBuilding(constants.ResearchLab), upgradeConstants)
	r.BuildingAttributes.Init(r.Level, upgradeConstants.MaxLevel, buildingConstants)
	if r.Level < upgradeConstants.MaxLevel {
		r.NextLevelRequirements = &repoModels.NextLevelRequirements{}
		r.NextLevelRequirements.Init(r.Level, upgradeConstants)
	}
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
	} else {
		p.WorkersMaxLimit.Next = researchLabBuildingConstants[maxLevelString]["workers_max_limit"].(float64)
		p.SoldiersMaxLimit.Next = researchLabBuildingConstants[maxLevelString]["soldiers_max_limit"].(float64)
	}
}

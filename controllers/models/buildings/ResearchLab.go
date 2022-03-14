package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type ResearchLab struct {
	BuildingId            string                `json:"building_id" example:"RESEARCH_LAB"`
	Level                 int                   `json:"level" example:"3"`
	Workers               int                   `json:"workers" example:"12"`
	BuildingState         State                 `json:"building_state"`
	NextLevelRequirements NextLevelRequirements `json:"next_level_requirements"`
}

func InitResearchLab(planetUser repoModels.PlanetUser,
	upgradeConstants constants.UpgradeConstants) *ResearchLab {

	r := new(ResearchLab)
	r.BuildingId = constants.ResearchLab
	r.Level = planetUser.Buildings[constants.ResearchLab].BuildingLevel
	r.Workers = planetUser.Buildings[constants.ResearchLab].Workers
	r.BuildingState.Init(planetUser.Buildings[constants.ResearchLab], upgradeConstants)
	r.NextLevelRequirements.Init(planetUser.Buildings[constants.ResearchLab].BuildingLevel, upgradeConstants)
	return r
}

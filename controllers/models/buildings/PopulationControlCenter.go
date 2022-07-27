package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PopulationControlCenter struct {
	BuildingId            string                            `json:"building_id" example:"POPULATION_CONTROL_CENTER"`
	Level                 int                               `json:"level" example:"3"`
	Workers               int                               `json:"workers" example:"12"`
	BuildingState         repoModels.State                  `json:"building_state"`
	BuildingAttributes    PopulationControlCenterAttributes `json:"building_attributes"`
	NextLevelRequirements *repoModels.NextLevelRequirements `json:"next_level_requirements"`
}
type PopulationControlCenterAttributes struct {
	MaxPopulationGenerationRate        FloatBuildingAttributes `json:"max_population_generation_rate"`
	PopulationGenerationRateMultiplier FloatBuildingAttributes `json:"population_generation_rate_multiplier"`
	MinimumWorkersRequired             FloatBuildingAttributes `json:"minimum_workers_required"`
	WorkersMaxLimit                    FloatBuildingAttributes `json:"workers_max_limit"`
}

func InitPopulationControlCenter(planetUser repoModels.PlanetUser,
	populationControlCenterUpgradeConstants constants.UpgradeConstants,
	populationControlCenterBuildingConstants map[string]map[string]interface{}) *PopulationControlCenter {

	p := new(PopulationControlCenter)
	p.BuildingId = constants.PopulationControlCenter
	populationControlCenter := planetUser.GetBuilding(constants.PopulationControlCenter)
	p.Level = populationControlCenter.BuildingLevel
	p.Workers = populationControlCenter.Workers
	p.BuildingState.Init(*populationControlCenter, populationControlCenterUpgradeConstants)
	p.BuildingAttributes.Init(populationControlCenter.BuildingLevel,
		populationControlCenterUpgradeConstants.MaxLevel, populationControlCenterBuildingConstants)
	if p.Level < populationControlCenterUpgradeConstants.MaxLevel {
		p.NextLevelRequirements = &repoModels.NextLevelRequirements{}
		p.NextLevelRequirements.Init(populationControlCenter.BuildingLevel, populationControlCenterUpgradeConstants)
	}
	return p
}

func (p *PopulationControlCenterAttributes) Init(currentLevel int, maxLevel int,
	populationControlCenterBuildingConstants map[string]map[string]interface{}) {

	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		p.MaxPopulationGenerationRate.Current = populationControlCenterBuildingConstants[currentLevelString]["max_population_generation_rate"].(float64)
		p.PopulationGenerationRateMultiplier.Current = populationControlCenterBuildingConstants[currentLevelString]["population_generation_rate_multiplier"].(float64)
		p.MinimumWorkersRequired.Current = populationControlCenterBuildingConstants[currentLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Current = populationControlCenterBuildingConstants[currentLevelString]["workers_max_limit"].(float64)
	}
	maxLevelString := strconv.Itoa(maxLevel)
	workersMaxLimit := populationControlCenterBuildingConstants[maxLevelString]["workers_max_limit"].(float64)
	p.MaxPopulationGenerationRate.Max = models.GetMaxPopulationGenerationRate(populationControlCenterBuildingConstants[maxLevelString], workersMaxLimit)
	p.WorkersMaxLimit.Max = populationControlCenterBuildingConstants[maxLevelString]["workers_max_limit"].(float64)

	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		p.MaxPopulationGenerationRate.Next = populationControlCenterBuildingConstants[nextLevelString]["max_population_generation_rate"].(float64)
		p.PopulationGenerationRateMultiplier.Next = populationControlCenterBuildingConstants[nextLevelString]["population_generation_rate_multiplier"].(float64)
		p.MinimumWorkersRequired.Next = populationControlCenterBuildingConstants[nextLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Next = populationControlCenterBuildingConstants[nextLevelString]["workers_max_limit"].(float64)
	} else {
		p.MaxPopulationGenerationRate.Next = populationControlCenterBuildingConstants[maxLevelString]["max_population_generation_rate"].(float64)
		p.PopulationGenerationRateMultiplier.Next = populationControlCenterBuildingConstants[maxLevelString]["population_generation_rate_multiplier"].(float64)
		p.MinimumWorkersRequired.Next = populationControlCenterBuildingConstants[maxLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Next = populationControlCenterBuildingConstants[maxLevelString]["workers_max_limit"].(float64)
	}
}

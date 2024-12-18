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
	if planetUser.GetBuilding(constants.PopulationControlCenter) != nil {
		p.Level = planetUser.GetBuilding(constants.PopulationControlCenter).BuildingLevel
		p.Workers = planetUser.GetBuilding(constants.PopulationControlCenter).Workers
	}
	p.BuildingState.Init(planetUser.GetBuilding(constants.PopulationControlCenter), populationControlCenterUpgradeConstants)
	p.BuildingAttributes.Init(p.Level, p.Workers,
		populationControlCenterUpgradeConstants.MaxLevel, populationControlCenterBuildingConstants)
	if p.Level < populationControlCenterUpgradeConstants.MaxLevel {
		p.NextLevelRequirements = &repoModels.NextLevelRequirements{}
		p.NextLevelRequirements.Init(p.Level, populationControlCenterUpgradeConstants)
	}
	return p
}

func (p *PopulationControlCenterAttributes) Init(currentLevel int, workersDeployed int, maxLevel int,
	populationControlCenterBuildingConstants map[string]map[string]interface{}) {

	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		p.PopulationGenerationRateMultiplier.Current = populationControlCenterBuildingConstants[currentLevelString]["population_generation_rate_multiplier"].(float64)
		p.MaxPopulationGenerationRate.Current =
			models.MaxSelectablePopulationGenerationRate(populationControlCenterBuildingConstants[currentLevelString], workersDeployed)
		p.MinimumWorkersRequired.Current = populationControlCenterBuildingConstants[currentLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Current = populationControlCenterBuildingConstants[currentLevelString]["workers_max_limit"].(float64)
	}
	maxLevelString := strconv.Itoa(maxLevel)
	p.PopulationGenerationRateMultiplier.Max = populationControlCenterBuildingConstants[maxLevelString]["population_generation_rate_multiplier"].(float64)
	p.MaxPopulationGenerationRate.Max = models.MaxPopulationGenerationRate(populationControlCenterBuildingConstants[maxLevelString])
	p.MinimumWorkersRequired.Max = populationControlCenterBuildingConstants[maxLevelString]["workers_required"].(float64)
	p.WorkersMaxLimit.Max = populationControlCenterBuildingConstants[maxLevelString]["workers_max_limit"].(float64)

	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		p.PopulationGenerationRateMultiplier.Next = populationControlCenterBuildingConstants[nextLevelString]["population_generation_rate_multiplier"].(float64)
		p.MaxPopulationGenerationRate.Next = models.MaxPopulationGenerationRate(populationControlCenterBuildingConstants[nextLevelString])
		p.MinimumWorkersRequired.Next = populationControlCenterBuildingConstants[nextLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Next = populationControlCenterBuildingConstants[nextLevelString]["workers_max_limit"].(float64)
	} else {
		p.PopulationGenerationRateMultiplier.Next = populationControlCenterBuildingConstants[maxLevelString]["population_generation_rate_multiplier"].(float64)
		p.MaxPopulationGenerationRate.Next = models.MaxPopulationGenerationRate(populationControlCenterBuildingConstants[maxLevelString])
		p.MinimumWorkersRequired.Next = populationControlCenterBuildingConstants[maxLevelString]["workers_required"].(float64)
		p.WorkersMaxLimit.Next = populationControlCenterBuildingConstants[maxLevelString]["workers_max_limit"].(float64)
	}
}

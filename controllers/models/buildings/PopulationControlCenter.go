package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PopulationControlCenter struct {
	BuildingId            string                            `json:"building_id,omitempty" example:"WMP101"`
	Level                 int                               `json:"level" example:"3"`
	Workers               int                               `json:"workers" example:"12"`
	BuildingState         BuildingState                     `json:"building_state"`
	BuildingAttributes    PopulationControlCenterAttributes `json:"building_attributes"`
	NextLevelRequirements NextLevelRequirements             `json:"next_level_requirements"`
}
type PopulationControlCenterAttributes struct {
	MaxPopulationGenerationRate        IntegerBuildingAttributes `json:"max_population_generation_rate"`
	PopulationGenerationRateMultiplier FloatBuildingAttributes   `json:"population_generation_rate_multiplier"`
	MinimumWorkersRequired             IntegerBuildingAttributes `json:"minimum_workers_required"`
	WorkersMaxLimit                    IntegerBuildingAttributes `json:"workers_max_limit" `
}

func InitPopulationControlCenter(planetUser repoModels.PlanetUser,
	populationControlCenterUpgradeConstants constants.UpgradeConstants,
	populationControlCenterBuildingConstants constants.BuildingConstants) *PopulationControlCenter {

	p := new(PopulationControlCenter)
	p.BuildingId = constants.PopulationControlCenter
	p.Level = planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel
	p.Workers = planetUser.Buildings[constants.PopulationControlCenter].Workers
	p.BuildingState.Init(planetUser.Buildings[constants.PopulationControlCenter], populationControlCenterUpgradeConstants)
	p.NextLevelRequirements.Init(planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel, populationControlCenterUpgradeConstants)
	p.BuildingAttributes.Init(planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel,
		populationControlCenterUpgradeConstants.MaxLevel, populationControlCenterBuildingConstants)
	return p
}

func (p *PopulationControlCenterAttributes) Init(currentLevel int, maxLevel int,
	populationControlCenterBuildingConstants constants.BuildingConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(maxLevel)

	p.MaxPopulationGenerationRate.Current, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["max_population_generation_rate"])
	p.PopulationGenerationRateMultiplier.Current, _ =
		strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[currentLevelString]["population_generation_rate_multiplier"], 64)
	p.MinimumWorkersRequired.Current, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["workers_required"])
	p.WorkersMaxLimit.Current, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["workers_max_limit"])

	workersMaxLimit, _ := strconv.Atoi(populationControlCenterBuildingConstants.Levels[maxLevelString]["workers_max_limit"])
	p.MaxPopulationGenerationRate.Max = models.GetMaxPopulationGenerationRate(populationControlCenterBuildingConstants.Levels[maxLevelString], workersMaxLimit)
	p.WorkersMaxLimit.Max, _ = strconv.Atoi(populationControlCenterBuildingConstants.Levels[maxLevelString]["workers_max_limit"])

	if currentLevel+1 < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		p.MaxPopulationGenerationRate.Next, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["max_population_generation_rate"])
		p.PopulationGenerationRateMultiplier.Next, _ =
			strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[nextLevelString]["population_generation_rate_multiplier"], 64)
		p.MinimumWorkersRequired.Next, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["workers_required"])
		p.WorkersMaxLimit.Next, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["workers_max_limit"])
	}
}

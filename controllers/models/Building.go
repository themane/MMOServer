package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PopulationControlCenter struct {
	BuildingId            string                                     `json:"building_id,omitempty" example:"WMP101"`
	Level                 int                                        `json:"level" example:"3"`
	Workers               int                                        `json:"workers" example:"12"`
	BuildingState         BuildingState                              `json:"building_state"`
	NextLevelAttributes   NextLevelPopulationControlCenterAttributes `json:"next_level_attributes"`
	NextLevelRequirements NextLevelRequirements                      `json:"next_level_requirements"`
}
type NextLevelPopulationControlCenterAttributes struct {
	CurrentMaxPopulationGenerationRate           int     `json:"current_max_population_generation_rate" example:"2"`
	NextMaxPopulationGenerationRate              int     `json:"next_max_population_generation_rate" example:"4"`
	MaxPopulationGenerationRate                  int     `json:"max_population_generation_rate" example:"16"`
	CurrentMaxPopulationGenerationRateMultiplier float64 `json:"current_max_population_generation_rate_multiplier" example:"0.1"`
	NextMaxPopulationGenerationRateMultiplier    float64 `json:"next_max_population_generation_rate_multiplier" example:"0.2"`
	MaxPopulationGenerationRateMultiplier        float64 `json:"max_population_generation_rate_multiplier" example:"0.5"`
	CurrentMinimumWorkersRequired                int     `json:"current_minimum_workers_required" example:"10"`
	NextMinimumWorkersRequired                   int     `json:"next_minimum_workers_required" example:"15"`
	CurrentMinimumSoldiersRequired               int     `json:"current_minimum_soldiers_required" example:"8"`
	NextMinimumSoldiersRequired                  int     `json:"next_minimum_soldiers_required" example:"12"`
}

func (p *PopulationControlCenter) InitPopulationControlCenter(planetUser models.PlanetUser,
	populationControlCenterUpgradeConstants constants.UpgradeConstants,
	populationControlCenterBuildingConstants constants.BuildingConstants) {

	p.BuildingId = constants.PopulationControlCenter
	p.Level = planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel
	p.Workers = planetUser.Buildings[constants.PopulationControlCenter].Workers
	p.BuildingState.Init(planetUser.Buildings[constants.PopulationControlCenter], populationControlCenterUpgradeConstants)
	p.NextLevelRequirements.Init(planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel, populationControlCenterUpgradeConstants)
	p.NextLevelAttributes.Init(planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel,
		populationControlCenterUpgradeConstants.MaxLevel, populationControlCenterBuildingConstants)
}

func (p *NextLevelPopulationControlCenterAttributes) Init(currentLevel int, maxLevel int,
	populationControlCenterBuildingConstants constants.BuildingConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(maxLevel)
	p.CurrentMaxPopulationGenerationRate, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["max_population_generation_rate"])
	p.CurrentMaxPopulationGenerationRateMultiplier, _ =
		strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[currentLevelString]["max_population_generation_rate_multiplier"], 64)
	p.CurrentMinimumWorkersRequired, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["workers_required"])
	p.CurrentMinimumSoldiersRequired, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[currentLevelString]["soldiers_required"])
	p.MaxPopulationGenerationRate, _ =
		strconv.Atoi(populationControlCenterBuildingConstants.Levels[maxLevelString]["max_population_generation_rate"])
	p.MaxPopulationGenerationRateMultiplier, _ =
		strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[maxLevelString]["max_population_generation_rate_multiplier"], 64)
	if currentLevel+1 < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		p.NextMaxPopulationGenerationRate, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["max_population_generation_rate"])
		p.NextMaxPopulationGenerationRateMultiplier, _ =
			strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[nextLevelString]["max_population_generation_rate_multiplier"], 64)
		p.NextMaxPopulationGenerationRateMultiplier, _ =
			strconv.ParseFloat(populationControlCenterBuildingConstants.Levels[nextLevelString]["max_population_generation_rate_multiplier"], 64)
		p.NextMinimumWorkersRequired, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["workers_required"])
		p.NextMinimumSoldiersRequired, _ =
			strconv.Atoi(populationControlCenterBuildingConstants.Levels[nextLevelString]["soldiers_required"])
	}
}

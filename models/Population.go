package models

import (
	"math"
)

func MaxPopulationGenerationRate(populationControlCenterConstants map[string]interface{}) float64 {
	maxPopulationGenerationRate := populationControlCenterConstants["max_population_generation_rate"].(float64)
	maxPopulationGenerationRateMultiplier := populationControlCenterConstants["population_generation_rate_multiplier"].(float64)
	maxSelectableWorkers := populationControlCenterConstants["workers_max_limit"].(float64)
	return math.Floor(maxPopulationGenerationRate + maxPopulationGenerationRateMultiplier*maxSelectableWorkers)
}

func MaxSelectablePopulationGenerationRate(
	populationControlCenterConstants map[string]interface{}, currentWorkers int) float64 {
	maxPopulationGenerationRate := populationControlCenterConstants["max_population_generation_rate"].(float64)
	maxPopulationGenerationRateMultiplier := populationControlCenterConstants["population_generation_rate_multiplier"].(float64)
	return math.Floor(maxPopulationGenerationRate + maxPopulationGenerationRateMultiplier*float64(currentWorkers))
}

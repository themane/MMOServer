package models

import (
	"math"
)

func GetMaxPopulationGenerationRate(
	populationControlCenterConstants map[string]interface{}, workers float64) float64 {
	maxPopulationGenerationRate := populationControlCenterConstants["max_population_generation_rate"].(float64)
	maxPopulationGenerationRateMultiplier := populationControlCenterConstants["population_generation_rate_multiplier"].(float64)
	return math.Floor(maxPopulationGenerationRate + maxPopulationGenerationRateMultiplier*float64(workers))
}

package models

import (
	"math"
	"strconv"
)

func GetMaxPopulationGenerationRate(
	populationControlCenterConstants map[string]string, workers int) int {
	maxPopulationGenerationRate, _ := strconv.Atoi(populationControlCenterConstants["max_population_generation_rate"])
	maxPopulationGenerationRateMultiplier, _ := strconv.ParseFloat(populationControlCenterConstants["population_generation_rate_multiplier"], 64)
	return int(math.Floor(float64(maxPopulationGenerationRate) + maxPopulationGenerationRateMultiplier*float64(workers)))
}

package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PlanetService struct {
	userRepository    repoModels.UserRepository
	upgradeConstants  map[string]constants.UpgradeConstants
	buildingConstants map[string]constants.BuildingConstants
	logger            *constants.LoggingUtils
}

func NewPlanetService(
	userRepository repoModels.UserRepository,
	buildingConstants map[string]constants.BuildingConstants,
	logLevel string,
) *PlanetService {
	return &PlanetService{
		userRepository:    userRepository,
		buildingConstants: buildingConstants,
		logger:            constants.NewLoggingUtils("PLANET_SERVICE", logLevel),
	}
}

func (p *PlanetService) UpdatePopulationRate(username string, planetId string, generationRate int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	currentDeployedWorkers := userData.OccupiedPlanets[planetId].Buildings[constants.PopulationControlCenter].Workers
	populationControlCenterLevel := userData.OccupiedPlanets[planetId].Buildings[constants.PopulationControlCenter].BuildingLevel

	populationControlCenterConstants := p.buildingConstants[constants.PopulationControlCenter].Levels[strconv.Itoa(populationControlCenterLevel)]
	maxPopulationGenerationRate, err := strconv.Atoi(populationControlCenterConstants["max_population_generation_rate"])
	if err != nil {
		return err
	}
	populationGenerationRateMultiplier, err := strconv.Atoi(populationControlCenterConstants["population_generation_rate_multiplier"])
	if err != nil {
		return err
	}
	finalMaxPopulationGenerationRate := maxPopulationGenerationRate + (populationGenerationRateMultiplier * currentDeployedWorkers)
	if finalMaxPopulationGenerationRate < generationRate {
		return errors.New("rate above maximum")
	}

	err = p.userRepository.UpdatePopulationRate(userData.Id, planetId, generationRate)
	if err != nil {
		return err
	}
	return nil
}

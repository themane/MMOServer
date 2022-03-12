package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PlanetService struct {
	userRepository    repoModels.UserRepository
	upgradeConstants  map[string]constants.UpgradeConstants
	buildingConstants map[string]map[string]map[string]interface{}
	logger            *constants.LoggingUtils
}

func NewPlanetService(
	userRepository repoModels.UserRepository,
	buildingConstants map[string]map[string]map[string]interface{},
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

	populationControlCenterConstants := p.buildingConstants[constants.PopulationControlCenter][strconv.Itoa(populationControlCenterLevel)]
	maxPopulationGenerationRate := int(models.GetMaxPopulationGenerationRate(populationControlCenterConstants, float64(currentDeployedWorkers)))
	if maxPopulationGenerationRate < generationRate {
		return errors.New("rate above maximum")
	}

	err = p.userRepository.UpdatePopulationRate(userData.Id, planetId, generationRate)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) EmployPopulation(username string, planetId string, workers int, soldiers int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	unemployedPopulation := userData.OccupiedPlanets[planetId].Population.Unemployed
	if workers+soldiers > unemployedPopulation {
		return errors.New("not enough population to employ")
	}

	err = p.userRepository.Recruit(userData.Id, planetId, workers, soldiers)
	if err != nil {
		return err
	}
	return nil
}

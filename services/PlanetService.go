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
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		currentDeployedWorkers := planetUser.Buildings[constants.PopulationControlCenter].Workers
		populationControlCenterLevel := planetUser.Buildings[constants.PopulationControlCenter].BuildingLevel

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
	return errors.New("planet not occupied")
}

func (p *PlanetService) EmployPopulation(username string, planetId string, workers int, soldiers int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		unemployedPopulation := planetUser.Population.Unemployed
		if workers+soldiers > unemployedPopulation {
			return errors.New("not enough population to employ")
		}

		err = p.userRepository.Recruit(userData.Id, planetId, workers, soldiers)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("planet not occupied")
}

func (p *PlanetService) KillPopulation(username string, planetId string, unemployed int, workers int, soldiers int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if unemployed > planetUser.Population.Unemployed ||
			workers > planetUser.Population.IdleWorkers ||
			soldiers > planetUser.Population.IdleSoldiers {
			return errors.New("not enough population to kill")
		}

		err = p.userRepository.KillPopulation(userData.Id, planetId, unemployed, workers, soldiers)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("planet not occupied")
}

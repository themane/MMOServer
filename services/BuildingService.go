package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type BuildingService struct {
	userRepository   repoModels.UserRepository
	upgradeConstants map[string]constants.UpgradeConstants
	logger           *constants.LoggingUtils
}

func NewBuildingService(
	userRepository repoModels.UserRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	logLevel string,
) *BuildingService {
	return &BuildingService{
		userRepository:   userRepository,
		upgradeConstants: upgradeConstants,
		logger:           constants.NewLoggingUtils("BUILDING_SERVICE", logLevel),
	}
}

func (b *BuildingService) UpgradeBuilding(username string, planetId string, buildingId string) error {
	userData, err := b.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	requirements, err := b.verifyAndGetRequiredResources(*userData, planetId, buildingId)
	if err != nil {
		return err
	}
	err = b.userRepository.UpgradeBuildingLevel(userData.Id, planetId, buildingId,
		requirements.WaterRequired, requirements.GrapheneRequired, requirements.ShelioRequired, requirements.MinutesRequiredPerWorker)
	if err != nil {
		return err
	}
	return nil
}

func (b *BuildingService) CancelUpgradeBuilding(username string, planetId string, buildingId string) error {
	userData, err := b.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	returns, err := b.verifyAndGetReturnedResources(*userData, planetId, buildingId)
	if err != nil {
		return err
	}
	err = b.userRepository.CancelUpgradeBuildingLevel(userData.Id, planetId, buildingId,
		returns.WaterReturned, returns.GrapheneReturned, returns.ShelioReturned)
	if err != nil {
		return err
	}
	return nil
}

func (b *BuildingService) UpdateWorkers(username string, planetId string, buildingId string, workers int) error {
	userData, err := b.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if constants.IsShieldId(buildingId) && userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingMinutesPerWorker == 0 {
		return errors.New("workers not employed at working shield")
	}
	currentWorkers := userData.OccupiedPlanets[planetId].Buildings[buildingId].Workers
	idleWorkers := userData.OccupiedPlanets[planetId].Population.IdleWorkers
	if currentWorkers == workers {
		return nil
	} else if workers > currentWorkers && idleWorkers < workers-currentWorkers {
		return errors.New("not enough workers")
	}
	return b.userRepository.UpdateWorkers(userData.Id, planetId, buildingId, workers-currentWorkers)
}

func (b *BuildingService) UpdateSoldiers(username string, planetId string, buildingId string, soldiers int) error {
	userData, err := b.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if _, ok := constants.GetSoldiersSupportedBuildingIds()[buildingId]; !ok {
		return errors.New("soldiers not employed at " + buildingId)
	}
	currentSoldiers := userData.OccupiedPlanets[planetId].Buildings[buildingId].Soldiers
	idleSoldiers := userData.OccupiedPlanets[planetId].Population.IdleSoldiers
	if currentSoldiers == soldiers {
		return nil
	} else if soldiers > currentSoldiers && idleSoldiers < soldiers-currentSoldiers {
		return errors.New("not enough soldiers")
	}
	return b.userRepository.UpdateSoldiers(userData.Id, planetId, buildingId, soldiers-currentSoldiers)
}

func (b *BuildingService) verifyAndGetRequiredResources(userData repoModels.UserData,
	planetId string, buildingId string) (*repoModels.NextLevelRequirements, error) {

	if userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingMinutesPerWorker > 0 {
		return nil, errors.New("building already under upgradation")
	}
	buildingLevel := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingLevel
	buildingType, err := constants.GetBuildingType(buildingId)
	if err != nil {
		return nil, errors.New("building not found")
	}

	if buildingConstants, ok := b.upgradeConstants[buildingType]; ok {
		if buildingConstants.MaxLevel <= buildingLevel {
			return nil, errors.New("max level reached")
		}
		requirements := repoModels.NextLevelRequirements{}
		requirements.Init(buildingLevel, buildingConstants)
		if userData.OccupiedPlanets[planetId].Water.Amount < requirements.WaterRequired ||
			userData.OccupiedPlanets[planetId].Graphene.Amount < requirements.GrapheneRequired ||
			userData.OccupiedPlanets[planetId].Shelio < requirements.ShelioRequired {
			return nil, errors.New("not enough resources")
		}
		return &requirements, nil
	}
	return nil, errors.New("building type not found")
}

func (b *BuildingService) verifyAndGetReturnedResources(userData repoModels.UserData,
	planetId string, buildingId string) (*repoModels.CancelReturns, error) {

	buildingMinutesPerWorker := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingMinutesPerWorker
	if buildingMinutesPerWorker == 0 {
		return nil, errors.New("building not under upgradation")
	}
	buildingLevel := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingLevel
	buildingType, err := constants.GetBuildingType(buildingId)
	if err != nil {
		return nil, errors.New("building not found")
	}
	if buildingConstants, ok := b.upgradeConstants[buildingType]; ok {
		if buildingConstants.MaxLevel <= buildingLevel {
			return nil, errors.New("max level reached")
		}
		returns := repoModels.CancelReturns{}
		returns.Init(buildingMinutesPerWorker, buildingLevel, buildingConstants)
		return &returns, nil
	}
	return nil, errors.New("building type not found")
}

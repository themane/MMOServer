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
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	building := planetUser.GetBuilding(buildingId)
	if building == nil {
		return errors.New("building id not valid")
	}
	if constants.IsShieldId(buildingId) && building.BuildingMinutesPerWorker == 0 {
		return errors.New("workers not employed at working shield")
	}
	currentWorkers := building.Workers
	idleWorkers := planetUser.Population.IdleWorkers
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
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	building := planetUser.GetBuilding(buildingId)
	if building == nil {
		return errors.New("building id not valid")
	}
	currentSoldiers := building.Soldiers
	idleSoldiers := planetUser.Population.IdleSoldiers
	if currentSoldiers == soldiers {
		return nil
	} else if soldiers > currentSoldiers && idleSoldiers < soldiers-currentSoldiers {
		return errors.New("not enough soldiers")
	}
	return b.userRepository.UpdateSoldiers(userData.Id, planetId, buildingId, soldiers-currentSoldiers)
}

func (b *BuildingService) verifyAndGetRequiredResources(userData repoModels.UserData,
	planetId string, buildingId string) (*repoModels.NextLevelRequirements, error) {

	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return nil, errors.New("planet not occupied")
	}
	building := planetUser.GetBuilding(buildingId)
	if building == nil {
		return nil, errors.New("building id not valid")
	}
	if building.BuildingMinutesPerWorker > 0 {
		return nil, errors.New("building already under upgradation")
	}
	buildingLevel := building.BuildingLevel
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
		if planetUser.Water.Amount < requirements.WaterRequired ||
			planetUser.Graphene.Amount < requirements.GrapheneRequired ||
			planetUser.Shelio < requirements.ShelioRequired {
			return nil, errors.New("not enough resources")
		}
		return &requirements, nil
	}
	return nil, errors.New("building type not found")
}

func (b *BuildingService) verifyAndGetReturnedResources(userData repoModels.UserData,
	planetId string, buildingId string) (*repoModels.CancelReturns, error) {

	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return nil, errors.New("planet not occupied")
	}
	building := planetUser.GetBuilding(buildingId)
	if building == nil {
		return nil, errors.New("building id not valid")
	}
	buildingMinutesPerWorker := building.BuildingMinutesPerWorker
	if buildingMinutesPerWorker == 0 {
		return nil, errors.New("building not under upgradation")
	}
	buildingLevel := building.BuildingLevel
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

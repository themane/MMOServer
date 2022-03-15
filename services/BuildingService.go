package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
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
	waterRequired, grapheneRequired, shelioRequired, minutesRequired, err := b.verifyAndGetRequiredResources(*userData, planetId, buildingId)
	if err != nil {
		return err
	}
	err = b.userRepository.UpgradeBuildingLevel(userData.Id, planetId, buildingId, waterRequired, grapheneRequired, shelioRequired, minutesRequired)
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
	planetId string, buildingId string) (int, int, int, int, error) {

	if userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingMinutesPerWorker > 0 {
		return 0, 0, 0, 0, errors.New("building already under upgradation")
	}
	buildingLevel := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingLevel
	nextBuildingLevelString := strconv.Itoa(buildingLevel + 1)
	buildingType, err := constants.GetBuildingType(buildingId)
	if err != nil {
		return 0, 0, 0, 0, errors.New("building not found")
	}
	if buildingConstants, ok := b.upgradeConstants[buildingType]; ok {
		if buildingConstants.MaxLevel <= buildingLevel {
			return 0, 0, 0, 0, errors.New("max level reached")
		}
		waterRequired := buildingConstants.Levels[nextBuildingLevelString].WaterRequired
		grapheneRequired := buildingConstants.Levels[nextBuildingLevelString].GrapheneRequired
		shelioRequired := buildingConstants.Levels[nextBuildingLevelString].ShelioRequired
		minutesRequired := buildingConstants.Levels[nextBuildingLevelString].MinutesRequired
		if userData.OccupiedPlanets[planetId].Water.Amount < waterRequired ||
			userData.OccupiedPlanets[planetId].Graphene.Amount < grapheneRequired ||
			userData.OccupiedPlanets[planetId].Shelio < shelioRequired {
			return 0, 0, 0, 0, errors.New("not enough resources")
		}
		return waterRequired, grapheneRequired, shelioRequired, minutesRequired, nil
	}
	return 0, 0, 0, 0, errors.New("building not found")
}

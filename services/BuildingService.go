package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type BuildingService struct {
	userRepository    repoModels.UserRepository
	buildingConstants map[string]constants.BuildingConstants
}

func NewBuildingService(
	userRepository repoModels.UserRepository,
	buildingConstants map[string]constants.BuildingConstants,
) *BuildingService {
	return &BuildingService{
		userRepository:    userRepository,
		buildingConstants: buildingConstants,
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

func (b *BuildingService) verifyAndGetRequiredResources(userData repoModels.UserData,
	planetId string, buildingId string) (int, int, int, int, error) {

	if userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingMinutesPerWorker > 0 {
		return 0, 0, 0, 0, errors.New("building already under upgradation")
	}
	buildingLevel := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingLevel
	nextBuildingLevelString := strconv.Itoa(buildingLevel + 1)
	buildingType := constants.GetBuildingType(buildingId)
	if buildingConstants, ok := b.buildingConstants[buildingType]; ok {
		if buildingConstants.MaxLevel <= buildingLevel {
			return 0, 0, 0, 0, errors.New("max level reached")
		}
		waterRequired := buildingConstants.Levels[nextBuildingLevelString].WaterRequired
		grapheneRequired := buildingConstants.Levels[nextBuildingLevelString].GrapheneRequired
		shelioRequired := buildingConstants.Levels[nextBuildingLevelString].ShelioRequired
		minutesRequired := buildingConstants.Levels[nextBuildingLevelString].MinutesRequired
		if userData.OccupiedPlanets[planetId].Water.Amount >= waterRequired &&
			userData.OccupiedPlanets[planetId].Graphene.Amount >= grapheneRequired &&
			userData.OccupiedPlanets[planetId].Shelio >= shelioRequired {
			return 0, 0, 0, 0, errors.New("not enough resources")
		}
		return waterRequired, grapheneRequired, shelioRequired, minutesRequired, nil
	}
	return 0, 0, 0, 0, errors.New("building not found")
}

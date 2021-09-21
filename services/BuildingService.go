package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
	"strings"
)

const (
	WATER_MINE    = "WATER_MINE"
	GRAPHENE_MINE = "GRAPHENE_MINE"
)

func UpgradeBuilding(username string, planetId string, buildingId string, userRepository repoModels.UserRepository) error {
	waterConstants := constants.GetWaterConstants()
	grapheneConstants := constants.GetGrapheneConstants()

	userData, err := userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	waterRequired, grapheneRequired, shelioRequired, err := verifyAndGetRequiredResources(*userData, planetId, buildingId, waterConstants, grapheneConstants)
	if err != nil {
		return err
	}
	err = userRepository.UpgradeBuildingLevel(userData.Id, planetId, buildingId, waterRequired, grapheneRequired, shelioRequired)
	if err != nil {
		return err
	}
	return nil
}

func verifyAndGetRequiredResources(userData repoModels.UserData, planetId string, buildingId string,
	waterConstants constants.ResourceConstants, grapheneConstants constants.ResourceConstants) (int, int, int, error) {

	buildingLevel := userData.OccupiedPlanets[planetId].Buildings[buildingId].BuildingLevel
	nextBuildingLevelString := strconv.Itoa(buildingLevel + 1)
	switch getBuildingType(buildingId) {
	case WATER_MINE:
		if waterConstants.MaxLevel <= buildingLevel {
			return 0, 0, 0, errors.New("max level reached")
		}
		waterRequired := waterConstants.Levels[nextBuildingLevelString].WaterRequired
		grapheneRequired := waterConstants.Levels[nextBuildingLevelString].GrapheneRequired
		shelioRequired := waterConstants.Levels[nextBuildingLevelString].ShelioRequired
		if userData.OccupiedPlanets[planetId].Water.Amount >= waterRequired &&
			userData.OccupiedPlanets[planetId].Graphene.Amount >= grapheneRequired &&
			userData.OccupiedPlanets[planetId].Shelio >= shelioRequired {
			return 0, 0, 0, errors.New("not enough resources")
		}
		return waterRequired, grapheneRequired, shelioRequired, nil
	case GRAPHENE_MINE:
		if grapheneConstants.MaxLevel <= buildingLevel {
			return 0, 0, 0, errors.New("max level reached")
		}
		waterRequired := grapheneConstants.Levels[nextBuildingLevelString].WaterRequired
		grapheneRequired := grapheneConstants.Levels[nextBuildingLevelString].GrapheneRequired
		shelioRequired := grapheneConstants.Levels[nextBuildingLevelString].ShelioRequired
		if userData.OccupiedPlanets[planetId].Water.Amount >= waterRequired &&
			userData.OccupiedPlanets[planetId].Graphene.Amount >= grapheneRequired &&
			userData.OccupiedPlanets[planetId].Shelio >= shelioRequired {
			return waterRequired, grapheneRequired, shelioRequired, nil
		}
		return 0, 0, 0, errors.New("not enough resources")
	}
	return 0, 0, 0, errors.New("building not found")
}

func getBuildingType(buildingId string) string {
	if strings.HasPrefix(buildingId, "WMP") {
		return WATER_MINE
	}
	if strings.HasPrefix(buildingId, "GMP") {
		return GRAPHENE_MINE
	}
	return ""
}

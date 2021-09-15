package services

import (
	"github.com/themane/MMOServer/models"
	"strconv"
)

func upgradeMiningPlant(planetUser models.PlanetUser, buildingId string, waterConstants models.ResourceConstants, grapheneConstants models.ResourceConstants) (string, string) {
	for _, mine := range planetUser.Mines {
		if mine.MiningPlant.BuildingId == buildingId {
			buildingLevelString := strconv.Itoa(mine.MiningPlant.BuildingLevel)
			if mine.Type == models.WATER {
				waterRequired := waterConstants.Levels[buildingLevelString].WaterRequired
				grapheneRequired := waterConstants.Levels[buildingLevelString].GrapheneRequired
				shelioRequired := waterConstants.Levels[buildingLevelString].ShelioRequired
				return upgradeMine(planetUser, buildingId, waterConstants.MaxLevel, mine.MiningPlant, waterRequired, grapheneRequired, shelioRequired)
			}
			if mine.Type == models.GRAPHENE {
				waterRequired := grapheneConstants.Levels[buildingLevelString].WaterRequired
				grapheneRequired := grapheneConstants.Levels[buildingLevelString].GrapheneRequired
				shelioRequired := grapheneConstants.Levels[buildingLevelString].ShelioRequired
				return upgradeMine(planetUser, buildingId, grapheneConstants.MaxLevel, mine.MiningPlant, waterRequired, grapheneRequired, shelioRequired)
			}
		}
	}
	return "", "Invalid building_id"
}

func upgradeMine(planetUser models.PlanetUser, buildingId string, maxLevel int, miningPlant models.MiningPlantUser, waterRequired int, grapheneRequired int, shelioRequired int) (string, string) {
	if miningPlant.BuildingLevel >= maxLevel {
		return "", "Max level Reached"
	}
	if waterRequired <= planetUser.Water.Amount && grapheneRequired <= planetUser.Graphene.Amount && shelioRequired <= planetUser.Shelio {
		miningPlant.BuildingLevel += 1
		return "Successfully updated: " + buildingId, ""
	}
	return "", "Insufficient resources"
}

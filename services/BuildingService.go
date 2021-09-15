package services

import (
	"github.com/themane/MMOServer/dao"
	"strings"
)

const (
	MINE = "MINE"
)

func UpgradeBuilding(username string, planetId string, buildingId string) (string, string) {
	waterConstants := dao.GetWaterConstants()
	grapheneConstants := dao.GetGrapheneConstants()
	userData := dao.GetUserData(username)
	buildingType := getBuildingType(buildingId)
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Position.PlanetId() == planetId {
			switch buildingType {
			case MINE:
				upgradeMiningPlant(planetUser, buildingId, waterConstants, grapheneConstants)
			}
		}
	}
	return "", "Invalid username or planet_id or building_id"
}

func getBuildingType(buildingId string) string {
	if strings.HasPrefix(buildingId, "WMP") || strings.HasPrefix(buildingId, "GMP") {
		return MINE
	}
	return ""
}

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
	for key, planetUser := range userData.OccupiedPlanets {
		if planetUser.Position.PlanetId() == planetId {
			switch buildingType {
			case MINE:
				msg, err := upgradeMiningPlant(&planetUser, buildingId, waterConstants, grapheneConstants)
				userData.OccupiedPlanets[key] = planetUser
				dao.UpdateUserData(username, userData)
				return msg, err
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

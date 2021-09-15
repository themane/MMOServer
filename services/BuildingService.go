package services

import (
	"github.com/themane/MMOServer/dao"
	"github.com/themane/MMOServer/models"
)

func UpgradeBuilding(username string, planetId string, buildingId string) *models.Population {
	userData := dao.GetUserData(username)
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Position.PlanetId() == planetId {
			response := models.Population{}
			response.Init(planetUser)
			return &response
		}
	}
	return nil
}

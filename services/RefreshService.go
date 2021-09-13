package services

import (
	"github.com/themane/MMOServer/dao"
	"github.com/themane/MMOServer/models"
)

func RefreshPopulation(username string, planetId string) *models.Population {
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

func RefreshResources(username string, planetId string) *models.Resources {
	userData := dao.GetUserData(username)
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Position.PlanetId() == planetId {
			response := models.Resources{}
			response.Init(planetUser)
			return &response
		}
	}
	return nil
}

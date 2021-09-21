package services

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

func RefreshPopulation(username string, inputPlanetId string, userRepository repoModels.UserRepository) (*models.Population, error) {
	userData, err := userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			response := models.Population{}
			response.Init(planetUser)
			return &response, nil
		}
	}
	return nil, nil
}

func RefreshResources(username string, inputPlanetId string, userRepository repoModels.UserRepository) (*models.Resources, error) {
	userData, err := userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			response := models.Resources{}
			response.Init(planetUser)
			return &response, nil
		}
	}
	return nil, nil
}

func RefreshMine(username string, inputPlanetId string, inputMineId string,
	userRepository repoModels.UserRepository, universeRepository repoModels.UniverseRepository,
	waterConstants constants.ResourceConstants, grapheneConstants constants.ResourceConstants) (*models.Mine, error) {
	userData, errUser := userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			planetUni, errUni := universeRepository.FindById(planetId)
			if errUni != nil {
				return nil, errUni
			}
			for mineId, mineUni := range planetUni.Mines {
				if mineId == inputMineId {
					response := models.Mine{}
					response.Init(mineUni, planetUser, waterConstants, grapheneConstants)
					return &response, nil
				}
			}
		}
	}
	return nil, nil
}

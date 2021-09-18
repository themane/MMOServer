package services

import (
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

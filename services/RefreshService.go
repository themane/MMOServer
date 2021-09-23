package services

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type QuickRefreshService struct {
	userRepository     repoModels.UserRepository
	universeRepository repoModels.UniverseRepository
	buildingConstants  map[string]constants.BuildingConstants
	waterConstants     constants.MiningConstants
	grapheneConstants  constants.MiningConstants
	defenceConstants   map[string]constants.DefenceConstants
}

func NewQuickRefreshService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
) *QuickRefreshService {
	return &QuickRefreshService{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		buildingConstants:  buildingConstants,
		waterConstants:     mineConstants[constants.Water],
		grapheneConstants:  mineConstants[constants.Graphene],
		defenceConstants:   defenceConstants,
	}
}

func (r *QuickRefreshService) RefreshPopulation(username string, inputPlanetId string) (*models.Population, error) {
	userData, err := r.userRepository.FindByUsername(username)
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

func (r *QuickRefreshService) RefreshResources(username string, inputPlanetId string) (*models.Resources, error) {
	userData, err := r.userRepository.FindByUsername(username)
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

func (r *QuickRefreshService) RefreshMine(username string, inputPlanetId string, inputMineId string) (*models.Mine, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			planetUni, errUni := r.universeRepository.FindById(planetId)
			if errUni != nil {
				return nil, errUni
			}
			for mineId, mineUni := range planetUni.Mines {
				if mineId == inputMineId {
					response := models.Mine{}
					response.Init(mineUni, planetUser,
						r.buildingConstants[constants.WaterMiningPlant], r.buildingConstants[constants.GrapheneMiningPlant],
						r.waterConstants, r.grapheneConstants)
					return &response, nil
				}
			}
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshShields(username string, inputPlanetId string) ([]models.Shield, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			return models.InitAllShields(planetUser, r.defenceConstants[constants.Shield], r.buildingConstants[constants.Shield]), nil
		}
	}
	return nil, nil
}

package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type QuickRefreshService struct {
	userRepository     repoModels.UserRepository
	universeRepository repoModels.UniverseRepository
	missionRepository  repoModels.MissionRepository
	buildingConstants  map[string]constants.BuildingConstants
	waterConstants     constants.MiningConstants
	grapheneConstants  constants.MiningConstants
	defenceConstants   map[string]constants.DefenceConstants
	logger             *constants.LoggingUtils
}

func NewQuickRefreshService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	logLevel string,
) *QuickRefreshService {
	return &QuickRefreshService{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		missionRepository:  missionRepository,
		buildingConstants:  buildingConstants,
		waterConstants:     mineConstants[constants.Water],
		grapheneConstants:  mineConstants[constants.Graphene],
		defenceConstants:   defenceConstants,
		logger:             constants.NewLoggingUtils("REFRESH_SERVICE", logLevel),
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
			return models.InitAllShields(planetUser, r.defenceConstants, r.buildingConstants[constants.Shield]), nil
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshMissions(username string, inputPlanetId string) (map[string][]models.ActiveMission, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			attackMissions, err := r.missionRepository.FindAttackMissionsFromPlanetId(planetId)
			if err != nil {
				r.logger.Error("error in retrieving attack missions for: "+planetId, err)
				return nil, errors.New("error in retrieving attack missions")
			}
			spyMissions, err := r.missionRepository.FindSpyMissionsFromPlanetId(planetId)
			if err != nil {
				r.logger.Error("error in retrieving spy missions for: "+planetId, err)
				return nil, errors.New("error in retrieving spy missions")
			}
			activeMissions := map[string][]models.ActiveMission{}
			activeMissions["attack_missions"] = []models.ActiveMission{}
			activeMissions["spy_missions"] = []models.ActiveMission{}
			for _, attackMission := range attackMissions {
				activeMission := models.ActiveMission{}
				activeMission.InitAttackMission(attackMission)
				activeMissions["attack_missions"] = append(activeMissions["attack_missions"], activeMission)
			}
			for _, spyMission := range spyMissions {
				activeMission := models.ActiveMission{}
				activeMission.InitSpyMission(spyMission)
				activeMissions["spy_missions"] = append(activeMissions["spy_missions"], activeMission)
			}
			return activeMissions, nil
		}
	}
	return nil, nil
}

package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
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
	shipConstants      map[string]constants.ShipConstants
	logger             *constants.LoggingUtils
}

func NewQuickRefreshService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
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
		shipConstants:      shipConstants,
		logger:             constants.NewLoggingUtils("REFRESH_SERVICE", logLevel),
	}
}

func (r *QuickRefreshService) RefreshPlanet(username string, inputPlanetId string) (*controllerModels.OccupiedPlanet, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			planetUni, err := r.universeRepository.FindById(planetId)
			if err != nil {
				return nil, err
			}
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
			planetResponse := controllerModels.OccupiedPlanet{}
			planetResponse.Init(*planetUni, planetUser, inputPlanetId, attackMissions, spyMissions,
				r.buildingConstants, r.waterConstants, r.grapheneConstants, r.defenceConstants, r.shipConstants)
			return &planetResponse, nil
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshUserPlanet(username string, inputPlanetId string) (*controllerModels.UserPlanetResponse, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId, planetUser := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			response := controllerModels.UserPlanetResponse{}
			var notifications []models.Notification
			response.Init(planetUser, r.buildingConstants, r.defenceConstants, r.shipConstants, notifications)
			return &response, nil
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshMine(username string, inputPlanetId string, inputMineId string) (*controllerModels.Mine, error) {
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
					response := controllerModels.Mine{}
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

func (r *QuickRefreshService) RefreshAttackMissions(username string, inputPlanetId string) ([]controllerModels.ActiveMission, error) {
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
			var activeMissions []controllerModels.ActiveMission
			for _, attackMission := range attackMissions {
				activeMission := controllerModels.ActiveMission{}
				activeMission.InitAttackMission(attackMission)
				activeMissions = append(activeMissions, activeMission)
			}
			return activeMissions, nil
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshSpyMissions(username string, inputPlanetId string) ([]controllerModels.ActiveMission, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for planetId := range userData.OccupiedPlanets {
		if planetId == inputPlanetId {
			spyMissions, err := r.missionRepository.FindSpyMissionsFromPlanetId(planetId)
			if err != nil {
				r.logger.Error("error in retrieving spy missions for: "+planetId, err)
				return nil, errors.New("error in retrieving spy missions")
			}
			var activeMissions []controllerModels.ActiveMission
			for _, spyMission := range spyMissions {
				activeMission := controllerModels.ActiveMission{}
				activeMission.InitSpyMission(spyMission)
				activeMissions = append(activeMissions, activeMission)
			}
			return activeMissions, nil
		}
	}
	return nil, nil
}

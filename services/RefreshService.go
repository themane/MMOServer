package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/controllers/models/buildings"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type QuickRefreshService struct {
	userRepository     repoModels.UserRepository
	universeRepository repoModels.UniverseRepository
	missionRepository  repoModels.MissionRepository
	upgradeConstants   map[string]constants.UpgradeConstants
	buildingConstants  map[string]map[string]map[string]interface{}
	waterConstants     constants.MiningConstants
	grapheneConstants  constants.MiningConstants
	defenceConstants   map[string]constants.DefenceConstants
	shipConstants      map[string]constants.ShipConstants
	speciesConstants   map[string]constants.SpeciesConstants
	logger             *constants.LoggingUtils
}

func NewQuickRefreshService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *QuickRefreshService {
	return &QuickRefreshService{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		missionRepository:  missionRepository,
		upgradeConstants:   upgradeConstants,
		buildingConstants:  buildingConstants,
		waterConstants:     mineConstants[constants.Water],
		grapheneConstants:  mineConstants[constants.Graphene],
		defenceConstants:   defenceConstants,
		shipConstants:      shipConstants,
		speciesConstants:   speciesConstants,
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
				r.upgradeConstants, r.buildingConstants, r.waterConstants, r.grapheneConstants,
				r.defenceConstants, r.shipConstants, r.speciesConstants[userData.Profile.Species])
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
			response.Init(planetUser, r.upgradeConstants, r.defenceConstants, r.speciesConstants[userData.Profile.Species], notifications)
			return &response, nil
		}
	}
	return nil, nil
}

func (r *QuickRefreshService) RefreshMine(username string, inputPlanetId string, inputMineId string) (*buildings.Mine, error) {
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
					response := buildings.Mine{}
					response.Init(mineUni, planetUser,
						r.upgradeConstants[constants.WaterMiningPlant], r.upgradeConstants[constants.GrapheneMiningPlant],
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

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
	militaryConstants  map[string]constants.MilitaryConstants
	researchConstants  map[string]constants.ResearchConstants
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
	militaryConstants map[string]constants.MilitaryConstants,
	researchConstants map[string]constants.ResearchConstants,
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
		militaryConstants:  militaryConstants,
		researchConstants:  researchConstants,
		speciesConstants:   speciesConstants,
		logger:             constants.NewLoggingUtils("REFRESH_SERVICE", logLevel),
	}
}

func (r *QuickRefreshService) RefreshPlanet(username string, inputPlanetId string) (*controllerModels.OccupiedPlanet, error) {
	userData, errUser := r.userRepository.FindByUsername(username)
	if errUser != nil {
		return nil, errUser
	}
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Id == inputPlanetId {
			planetUni, err := r.universeRepository.FindById(planetUser.Id)
			if err != nil {
				return nil, err
			}
			attackMissions, err := r.missionRepository.FindAttackMissionsFromPlanetId(planetUser.Id)
			if err != nil {
				r.logger.Error("error in retrieving attack missions for: "+planetUser.Id, err)
				return nil, errors.New("error in retrieving attack missions")
			}
			spyMissions, err := r.missionRepository.FindSpyMissionsFromPlanetId(planetUser.Id)
			if err != nil {
				r.logger.Error("error in retrieving spy missions for: "+planetUser.Id, err)
				return nil, errors.New("error in retrieving spy missions")
			}
			planetResponse := controllerModels.OccupiedPlanet{}
			planetResponse.Init(*planetUni, planetUser, inputPlanetId, attackMissions, spyMissions,
				r.upgradeConstants, r.buildingConstants, r.waterConstants, r.grapheneConstants,
				r.militaryConstants, r.researchConstants, r.speciesConstants[userData.Profile.Species])
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
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Id == inputPlanetId {
			response := controllerModels.UserPlanetResponse{}
			notifications := models.Notification{Tutorial: "", Error: "", Warning: ""}
			response.Init(planetUser, r.upgradeConstants, r.buildingConstants[constants.Shield], r.militaryConstants, r.speciesConstants[userData.Profile.Species], notifications)
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
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Id == inputPlanetId {
			planetUni, errUni := r.universeRepository.FindById(planetUser.Id)
			if errUni != nil {
				return nil, errUni
			}

			for _, mineUni := range planetUni.Mines {
				if mineUni.Id == inputMineId {
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
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Id == inputPlanetId {
			attackMissions, err := r.missionRepository.FindAttackMissionsFromPlanetId(planetUser.Id)
			if err != nil {
				r.logger.Error("error in retrieving attack missions for: "+planetUser.Id, err)
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
	for _, planetUser := range userData.OccupiedPlanets {
		if planetUser.Id == inputPlanetId {
			spyMissions, err := r.missionRepository.FindSpyMissionsFromPlanetId(planetUser.Id)
			if err != nil {
				r.logger.Error("error in retrieving spy missions for: "+planetUser.Id, err)
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

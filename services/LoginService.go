package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type LoginService struct {
	userRepository          repoModels.UserRepository
	clanRepository          repoModels.ClanRepository
	universeRepository      repoModels.UniverseRepository
	missionRepository       repoModels.MissionRepository
	notificationService     *NotificationService
	userExperienceConstants constants.ExperienceConstants
	clanExperienceConstants constants.ExperienceConstants
	buildingConstants       map[string]constants.BuildingConstants
	waterConstants          constants.MiningConstants
	grapheneConstants       constants.MiningConstants
	defenceConstants        map[string]constants.DefenceConstants
	shipConstants           map[string]constants.ShipConstants
	logger                  *constants.LoggingUtils
}

func NewLoginService(
	userRepository repoModels.UserRepository,
	clanRepository repoModels.ClanRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *LoginService {
	return &LoginService{
		userRepository:          userRepository,
		clanRepository:          clanRepository,
		universeRepository:      universeRepository,
		missionRepository:       missionRepository,
		notificationService:     NewNotificationService(experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		buildingConstants:       buildingConstants,
		userExperienceConstants: experienceConstants[constants.UserExperiences],
		clanExperienceConstants: experienceConstants[constants.ClanExperiences],
		waterConstants:          mineConstants[constants.Water],
		grapheneConstants:       mineConstants[constants.Graphene],
		defenceConstants:        defenceConstants,
		shipConstants:           shipConstants,
		logger:                  constants.NewLoggingUtils("LOGIN_SERVICE", logLevel),
	}
}

func (l *LoginService) Login(username string) (*controllerModels.UserResponse, error) {
	userData, err := l.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	clanData, err := getClanData(userData.Profile.ClanId, l.clanRepository)
	if err != nil {
		return nil, err
	}
	homePlanetPosition, homeSectorData, err := getHomeSectorData(userData, l.universeRepository)
	if err != nil {
		return nil, err
	}

	var response controllerModels.UserResponse
	response.Profile.Init(*userData, clanData, l.userExperienceConstants)
	homeSector, err := l.home(userData.OccupiedPlanets, *homePlanetPosition, homeSectorData)
	if err != nil {
		return nil, err
	}
	response.HomeSector = *homeSector
	response.OccupiedPlanets, err = l.occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), homeSectorData)
	if err != nil {
		return nil, err
	}
	//for _, userPlanet := range userData.OccupiedPlanets {
	//notifications, err1 := l.notificationService.getNotifications(userPlanet)
	//if err1 != nil {
	//	return nil, err1
	//}
	//response.Notifications = append(response.Notifications, notifications...)
	//}
	return &response, nil
}

func (l *LoginService) home(allOccupiedPlanetIds map[string]repoModels.PlanetUser,
	homePlanetPosition models.PlanetPosition, homeSectorData map[string]repoModels.PlanetUni) (*controllerModels.Sector, error) {

	var homeSector controllerModels.Sector
	homeSector.Position.Init(homePlanetPosition)

	for planetId, planetUni := range homeSectorData {
		if planetUser, ok := allOccupiedPlanetIds[planetId]; ok {
			planetData := controllerModels.OccupiedPlanet{}
			attackMissions, err := l.missionRepository.FindAttackMissionsFromPlanetId(planetId)
			if err != nil {
				l.logger.Error("error in retrieving attack missions for: "+planetId, err)
				return nil, errors.New("error in retrieving attack missions for: " + planetId)
			}
			spyMissions, err := l.missionRepository.FindSpyMissionsFromPlanetId(planetId)
			if err != nil {
				l.logger.Error("error in retrieving spy missions for: "+planetId, err)
				return nil, errors.New("error in retrieving spy missions for: " + planetId)
			}
			planetData.Init(planetUni, planetUser,
				attackMissions, spyMissions,
				l.buildingConstants,
				l.waterConstants, l.grapheneConstants,
				l.defenceConstants, l.shipConstants,
			)
			homeSector.OccupiedPlanets = append(homeSector.OccupiedPlanets, planetData)
			continue
		}
		planetData := controllerModels.UnoccupiedPlanet{}
		userData, err := l.userRepository.FindById(planetUni.Occupied)
		if err != nil {
			planetData.Init(planetUni, repoModels.PlanetUser{}, "")
		} else {
			planetData.Init(planetUni, userData.OccupiedPlanets[planetId], userData.Profile.Username)
		}
		homeSector.UnoccupiedPlanets = append(homeSector.UnoccupiedPlanets, planetData)
	}
	return &homeSector, nil
}

func (l *LoginService) occupiedPlanets(occupiedPlanets map[string]repoModels.PlanetUser,
	homeSectorId string, homeSectorData map[string]repoModels.PlanetUni) ([]controllerModels.StaticPlanetData, error) {
	var staticPlanets []controllerModels.StaticPlanetData
	for planetId := range occupiedPlanets {
		planetUni := homeSectorData[planetId]
		if planetUni.Id == "" {
			planet, err := l.universeRepository.FindById(planetId)
			if err != nil {
				return nil, errors.New("Error in retrieving universe data for planetId: " + planetId)
			}
			planetUni = *planet
		}
		staticPlanet := controllerModels.StaticPlanetData{}
		staticPlanet.Init(planetUni, homeSectorId)
		staticPlanets = append(staticPlanets, staticPlanet)
	}
	return staticPlanets, nil
}

func getHomeSectorData(userData *repoModels.UserData, universeRepository repoModels.UniverseRepository) (*models.PlanetPosition, map[string]repoModels.PlanetUni, error) {
	var homePlanetPosition models.PlanetPosition
	for planetId, planet := range userData.OccupiedPlanets {
		if planet.HomePlanet {
			homePlanetPosition = models.InitPlanetPositionById(planetId)
			break
		}
	}
	homeSectorData, err := universeRepository.GetSector(homePlanetPosition.System, homePlanetPosition.Sector)
	if err != nil {
		return nil, nil, err
	}
	return &homePlanetPosition, homeSectorData, nil
}

func getClanData(clanId string, clanRepository repoModels.ClanRepository) (*repoModels.ClanData, error) {
	var clanData *repoModels.ClanData
	var err error
	if len(clanId) > 0 {
		clanData, err = clanRepository.FindById(clanId)
		if err != nil {
			return nil, err
		}
		return clanData, nil
	}
	return nil, nil
}

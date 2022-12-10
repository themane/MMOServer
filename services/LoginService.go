package services

import (
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
	upgradeConstants        map[string]constants.UpgradeConstants
	buildingConstants       map[string]map[string]map[string]interface{}
	waterConstants          constants.MiningConstants
	grapheneConstants       constants.MiningConstants
	militaryConstants       map[string]constants.MilitaryConstants
	researchConstants       map[string]constants.ResearchConstants
	speciesConstants        map[string]constants.SpeciesConstants
	logger                  *constants.LoggingUtils
}

func NewLoginService(
	userRepository repoModels.UserRepository,
	clanRepository repoModels.ClanRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	researchConstants map[string]constants.ResearchConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *LoginService {
	return &LoginService{
		userRepository:          userRepository,
		clanRepository:          clanRepository,
		universeRepository:      universeRepository,
		missionRepository:       missionRepository,
		notificationService:     NewNotificationService(experienceConstants, buildingConstants, mineConstants, militaryConstants, logLevel),
		upgradeConstants:        upgradeConstants,
		buildingConstants:       buildingConstants,
		userExperienceConstants: experienceConstants[constants.UserExperiences],
		clanExperienceConstants: experienceConstants[constants.ClanExperiences],
		waterConstants:          mineConstants[constants.Water],
		grapheneConstants:       mineConstants[constants.Graphene],
		militaryConstants:       militaryConstants,
		researchConstants:       researchConstants,
		speciesConstants:        speciesConstants,
		logger:                  constants.NewLoggingUtils("LOGIN_SERVICE", logLevel),
	}
}

func (l *LoginService) LegacyLogin(username string) (*controllerModels.UserResponse, error) {
	userData, err := l.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return l.login(userData)
}

func (l *LoginService) GoogleLogin(userId string) (*controllerModels.UserResponse, error) {
	userData, err := l.userRepository.FindByGoogleId(userId)
	if err != nil {
		return nil, err
	}
	return l.login(userData)
}

func (l *LoginService) FacebookLogin(userId string) (*controllerModels.UserResponse, error) {
	userData, err := l.userRepository.FindByFacebookId(userId)
	if err != nil {
		return nil, err
	}
	return l.login(userData)
}

func (l *LoginService) login(userData *repoModels.UserData) (*controllerModels.UserResponse, error) {
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
	occupiedPlanets := map[string]repoModels.PlanetUser{}
	for _, occupiedPlanet := range userData.OccupiedPlanets {
		occupiedPlanets[occupiedPlanet.Id] = occupiedPlanet
	}
	homeSector, err := generateSectorData(occupiedPlanets,
		homePlanetPosition.SectorPosition(), homeSectorData, "",
		l.userRepository, l.missionRepository,
		l.upgradeConstants, l.buildingConstants, l.waterConstants, l.grapheneConstants,
		l.militaryConstants, l.researchConstants, l.speciesConstants[userData.Profile.Species],
		l.logger,
	)
	if err != nil {
		return nil, err
	}
	response.HomeSector = *homeSector
	response.OccupiedPlanets, err = generateOccupiedPlanetsData(occupiedPlanets,
		homePlanetPosition.SectorId(), homeSectorData, l.universeRepository)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func getHomeSectorData(userData *repoModels.UserData, universeRepository repoModels.UniverseRepository) (*models.PlanetPosition, map[string]repoModels.PlanetUni, error) {
	var homePlanetPosition *models.PlanetPosition
	for _, planet := range userData.OccupiedPlanets {
		if planet.HomePlanet {
			var err error
			homePlanetPosition, err = models.InitPlanetPositionById(planet.Id)
			if err != nil {
				return nil, nil, err
			}
			break
		}
	}
	homeSectorData, err := universeRepository.GetSector(homePlanetPosition.System, homePlanetPosition.Sector)
	if err != nil {
		return nil, nil, err
	}
	return homePlanetPosition, homeSectorData, nil
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

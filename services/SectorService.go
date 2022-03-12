package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type SectorService struct {
	userRepository      repoModels.UserRepository
	universeRepository  repoModels.UniverseRepository
	missionRepository   repoModels.MissionRepository
	notificationService *NotificationService
	upgradeConstants    map[string]constants.UpgradeConstants
	buildingConstants   map[string]map[string]map[string]interface{}
	waterConstants      constants.MiningConstants
	grapheneConstants   constants.MiningConstants
	defenceConstants    map[string]constants.DefenceConstants
	shipConstants       map[string]constants.ShipConstants
	speciesConstants    map[string]constants.SpeciesConstants
	logger              *constants.LoggingUtils
}

func NewSectorService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *SectorService {
	return &SectorService{
		userRepository:      userRepository,
		universeRepository:  universeRepository,
		missionRepository:   missionRepository,
		notificationService: NewNotificationService(experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		upgradeConstants:    upgradeConstants,
		buildingConstants:   buildingConstants,
		waterConstants:      mineConstants[constants.Water],
		grapheneConstants:   mineConstants[constants.Graphene],
		defenceConstants:    defenceConstants,
		shipConstants:       shipConstants,
		speciesConstants:    speciesConstants,
		logger:              constants.NewLoggingUtils("SECTOR_SERVICE", logLevel),
	}
}

func (s *SectorService) Visit(username string, sectorId string) (*controllerModels.SectorResponse, error) {
	userData, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	sectorPosition, err := models.InitSectorPositionById(sectorId)
	if err != nil {
		return nil, err
	}
	sectorData, err := s.universeRepository.GetSector(sectorPosition.System, sectorPosition.Sector)
	if err != nil {
		return nil, err
	}

	var response controllerModels.SectorResponse
	sector, err := generateSectorData(userData.OccupiedPlanets, *sectorPosition, sectorData, "",
		s.userRepository, s.missionRepository,
		s.upgradeConstants, s.buildingConstants, s.waterConstants, s.grapheneConstants,
		s.defenceConstants, s.shipConstants, s.speciesConstants[userData.Profile.Species],
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	response.Sector = *sector

	response.OccupiedPlanets, err = generateOccupiedPlanetsData(userData.OccupiedPlanets,
		sectorPosition.SectorId(), sectorData, s.universeRepository)
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

func (s *SectorService) Teleport(username string, planetId string) (*controllerModels.SectorResponse, error) {
	userData, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if _, ok := userData.OccupiedPlanets[planetId]; !ok {
		return nil, errors.New("not a user occupied planet")
	}
	planetPosition, err := models.InitPlanetPositionById(planetId)
	if err != nil {
		return nil, err
	}
	sectorData, err := s.universeRepository.GetSector(planetPosition.System, planetPosition.Sector)
	if err != nil {
		return nil, err
	}

	var response controllerModels.SectorResponse
	sector, err := generateSectorData(userData.OccupiedPlanets, planetPosition.SectorPosition(), sectorData, planetPosition.PlanetId(),
		s.userRepository, s.missionRepository,
		s.upgradeConstants, s.buildingConstants, s.waterConstants, s.grapheneConstants,
		s.defenceConstants, s.shipConstants, s.speciesConstants[userData.Profile.Species],
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	response.Sector = *sector

	response.OccupiedPlanets, err = generateOccupiedPlanetsData(userData.OccupiedPlanets,
		planetPosition.SectorId(), sectorData, s.universeRepository)
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

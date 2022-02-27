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
	buildingConstants   map[string]constants.BuildingConstants
	waterConstants      constants.MiningConstants
	grapheneConstants   constants.MiningConstants
	defenceConstants    map[string]constants.DefenceConstants
	shipConstants       map[string]constants.ShipConstants
	logger              *constants.LoggingUtils
}

func NewSectorService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *SectorService {
	return &SectorService{
		userRepository:      userRepository,
		universeRepository:  universeRepository,
		missionRepository:   missionRepository,
		notificationService: NewNotificationService(experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		buildingConstants:   buildingConstants,
		waterConstants:      mineConstants[constants.Water],
		grapheneConstants:   mineConstants[constants.Graphene],
		defenceConstants:    defenceConstants,
		shipConstants:       shipConstants,
		logger:              constants.NewLoggingUtils("SECTOR_SERVICE", logLevel),
	}
}

func (s *SectorService) Visit(visitRequest controllerModels.VisitSectorRequest) (*controllerModels.SectorResponse, error) {
	userData, err := s.userRepository.FindByUsername(visitRequest.Username)
	if err != nil {
		return nil, err
	}
	sectorPosition, err := models.InitSectorPositionById(visitRequest.Sector)
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
		s.buildingConstants, s.waterConstants, s.grapheneConstants, s.defenceConstants, s.shipConstants, s.logger,
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

func (s *SectorService) Teleport(teleportRequest controllerModels.TeleportRequest) (*controllerModels.SectorResponse, error) {
	userData, err := s.userRepository.FindByUsername(teleportRequest.Username)
	if err != nil {
		return nil, err
	}
	if _, ok := userData.OccupiedPlanets[teleportRequest.Planet]; !ok {
		return nil, errors.New("not a user occupied planet")
	}
	planetPosition, err := models.InitPlanetPositionById(teleportRequest.Planet)
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
		s.buildingConstants, s.waterConstants, s.grapheneConstants, s.defenceConstants, s.shipConstants, s.logger,
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

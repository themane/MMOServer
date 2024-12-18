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
	militaryConstants   map[string]constants.MilitaryConstants
	researchConstants   map[string]constants.ResearchConstants
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
	militaryConstants map[string]constants.MilitaryConstants,
	researchConstants map[string]constants.ResearchConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *SectorService {
	return &SectorService{
		userRepository:      userRepository,
		universeRepository:  universeRepository,
		missionRepository:   missionRepository,
		notificationService: NewNotificationService(experienceConstants, buildingConstants, mineConstants, militaryConstants, logLevel),
		upgradeConstants:    upgradeConstants,
		buildingConstants:   buildingConstants,
		waterConstants:      mineConstants[constants.Water],
		grapheneConstants:   mineConstants[constants.Graphene],
		militaryConstants:   militaryConstants,
		researchConstants:   researchConstants,
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
	occupiedPlanets := map[string]repoModels.PlanetUser{}
	for _, occupiedPlanet := range userData.OccupiedPlanets {
		occupiedPlanets[occupiedPlanet.Id] = occupiedPlanet
	}
	sector, err := generateSectorData(occupiedPlanets, *sectorPosition, sectorData, "",
		s.userRepository, s.missionRepository,
		s.upgradeConstants, s.buildingConstants, s.waterConstants, s.grapheneConstants,
		s.militaryConstants, s.researchConstants, s.speciesConstants[userData.Profile.Species],
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	response.Sector = *sector

	response.OccupiedPlanets, err = generateOccupiedPlanetsData(occupiedPlanets,
		sectorPosition.SectorId(), sectorData, s.universeRepository)
	response.Notifications = models.Notification{Tutorial: "", Error: "", Warning: ""}
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *SectorService) Teleport(username string, planetId string) (*controllerModels.SectorResponse, error) {
	userData, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if userData.GetOccupiedPlanet(planetId) == nil {
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
	occupiedPlanets := map[string]repoModels.PlanetUser{}
	for _, occupiedPlanet := range userData.OccupiedPlanets {
		occupiedPlanets[occupiedPlanet.Id] = occupiedPlanet
	}
	sector, err := generateSectorData(occupiedPlanets, planetPosition.SectorPosition(), sectorData, planetPosition.PlanetId(),
		s.userRepository, s.missionRepository,
		s.upgradeConstants, s.buildingConstants, s.waterConstants, s.grapheneConstants,
		s.militaryConstants, s.researchConstants, s.speciesConstants[userData.Profile.Species],
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	response.Sector = *sector

	response.OccupiedPlanets, err = generateOccupiedPlanetsData(occupiedPlanets,
		planetPosition.SectorId(), sectorData, s.universeRepository)
	response.Notifications = models.Notification{Tutorial: "", Error: "", Warning: ""}
	if err != nil {
		return nil, err
	}

	return &response, nil
}

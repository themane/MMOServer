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
	userExperienceConstants constants.ExperienceConstants
	clanExperienceConstants constants.ExperienceConstants
	buildingConstants       map[string]constants.BuildingConstants
	waterConstants          constants.MiningConstants
	grapheneConstants       constants.MiningConstants
	defenceConstants        map[string]constants.DefenceConstants
	shipConstants           map[string]constants.ShipConstants
}

func NewLoginService(
	userRepository repoModels.UserRepository,
	clanRepository repoModels.ClanRepository,
	universeRepository repoModels.UniverseRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
) *LoginService {
	return &LoginService{
		userRepository:          userRepository,
		clanRepository:          clanRepository,
		universeRepository:      universeRepository,
		buildingConstants:       buildingConstants,
		userExperienceConstants: experienceConstants[constants.UserExperiences],
		clanExperienceConstants: experienceConstants[constants.ClanExperiences],
		waterConstants:          mineConstants[constants.Water],
		grapheneConstants:       mineConstants[constants.Graphene],
		defenceConstants:        defenceConstants,
		shipConstants:           shipConstants,
	}
}

func (l *LoginService) Login(username string) (*controllerModels.LoginResponse, error) {
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

	var response controllerModels.LoginResponse
	response.Profile.Init(*userData, clanData, l.userExperienceConstants)
	response.HomeSector = l.home(userData.OccupiedPlanets, *homePlanetPosition, homeSectorData)
	response.OccupiedPlanets, err = l.occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), homeSectorData)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (l *LoginService) home(allOccupiedPlanetIds map[string]repoModels.PlanetUser,
	homePlanetPosition models.PlanetPosition, homeSectorData map[string]repoModels.PlanetUni) controllerModels.Sector {

	var homeSector controllerModels.Sector
	homeSector.Position.Init(homePlanetPosition)

	for planetId, planetUni := range homeSectorData {
		if planetUser, ok := allOccupiedPlanetIds[planetId]; ok {
			planetData := controllerModels.OccupiedPlanet{}
			planetData.Init(planetUni, planetUser,
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
	return homeSector
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
		if planet.Home {
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

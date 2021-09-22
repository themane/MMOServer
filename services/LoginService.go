package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

func Login(username string,
	userRepository repoModels.UserRepository,
	clanRepository repoModels.ClanRepository,
	universeRepository repoModels.UniverseRepository,
	waterConstants constants.ResourceConstants,
	grapheneConstants constants.ResourceConstants,
	experienceConstants constants.ExperienceConstants,
) (*controllerModels.LoginResponse, error) {

	userData, err := userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	clanData, err := getClanData(userData.Profile.ClanId, clanRepository)
	if err != nil {
		return nil, err
	}
	homePlanetPosition, homeSectorData, err := getHomeSectorData(userData, universeRepository)
	if err != nil {
		return nil, err
	}

	var response controllerModels.LoginResponse
	response.Profile.Init(*userData, clanData, experienceConstants)
	response.HomeSector = home(userData.OccupiedPlanets, *homePlanetPosition, homeSectorData, waterConstants, grapheneConstants)
	response.OccupiedPlanets, err = occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), homeSectorData, universeRepository)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func home(allOccupiedPlanetIds map[string]repoModels.PlanetUser,
	homePlanetPosition models.PlanetPosition,
	homeSectorData map[string]repoModels.PlanetUni,
	waterConstants constants.ResourceConstants,
	grapheneConstants constants.ResourceConstants) controllerModels.Sector {

	var homeSector controllerModels.Sector

	homeSector.Position.Init(homePlanetPosition)

	for planetId, planetUni := range homeSectorData {
		if planetUser, ok := allOccupiedPlanetIds[planetId]; ok {
			planetData := controllerModels.OccupiedPlanet{}
			planetData.Init(planetUni, planetUser, waterConstants, grapheneConstants)
			homeSector.OccupiedPlanets = append(homeSector.OccupiedPlanets, planetData)
			continue
		}
		planetData := controllerModels.UnoccupiedPlanet{}
		planetData.Init(planetUni)
		homeSector.UnoccupiedPlanets = append(homeSector.UnoccupiedPlanets, planetData)
	}
	return homeSector
}

func occupiedPlanets(occupiedPlanets map[string]repoModels.PlanetUser, homeSectorId string, homeSectorData map[string]repoModels.PlanetUni,
	universeRepository repoModels.UniverseRepository) ([]controllerModels.StaticPlanetData, error) {
	var staticPlanets []controllerModels.StaticPlanetData
	for planetId := range occupiedPlanets {
		planetUni := homeSectorData[planetId]
		if planetUni.Id == "" {
			planet, err := universeRepository.FindById(planetId)
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

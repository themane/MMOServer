package services

import (
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

func Login(username string, userRepository repoModels.UserRepository, clanRepository repoModels.ClanRepository, universeRepository repoModels.UniverseRepository) (*controllerModels.LoginResponse, error) {
	waterConstants := constants.GetWaterConstants()
	grapheneConstants := constants.GetGrapheneConstants()
	experienceConstants := constants.GetExperienceConstants()

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
	response.HomeSector, response.HomePlanetId = home(userData.OccupiedPlanets, *homePlanetPosition, homeSectorData, waterConstants, grapheneConstants)
	response.OccupiedPlanets = occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), homeSectorData)
	return &response, nil
}

func home(allOccupiedPlanetIds map[string]repoModels.PlanetUser,
	homePlanetPosition models.PlanetPosition,
	homeSectorData map[string]repoModels.PlanetUni,
	waterConstants constants.ResourceConstants,
	grapheneConstants constants.ResourceConstants) (controllerModels.Sector, string) {

	var homeSector controllerModels.Sector
	var homePlanetId string

	homeSector.Position.Init(homePlanetPosition)

	for planetId, planetUni := range homeSectorData {
		if planetUser, ok := allOccupiedPlanetIds[planetId]; ok {
			planetData := controllerModels.OccupiedPlanet{}
			planetData.Init(planetUni, planetUser, waterConstants, grapheneConstants)
			if planetData.Home {
				homePlanetId = planetId
			}
			homeSector.OccupiedPlanets = append(homeSector.OccupiedPlanets, planetData)
			continue
		}
		planetData := controllerModels.UnoccupiedPlanet{}
		planetData.Init(planetUni)
		homeSector.UnoccupiedPlanets = append(homeSector.UnoccupiedPlanets, planetData)
	}
	return homeSector, homePlanetId
}

func occupiedPlanets(occupiedPlanets map[string]repoModels.PlanetUser, homeSectorId string, homeSectorData map[string]repoModels.PlanetUni) []controllerModels.StaticPlanetData {
	var staticPlanets []controllerModels.StaticPlanetData
	for planetId := range occupiedPlanets {
		planetUni := homeSectorData[planetId]
		staticPlanet := controllerModels.StaticPlanetData{}
		staticPlanet.Init(planetUni, homeSectorId)
		staticPlanets = append(staticPlanets, staticPlanet)
	}
	return staticPlanets
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
	if clanId != "" {
		clanData, err = clanRepository.FindById(clanId)
		if err != nil {
			return nil, err
		}
		return clanData, nil
	}
	return nil, nil
}

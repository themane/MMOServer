package services

import (
	"github.com/themane/MMOServer/dao"
	"github.com/themane/MMOServer/models"
)

func Login(username string) models.LoginResponse {
	universe := dao.GetUniverse()
	waterConstants := dao.GetWaterConstants()
	grapheneConstants := dao.GetGrapheneConstants()
	userData := dao.GetUserData(username)
	homePlanetPosition := findHomePlanet(userData.OccupiedPlanets)

	var response models.LoginResponse
	response.Profile = userData.Profile
	response.HomeSector, response.HomePlanet = home(userData.OccupiedPlanets, *homePlanetPosition, universe, waterConstants, grapheneConstants)
	response.OccupiedPlanets = occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), universe)

	return response
}

func home(allOccupiedPlanetIds map[string]models.PlanetUser, homePlanetPosition models.PlanetPosition, universe models.Universe,
	waterConstants models.ResourceConstants, grapheneConstants models.ResourceConstants) (models.Sector, models.OccupiedPlanet) {
	sectorU := universe.Systems[homePlanetPosition.SystemId()].Sectors[homePlanetPosition.SectorId()]
	var homeSector models.Sector
	var homePlanet models.OccupiedPlanet

	homeSector.Position.Init(homePlanetPosition)

	for key, planetUni := range sectorU.Planets {
		if planetUser, ok := allOccupiedPlanetIds[key]; ok {
			planetData := models.OccupiedPlanet{}
			planetData.Init(planetUni, planetUser, homeSector.Position, waterConstants, grapheneConstants)
			if planetData.Home {
				homePlanet = planetData
			}
			homeSector.OccupiedPlanets = append(homeSector.OccupiedPlanets, planetData)
			continue
		}
		planetData := models.UnoccupiedPlanet{}
		planetData.Init(planetUni, homeSector.Position)
		homeSector.UnoccupiedPlanets = append(homeSector.UnoccupiedPlanets, planetData)
	}
	return homeSector, homePlanet
}

func occupiedPlanets(occupiedPlanets map[string]models.PlanetUser, homeSectorId string, universe models.Universe) []models.StaticPlanetData {
	var staticPlanets []models.StaticPlanetData
	for _, planetUser := range occupiedPlanets {
		planetUni := universe.Systems[planetUser.Position.SystemId()].Sectors[planetUser.Position.SectorId()].Planets[planetUser.Position.PlanetId()]
		staticPlanet := models.StaticPlanetData{}
		staticPlanet.Init(planetUser, planetUni, homeSectorId)
		staticPlanets = append(staticPlanets, staticPlanet)
	}
	return staticPlanets
}

func findHomePlanet(occupiedPlanets map[string]models.PlanetUser) *models.PlanetPosition {
	for _, planet := range occupiedPlanets {
		if planet.Home {
			return &planet.Position
		}
	}
	return nil
}

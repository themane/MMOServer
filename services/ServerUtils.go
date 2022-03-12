package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

func generateSectorData(allOccupiedPlanetIds map[string]repoModels.PlanetUser,
	sectorPosition models.SectorPosition, sectorData map[string]repoModels.PlanetUni, customHomePlanetId string,
	userRepository repoModels.UserRepository,
	missionRepository repoModels.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	waterConstants constants.MiningConstants,
	grapheneConstants constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	speciesConstants constants.SpeciesConstants,
	logger *constants.LoggingUtils,
) (*controllerModels.Sector, error) {

	var homeSector controllerModels.Sector
	homeSector.Position = sectorPosition

	for planetId, planetUni := range sectorData {
		if planetUser, ok := allOccupiedPlanetIds[planetId]; ok {
			planetData := controllerModels.OccupiedPlanet{}
			attackMissions, err := missionRepository.FindAttackMissionsFromPlanetId(planetId)
			if err != nil {
				logger.Error("error in retrieving attack missions for: "+planetId, err)
				return nil, errors.New("error in retrieving attack missions for: " + planetId)
			}
			spyMissions, err := missionRepository.FindSpyMissionsFromPlanetId(planetId)
			if err != nil {
				logger.Error("error in retrieving spy missions for: "+planetId, err)
				return nil, errors.New("error in retrieving spy missions for: " + planetId)
			}
			planetData.Init(planetUni, planetUser, customHomePlanetId,
				attackMissions, spyMissions,
				upgradeConstants, buildingConstants,
				waterConstants, grapheneConstants,
				defenceConstants, shipConstants, speciesConstants,
			)
			homeSector.OccupiedPlanets = append(homeSector.OccupiedPlanets, planetData)
			continue
		}
		planetData := controllerModels.UnoccupiedPlanet{}
		userData, err := userRepository.FindById(planetUni.Occupied)
		if err != nil {
			planetData.Init(planetUni, repoModels.PlanetUser{}, "")
		} else {
			planetData.Init(planetUni, userData.OccupiedPlanets[planetId], userData.Profile.Username)
		}
		homeSector.UnoccupiedPlanets = append(homeSector.UnoccupiedPlanets, planetData)
	}
	return &homeSector, nil
}

func generateOccupiedPlanetsData(occupiedPlanets map[string]repoModels.PlanetUser,
	sectorId string, sectorData map[string]repoModels.PlanetUni,
	universeRepository repoModels.UniverseRepository,
) ([]controllerModels.StaticPlanetData, error) {
	var staticPlanets []controllerModels.StaticPlanetData
	for planetId := range occupiedPlanets {
		planetUni := sectorData[planetId]
		if planetUni.Id == "" {
			planet, err := universeRepository.FindById(planetId)
			if err != nil {
				return nil, errors.New("Error in retrieving universe data for planetId: " + planetId)
			}
			planetUni = *planet
		}
		staticPlanet := controllerModels.StaticPlanetData{}
		staticPlanet.Init(planetUni, sectorId)
		staticPlanets = append(staticPlanets, staticPlanet)
	}
	return staticPlanets, nil
}

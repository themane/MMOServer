package schedulers

import (
	"github.com/go-co-op/gocron"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"log"
	"strconv"
	"time"
)

type ScheduledJobManager struct {
	userRepository     models.UserRepository
	universeRepository models.UniverseRepository
	waterConstants     constants.ResourceConstants
	grapheneConstants  constants.ResourceConstants
	maxSystem          int
}

func NewScheduledJobManager(userRepository *models.UserRepository, universeRepository *models.UniverseRepository,
	waterConstants constants.ResourceConstants, grapheneConstants constants.ResourceConstants,
	maxSystem int) *ScheduledJobManager {
	return &ScheduledJobManager{
		userRepository:     *userRepository,
		universeRepository: *universeRepository,
		waterConstants:     waterConstants,
		grapheneConstants:  grapheneConstants,
		maxSystem:          maxSystem,
	}
}

func (j *ScheduledJobManager) SchedulePlanetUpdates() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Hour().Do(j.scheduledPopulationIncrease)
	if err != nil {
		log.Print(err)
	}
	_, err1 := s.Every(1).Minutes().Do(j.scheduledMining)
	if err1 != nil {
		log.Print(err1)
	}
}

func (j *ScheduledJobManager) scheduledPopulationIncrease() {
	for system := 0; system < j.maxSystem; system++ {
		occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
		if err != nil {
			log.Print(err)
			return
		}
		var userIdplanetsMap map[string][]string
		for planetId, occupiedPlanet := range occupiedPlanets {
			userIdplanetsMap[occupiedPlanet.Occupied] = append(userIdplanetsMap[occupiedPlanet.Occupied], planetId)
		}
		for userId, planets := range userIdplanetsMap {
			planetIdGenerationRateMap := j.getPopulationGenerationRate(userId, planets)
			err := j.userRepository.ScheduledPopulationIncrease(userId, planetIdGenerationRateMap)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func (j *ScheduledJobManager) scheduledMining() {
	for system := 0; system < j.maxSystem; system++ {
		occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
		if err != nil {
			log.Print(err)
			return
		}
		var userIdplanetsMap map[string][]models.PlanetUni
		for _, occupiedPlanet := range occupiedPlanets {
			userIdplanetsMap[occupiedPlanet.Occupied] = append(userIdplanetsMap[occupiedPlanet.Occupied], occupiedPlanet)
		}
		for userId, planets := range userIdplanetsMap {
			planetIdWaterMiningRateMap, planetIdGrapheneMiningRateMap := j.getMiningRate(userId, planets)
			err := j.userRepository.ScheduledWaterIncrease(userId, planetIdWaterMiningRateMap)
			if err != nil {
				log.Print(err)
				return
			}
			err = j.userRepository.ScheduledGrapheneIncrease(userId, planetIdGrapheneMiningRateMap)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func (j *ScheduledJobManager) getPopulationGenerationRate(userId string, occupiedPlanets []string) map[string]int {
	userData, err := j.userRepository.FindById(userId)
	if err != nil {
		log.Print(err)
		return nil
	}
	var planetIdGenerationRateMap map[string]int
	for _, planetId := range occupiedPlanets {
		generationRate := userData.OccupiedPlanets[planetId].Population.GenerationRate
		planetIdGenerationRateMap[planetId] = generationRate
	}
	return planetIdGenerationRateMap
}

func (j *ScheduledJobManager) getMiningRate(userId string, occupiedPlanets []models.PlanetUni) (map[string]map[string]int, map[string]map[string]int) {
	userData, err := j.userRepository.FindById(userId)
	if err != nil {
		log.Print(err)
		return nil, nil
	}
	var planetIdWaterMiningRateMap map[string]map[string]int
	var planetIdGrapheneMiningRateMap map[string]map[string]int
	for _, planetUni := range occupiedPlanets {
		planetUser := userData.OccupiedPlanets[planetUni.Id]
		for _, mineUni := range planetUni.Mines {
			mineUser := planetUser.Mines[mineUni.Id]
			miningPlant := planetUser.Buildings[mineUser.MiningPlantId]

			var miningRatePerWorker int
			if mineUni.Type == constants.WATER {
				miningRatePerWorker = j.waterConstants.Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := j.getTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				planetIdWaterMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
			if mineUni.Type == constants.GRAPHENE {
				miningRatePerWorker = j.grapheneConstants.Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := j.getTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				planetIdGrapheneMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
		}
	}
	return planetIdWaterMiningRateMap, planetIdGrapheneMiningRateMap
}

func (j *ScheduledJobManager) getTotalMiningRate(miningRatePerWorker int, miningPlantWorkers int, maxMinedLimit int, minedResource int) int {
	miningRate := miningRatePerWorker * miningPlantWorkers
	if maxMinedLimit < (minedResource + miningRate) {
		miningRate = maxMinedLimit - (minedResource + miningRate)
	}
	return miningRate
}

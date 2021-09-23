package schedulers

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"log"
	"strconv"
)

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
			if mineUni.Type == constants.Water {
				miningRatePerWorker = j.waterConstants.Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := j.getTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				planetIdWaterMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
			if mineUni.Type == constants.Graphene {
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

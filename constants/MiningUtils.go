package constants

import (
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

func GetMiningRate(userData models.UserData, occupiedPlanets []models.PlanetUni,
	mineConstants map[string]MiningConstants,
) (map[string]map[string]int, map[string]map[string]int) {

	planetIdWaterMiningRateMap := map[string]map[string]int{}
	planetIdGrapheneMiningRateMap := map[string]map[string]int{}
	for _, planetUni := range occupiedPlanets {
		planetUser := userData.OccupiedPlanets[planetUni.Id]
		for _, mineUni := range planetUni.Mines {
			mineUser := planetUser.Mines[mineUni.Id]
			miningPlant := planetUser.Buildings[mineUser.MiningPlantId]

			var miningRatePerWorker int
			if mineUni.Type == Water {
				miningRatePerWorker = mineConstants[Water].Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := GetTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				if planetIdWaterMiningRateMap[planetUni.Id] == nil {
					planetIdWaterMiningRateMap[planetUni.Id] = map[string]int{}
				}
				planetIdWaterMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
			if mineUni.Type == Graphene {
				miningRatePerWorker = mineConstants[Graphene].Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := GetTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				if planetIdGrapheneMiningRateMap[planetUni.Id] == nil {
					planetIdGrapheneMiningRateMap[planetUni.Id] = map[string]int{}
				}
				planetIdGrapheneMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
		}
	}
	return planetIdWaterMiningRateMap, planetIdGrapheneMiningRateMap
}

func GetTotalMiningRate(miningRatePerWorker int, miningPlantWorkers int, maxMinedLimit int, minedResource int) int {
	miningRate := miningRatePerWorker * miningPlantWorkers
	if maxMinedLimit < (minedResource + miningRate) {
		miningRate = maxMinedLimit - (minedResource + miningRate)
	}
	return miningRate
}

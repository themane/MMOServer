package models

import (
	"github.com/themane/MMOServer/constants"
	"strconv"
)

func GetMiningRate(userData UserData, occupiedPlanets []PlanetUni,
	waterConstants constants.MiningConstants,
	grapheneConstants constants.MiningConstants,
) (map[string]map[string]int, map[string]map[string]int) {

	planetIdWaterMiningRateMap := map[string]map[string]int{}
	planetIdGrapheneMiningRateMap := map[string]map[string]int{}
	for _, planetUni := range occupiedPlanets {
		planetUser := userData.OccupiedPlanets[planetUni.Id]
		for _, mineUni := range planetUni.Mines {
			mineUser := planetUser.Mines[mineUni.Id]
			miningPlant := planetUser.Buildings[mineUser.MiningPlantId]

			var miningRatePerWorker int
			if mineUni.Type == constants.Water {
				miningRatePerWorker = waterConstants.Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
				miningRate := GetTotalMiningRate(miningRatePerWorker, miningPlant.Workers, mineUni.MaxLimit, mineUser.Mined)
				if planetIdWaterMiningRateMap[planetUni.Id] == nil {
					planetIdWaterMiningRateMap[planetUni.Id] = map[string]int{}
				}
				planetIdWaterMiningRateMap[planetUni.Id][mineUni.Id] = miningRate
			}
			if mineUni.Type == constants.Graphene {
				miningRatePerWorker = grapheneConstants.Levels[strconv.Itoa(miningPlant.BuildingLevel)].MiningRatePerWorker
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
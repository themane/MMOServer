package schedulers

import (
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/themane/MMOServer/dao"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func SchedulePlanetUpdates() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Hour().Do(scheduledPopulationUpdates)
	if err != nil {
		log.Print(err)
	}
	_, err1 := s.Every(1).Minutes().Do(scheduledResourcesUpdates)
	if err1 != nil {
		log.Print(err1)
	}
}

func scheduledPopulationUpdates() {
	files, err := ioutil.ReadDir("sample_data/users")
	if err != nil {
		log.Print(err)
	}
	for _, file := range files {
		var userData models.UserData
		jsonFile, _ := os.Open(file.Name())
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, planet := range userData.OccupiedPlanets {
			totalPopulation := planet.Population.Unemployed + planet.Population.Workers.Total + planet.Population.Soldiers.Total
			if planet.Water.Amount >= totalPopulation {
				planet.Water.Amount -= totalPopulation
			} else {
				extraPopulation := totalPopulation - planet.Water.Amount
				for i := 0; i < extraPopulation; i++ {
					isWorker := rand.Intn(2)
					if isWorker == 0 {
						planet.Population.Soldiers.Total -= extraPopulation
					} else {
						planet.Population.Workers.Total -= extraPopulation
					}
				}
				planet.Water.Amount = 0
			}
			planet.Population.Unemployed += planet.Population.GenerationRate
		}
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile(file.Name(), updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	}
}

func scheduledResourcesUpdates() {
	waterConstants := dao.GetWaterConstants()
	grapheneConstants := dao.GetGrapheneConstants()

	files, err := ioutil.ReadDir("sample_data/users")
	if err != nil {
		log.Print(err)
	}

	for _, file := range files {
		var userData models.UserData
		jsonFile, _ := os.Open(file.Name())
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, planet := range userData.OccupiedPlanets {
			totalWaterMined := 0
			totalGrapheneMined := 0
			for _, mine := range planet.Mines {
				var miningRate int
				levelString := strconv.Itoa(mine.MiningPlant.BuildingLevel)
				if mine.Type == models.WATER {
					miningRate = mine.MiningPlant.Workers * waterConstants.Levels[levelString].MiningRatePerWorker
					totalWaterMined += miningRate
				}
				if mine.Type == models.GRAPHENE {
					miningRate = mine.MiningPlant.Workers * grapheneConstants.Levels[levelString].MiningRatePerWorker
					totalGrapheneMined += miningRate
				}
				mine.Mined += miningRate
			}
			planet.Water.Amount += totalWaterMined
			planet.Graphene.Amount += totalGrapheneMined
		}
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile(file.Name(), updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	}
}

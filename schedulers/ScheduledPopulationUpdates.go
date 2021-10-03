package schedulers

import (
	"github.com/themane/MMOServer/constants"
	"log"
)

func (j *ScheduledJobManager) scheduledPopulationIncrease() {
	log.Println("Scheduled run of increasing population")
	for system := 0; system < j.maxSystem; system++ {
		occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
		if err != nil {
			log.Print("Error in getting all occupied planets", err)
			return
		}
		userIdplanetsMap := map[string][]string{}
		for planetId, occupiedPlanet := range occupiedPlanets {
			planetType := constants.GetPlanetType(occupiedPlanet)
			if planetType == constants.User {
				if userIdplanetsMap[occupiedPlanet.Occupied] == nil {
					userIdplanetsMap[occupiedPlanet.Occupied] = []string{}
				}
				userIdplanetsMap[occupiedPlanet.Occupied] = append(userIdplanetsMap[occupiedPlanet.Occupied], planetId)
			}
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

func (j *ScheduledJobManager) getPopulationGenerationRate(userId string, occupiedPlanets []string) map[string]int {
	userData, err := j.userRepository.FindById(userId)
	if err != nil {
		log.Print(err)
		return nil
	}
	planetIdGenerationRateMap := map[string]int{}
	for _, planetId := range occupiedPlanets {
		generationRate := userData.OccupiedPlanets[planetId].Population.GenerationRate
		planetIdGenerationRateMap[planetId] = generationRate
	}
	return planetIdGenerationRateMap
}

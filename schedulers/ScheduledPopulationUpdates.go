package schedulers

import (
	"github.com/themane/MMOServer/constants"
)

func (j *ScheduledJobManager) scheduledPopulationIncrease() {
	j.logger.Info("Scheduled run of increasing population")
	for system := 0; system < j.maxSystem; system++ {
		occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
		if err != nil {
			j.logger.Error("error in getting all occupied planets", err)
			return
		}
		userIdplanetsMap := map[string][]string{}
		for planetId, occupiedPlanet := range occupiedPlanets {
			planetType := occupiedPlanet.GetPlanetType()
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
				j.logger.Error("error in firing update query for scheduled population increase, userId: "+userId, err)
				return
			}
		}
	}
}

func (j *ScheduledJobManager) getPopulationGenerationRate(userId string, occupiedPlanets []string) map[string]int {
	userData, err := j.userRepository.FindById(userId)
	if err != nil {
		j.logger.Error("error in getting user data for userId: "+userId, err)
		return nil
	}
	planetIdGenerationRateMap := map[string]int{}
	for _, planetId := range occupiedPlanets {
		generationRate := userData.OccupiedPlanets[planetId].Population.GenerationRate
		planetIdGenerationRateMap[planetId] = generationRate
	}
	return planetIdGenerationRateMap
}

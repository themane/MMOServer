package schedulers

//import (
//	"strconv"
//)

func (j *ScheduledJobManager) retrieveMiningAndPopulationData() {
	//for system := 0; system < j.maxSystem; system++ {
	//	occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
	//	if err != nil {
	//		j.logger.Error("error in retrieving all occupied planets for system: "+strconv.Itoa(system), err)
	//		return
	//	}
	//	var allOccupierIds []string
	//	for _, occupiedPlanet := range occupiedPlanets {
	//		allOccupierIds = append(allOccupierIds, occupiedPlanet.Occupied)
	//	}
	//
	//	j.userRepository.GetMiningAndPopulationData(allOccupierIds)
	//}
}

func (j *ScheduledJobManager) miningUpdates() {
	//j.logger.Info("Scheduled run of mining")
	//for system := 0; system < j.maxSystem; system++ {
	//	occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
	//	if err != nil {
	//		j.logger.Error("error in retrieving all occupied planets for system: "+strconv.Itoa(system), err)
	//		return
	//	}
	//	userIdplanetsMap := map[string][]models.PlanetUni{}
	//	for _, occupiedPlanet := range occupiedPlanets {
	//		planetType := occupiedPlanet.GetPlanetType()
	//		if planetType == constants.User {
	//			if userIdplanetsMap[occupiedPlanet.Occupied] == nil {
	//				userIdplanetsMap[occupiedPlanet.Occupied] = []models.PlanetUni{}
	//			}
	//			userIdplanetsMap[occupiedPlanet.Occupied] = append(userIdplanetsMap[occupiedPlanet.Occupied], occupiedPlanet)
	//		}
	//	}
	//	for userId, planets := range userIdplanetsMap {
	//		userData, err1 := j.userRepository.FindById(userId)
	//		if err1 != nil {
	//			j.logger.Error("error in retrieving user data for: "+userId, err)
	//			return
	//		}
	//		planetIdWaterMiningRateMap, planetIdGrapheneMiningRateMap := models.GetMiningRate(*userData, planets, j.waterConstants, j.grapheneConstants)
	//		err1 = j.userRepository.ScheduledWaterIncrease(userId, planetIdWaterMiningRateMap)
	//		if err1 != nil {
	//			j.logger.Error("error in water increase update for user: "+userId, err)
	//			return
	//		}
	//		err1 = j.userRepository.ScheduledGrapheneIncrease(userId, planetIdGrapheneMiningRateMap)
	//		if err1 != nil {
	//			j.logger.Error("error in graphene increase update for user: "+userId, err)
	//			return
	//		}
	//	}
	//}
}

//func (j *ScheduledJobManager) scheduledMining() {
//j.logger.Info("Scheduled run of mining")
//for system := 0; system < j.maxSystem; system++ {
//	occupiedPlanets, err := j.universeRepository.GetAllOccupiedPlanets(system)
//	if err != nil {
//		j.logger.Error("error in retrieving all occupied planets for system: "+strconv.Itoa(system), err)
//		return
//	}
//	userIdplanetsMap := map[string][]models.PlanetUni{}
//	for _, occupiedPlanet := range occupiedPlanets {
//		planetType := occupiedPlanet.GetPlanetType()
//		if planetType == constants.User {
//			if userIdplanetsMap[occupiedPlanet.Occupied] == nil {
//				userIdplanetsMap[occupiedPlanet.Occupied] = []models.PlanetUni{}
//			}
//			userIdplanetsMap[occupiedPlanet.Occupied] = append(userIdplanetsMap[occupiedPlanet.Occupied], occupiedPlanet)
//		}
//	}
//	for userId, planets := range userIdplanetsMap {
//		userData, err1 := j.userRepository.FindById(userId)
//		if err1 != nil {
//			j.logger.Error("error in retrieving user data for: "+userId, err)
//			return
//		}
//		planetIdWaterMiningRateMap, planetIdGrapheneMiningRateMap := models.GetMiningRate(*userData, planets, j.waterConstants, j.grapheneConstants)
//		err1 = j.userRepository.ScheduledWaterIncrease(userId, planetIdWaterMiningRateMap)
//		if err1 != nil {
//			j.logger.Error("error in water increase update for user: "+userId, err)
//			return
//		}
//		err1 = j.userRepository.ScheduledGrapheneIncrease(userId, planetIdGrapheneMiningRateMap)
//		if err1 != nil {
//			j.logger.Error("error in graphene increase update for user: "+userId, err)
//			return
//		}
//	}
//}
//}

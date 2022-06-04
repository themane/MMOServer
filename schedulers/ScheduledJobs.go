package schedulers

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
)

type ScheduledJobManager struct {
	userRepository     models.UserRepository
	universeRepository models.UniverseRepository
	waterConstants     constants.MiningConstants
	grapheneConstants  constants.MiningConstants
	maxSystem          int
	logger             *constants.LoggingUtils
}

func NewScheduledJobManager(userRepository models.UserRepository, universeRepository models.UniverseRepository,
	mineConstants map[string]constants.MiningConstants, maxSystem int,
	logLevel string,
) *ScheduledJobManager {
	return &ScheduledJobManager{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		waterConstants:     mineConstants[constants.Water],
		grapheneConstants:  mineConstants[constants.Graphene],
		maxSystem:          maxSystem,
		logger:             constants.NewLoggingUtils("SCHEDULER_JOBS", logLevel),
	}
}

//func (j *ScheduledJobManager) SchedulePlanetUpdates() {
//	s := gocron.NewScheduler(time.UTC)
//	_, err := s.Every(10).Hour().Do(j.scheduledPopulationIncrease)
//	if err != nil {
//		j.logger.Error("error in scheduled population increase", err)
//	}
//	_, err1 := s.Every(10).Minute().Do(j.scheduledMining)
//	if err1 != nil {
//		j.logger.Error("error in scheduled mining increase", err1)
//	}
//	s.StartAsync()
//}

package schedulers

import (
	"github.com/go-co-op/gocron"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"log"
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

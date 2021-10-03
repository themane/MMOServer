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
}

func NewScheduledJobManager(userRepository models.UserRepository, universeRepository models.UniverseRepository,
	mineConstants map[string]constants.MiningConstants, maxSystem int) *ScheduledJobManager {
	return &ScheduledJobManager{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		waterConstants:     mineConstants[constants.Water],
		grapheneConstants:  mineConstants[constants.Graphene],
		maxSystem:          maxSystem,
	}
}

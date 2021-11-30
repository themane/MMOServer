package services

import (
	"github.com/themane/MMOServer/constants"
)

type NotificationService struct {
	userExperienceConstants constants.ExperienceConstants
	clanExperienceConstants constants.ExperienceConstants
	buildingConstants       map[string]constants.BuildingConstants
	waterConstants          constants.MiningConstants
	grapheneConstants       constants.MiningConstants
	defenceConstants        map[string]constants.DefenceConstants
	shipConstants           map[string]constants.ShipConstants
	logger                  *constants.LoggingUtils
}

func NewNotificationService(
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *NotificationService {
	return &NotificationService{
		buildingConstants:       buildingConstants,
		userExperienceConstants: experienceConstants[constants.UserExperiences],
		clanExperienceConstants: experienceConstants[constants.ClanExperiences],
		waterConstants:          mineConstants[constants.Water],
		grapheneConstants:       mineConstants[constants.Graphene],
		defenceConstants:        defenceConstants,
		shipConstants:           shipConstants,
		logger:                  constants.NewLoggingUtils("NOTIFICATION_SERVICE", logLevel),
	}
}

//func (l *NotificationService) getNotifications(occupiedPlanet controllerModels.OccupiedPlanet) ([]controllerModels.Notification, error) {
//	waterAvailable := occupiedPlanet.Resources.Water.Amount
//	for _, mineData := range occupiedPlanet.Mines {
//		workers :=
//		currentMiningRate := mineData.MiningPlant.NextLevelMiningAttributes.CurrentMiningRatePerWorker*mineData.MiningPlant.Workers
//	}
//	userPlanet.Mines
//
//	userData, err := l.userRepository.FindByUsername(username)
//	if err != nil {
//		return nil, err
//	}
//	clanData, err := getClanData(userData.Profile.ClanId, l.clanRepository)
//	if err != nil {
//		return nil, err
//	}
//	homePlanetPosition, homeSectorData, err := getHomeSectorData(userData, l.universeRepository)
//	if err != nil {
//		return nil, err
//	}
//
//	var response controllerModels.UserResponse
//	response.Profile.Init(*userData, clanData, l.userExperienceConstants)
//	homeSector, err := l.home(userData.OccupiedPlanets, *homePlanetPosition, homeSectorData)
//	if err != nil {
//		return nil, err
//	}
//	response.HomeSector = *homeSector
//	response.OccupiedPlanets, err = l.occupiedPlanets(userData.OccupiedPlanets, homePlanetPosition.SectorId(), homeSectorData)
//	if err != nil {
//		return nil, err
//	}
//	response.Notifications
//	return &response, nil
//}

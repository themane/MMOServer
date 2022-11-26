package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"hash/fnv"
	"strings"
)

type RegistrationService struct {
	userRepository          repoModels.UserRepository
	universeRepository      repoModels.UniverseRepository
	notificationService     *NotificationService
	userExperienceConstants constants.ExperienceConstants
	maxSystems              int
	retries                 int
	logger                  *constants.LoggingUtils
}

func NewRegistrationService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	maxSystems int,
	retries int,
	logLevel string,
) *RegistrationService {
	return &RegistrationService{
		userRepository:          userRepository,
		universeRepository:      universeRepository,
		notificationService:     NewNotificationService(experienceConstants, buildingConstants, mineConstants, militaryConstants, logLevel),
		userExperienceConstants: experienceConstants[constants.UserExperiences],
		maxSystems:              maxSystems,
		retries:                 retries,
		logger:                  constants.NewLoggingUtils("REGISTRATION_SERVICE", logLevel),
	}
}

func (r *RegistrationService) Register(request controllerModels.RegistrationRequest,
	id string, email string, name string) error {

	if r.UsernameExists(request.Username) {
		return errors.New("username already taken")
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	basePalnetUni, homePlanetUni, err := r.findNewPlanet(request.Location)
	if err != nil {
		return err
	}
	var mines []repoModels.MineUser
	for _, mine := range homePlanetUni.Mines {
		mines = append(mines, repoModels.MineUser{
			Id:    mine.Id,
			Mined: 0,
		})
	}
	userData := repoModels.UserData{
		Id: newUUID.String(),
		Profile: repoModels.ProfileUser{
			Username:   request.Username,
			Experience: 0,
			Species:    strings.ToUpper(request.Species),
			GoogleCredentials: repoModels.GoogleCredentials{
				Id:    id,
				Email: email,
				Name:  name,
			},
		},
		OccupiedPlanets: []repoModels.PlanetUser{
			{
				Id:         basePalnetUni.Id,
				BasePlanet: true,
			},
			{
				Id: homePlanetUni.Id,
				Water: repoModels.Resource{
					Amount:    0,
					Reserved:  0,
					Reserving: 0,
				},
				Graphene: repoModels.Resource{
					Amount:    0,
					Reserved:  0,
					Reserving: 0,
				},
				Shelio: 5,
				Population: repoModels.Population{
					GenerationRate: 1,
					Unemployed:     0,
					IdleWorkers:    0,
					IdleSoldiers:   0,
				},
				Mines:      mines,
				HomePlanet: true,
			},
		},
	}

	err = r.userRepository.AddUser(userData)
	if err != nil {
		return err
	}

	err = r.universeRepository.MarkOccupied(
		basePalnetUni.Position.System, basePalnetUni.Position.Sector, basePalnetUni.Position.Planet, userData.Id)
	if err != nil {
		return err
	}
	err = r.universeRepository.MarkOccupied(
		homePlanetUni.Position.System, homePlanetUni.Position.Sector, homePlanetUni.Position.Planet, userData.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RegistrationService) UsernameExists(username string) bool {
	user, err := r.userRepository.FindByUsername(username)
	return err == nil && user != nil && user.Profile.Username == username
}

func (r *RegistrationService) findNewPlanet(location string) (*repoModels.PlanetUni, *repoModels.PlanetUni, error) {
	system, err := r.hashPosition(location)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < r.retries; i++ {
		basePlanetUni, err1 := r.universeRepository.GetRandomUnoccupiedBasePlanet(system)
		if _, ok := err1.(*mongoRepository.NoSuchCombinationError); ok {
			continue
		}
		if err1 != nil {
			return nil, nil, err1
		}
		for j := 0; j < r.retries; j++ {
			homePlanetUni, err2 := r.universeRepository.GetRandomUnoccupiedHomePlanet(
				basePlanetUni.Position.System, basePlanetUni.Position.Sector, basePlanetUni.Position.Planet)
			if _, ok := err2.(*mongoRepository.NoSuchCombinationError); ok {
				continue
			}
			if err2 != nil {
				return nil, nil, err2
			}
			return basePlanetUni, homePlanetUni, nil
		}
	}
	return nil, nil, &mongoRepository.NoSuchCombinationError{}
}

func (r *RegistrationService) hashPosition(s string) (int, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return int(h.Sum32()) % r.maxSystems, nil
}

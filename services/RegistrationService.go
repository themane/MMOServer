package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/mongoRepository/exceptions"
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

func (r *RegistrationService) Register(request controllerModels.RegistrationRequest, userDetails models.UserSocialDetails) error {
	if r.UsernameExists(request.Username) {
		return errors.New("username already taken")
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	basePlanetUni, homePlanetUni, err := r.findNewPlanet(request.Location)
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
	profile := repoModels.ProfileUser{
		Username:   request.Username,
		Experience: 0,
		Species:    strings.ToUpper(request.Species),
	}
	if userDetails.Authenticator == constants.GoogleAuthenticator {
		profile.GoogleCredentials = userDetails
	}
	if userDetails.Authenticator == constants.FacebookAuthenticator {
		profile.FacebookCredentials = userDetails
	}

	userData := repoModels.UserData{
		Id:      newUUID.String(),
		Profile: profile,
		OccupiedPlanets: []repoModels.PlanetUser{
			{
				Id:         basePlanetUni.Id,
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
		basePlanetUni.Position.System, basePlanetUni.Position.Sector, basePlanetUni.Position.Planet, userData.Id)
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

func (r *RegistrationService) AddSocialLogin(username string, userDetails models.UserSocialDetails) error {
	user, err := r.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if userDetails.Authenticator == constants.GoogleAuthenticator {
		if len(user.Profile.GoogleCredentials.Id) == 0 {
			return r.userRepository.AddGoogleId(user.Id, userDetails)
		}
		return &exceptions.AlreadyExistsError{Message: "google sign in already linked with this account"}
	}
	if userDetails.Authenticator == constants.FacebookAuthenticator {
		if len(user.Profile.FacebookCredentials.Id) == 0 {
			return r.userRepository.AddFacebookId(user.Id, userDetails)
		}
		return &exceptions.AlreadyExistsError{Message: "facebook sign in already linked with this account"}
	}
	return errors.New("cannot register this login")
}

func (r *RegistrationService) findNewPlanet(location string) (*repoModels.PlanetUni, *repoModels.PlanetUni, error) {
	system, err := r.hashPosition(location)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < r.retries; i++ {
		basePlanetUni, err1 := r.universeRepository.GetRandomUnoccupiedBasePlanet(system)
		if _, ok := err1.(*exceptions.NoSuchCombinationError); ok {
			continue
		}
		if err1 != nil {
			return nil, nil, err1
		}
		for j := 0; j < r.retries; j++ {
			homePlanetUni, err2 := r.universeRepository.GetRandomUnoccupiedHomePlanet(
				basePlanetUni.Position.System, basePlanetUni.Position.Sector, basePlanetUni.Position.Planet)
			if _, ok := err2.(*exceptions.NoSuchCombinationError); ok {
				continue
			}
			if err2 != nil {
				return nil, nil, err2
			}
			return basePlanetUni, homePlanetUni, nil
		}
	}
	return nil, nil, &exceptions.NoSuchCombinationError{}
}

func (r *RegistrationService) hashPosition(s string) (int, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return int(h.Sum32()) % r.maxSystems, nil
}

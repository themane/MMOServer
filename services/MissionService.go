package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type MissionService struct {
	userRepository     repoModels.UserRepository
	universeRepository repoModels.UniverseRepository
	missionRepository  repoModels.MissionRepository
	shipConstants      map[string]constants.ShipConstants
	logger             *constants.LoggingUtils
}

func NewMissionService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *MissionService {
	return &MissionService{
		userRepository:     userRepository,
		universeRepository: universeRepository,
		missionRepository:  missionRepository,
		shipConstants:      shipConstants,
		logger:             constants.NewLoggingUtils("MISSION_SERVICE", logLevel),
	}
}

func (a *MissionService) Spy(spyRequest controllerModels.SpyRequest) error {
	userData, err := a.userRepository.FindByUsername(spyRequest.Username)
	if err != nil {
		return err
	}
	fromPlanetUni, err := a.universeRepository.FindById(spyRequest.FromPlanetId)
	if err != nil {
		return err
	}
	toPlanetUni, err := a.universeRepository.FindById(spyRequest.ToPlanetId)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[spyRequest.FromPlanetId]; ok {
		var squadSpeed float64
		availableShips := planetUser.GetAvailableShips()
		for _, formation := range spyRequest.Scouts {
			if availableShips[formation.ShipName] < formation.Quantity {
				return errors.New("error! found insufficient ships for attack formation")
			}
			availableShips[formation.ShipName] -= formation.Quantity
			currentLevel := strconv.Itoa(planetUser.Ships[formation.ShipName].Level)
			speed := a.shipConstants[formation.ShipName].Levels[currentLevel].Speed
			if squadSpeed < float64(speed) {
				squadSpeed = float64(speed)
			}
		}

		blocks := distance(*fromPlanetUni, *toPlanetUni)
		totalSecondsRequired := blocks * squadSpeed
		missionTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired))
		returnTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired) * 2)
		spyMission, err := spyRequest.GetSpyMission(primitive.NewDateTimeFromTime(missionTime), primitive.NewDateTimeFromTime(returnTime))
		if err != nil {
			return err
		}
		err = a.missionRepository.AddSpyMission(*spyMission)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("error occurred in retrieving planet data")
}

func (a *MissionService) Attack(attackRequest controllerModels.AttackRequest) error {
	userData, err := a.userRepository.FindByUsername(attackRequest.Username)
	if err != nil {
		return err
	}
	fromPlanetUni, err := a.universeRepository.FindById(attackRequest.FromPlanetId)
	if err != nil {
		return err
	}
	toPlanetUni, err := a.universeRepository.FindById(attackRequest.ToPlanetId)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[attackRequest.FromPlanetId]; ok {
		var squadSpeed float64
		availableShips := planetUser.GetAvailableShips()
		for attackPointId, formationMap := range attackRequest.Formation {
			if !validAttackPointId(attackPointId) {
				return errors.New("error! found invalid attack point id: " + attackPointId)
			}
			err := validateAttackLineIds(formationMap, attackPointId)
			if err != nil {
				return err
			}
			for _, formations := range formationMap {
				for _, formation := range formations {
					if availableShips[formation.ShipName] < formation.Quantity {
						return errors.New("error! found insufficient ships for attack formation")
					}
					availableShips[formation.ShipName] -= formation.Quantity
					currentLevel := strconv.Itoa(planetUser.Ships[formation.ShipName].Level)
					speed := a.shipConstants[formation.ShipName].Levels[currentLevel].Speed
					if squadSpeed < float64(speed) {
						squadSpeed = float64(speed)
					}
				}
			}
		}
		blocks := distance(*fromPlanetUni, *toPlanetUni)
		totalSecondsRequired := blocks * squadSpeed
		missionTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired))
		returnTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired) * 2)
		attackMission, err := attackRequest.GetAttackMission(primitive.NewDateTimeFromTime(missionTime), primitive.NewDateTimeFromTime(returnTime))
		if err != nil {
			return err
		}
		err = a.missionRepository.AddAttackMission(*attackMission)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("error occurred in retrieving planet data")
}

func distance(fromPlanet repoModels.PlanetUni, toPlanet repoModels.PlanetUni) float64 {
	if math.Abs(float64(fromPlanet.Position.System-toPlanet.Position.System)) > 0 {
		return constants.SystemDistanceBlocks
	}
	sectorDifference := math.Abs(float64(fromPlanet.Position.Sector - toPlanet.Position.Sector))
	if sectorDifference > 0 {
		return sectorDifference * constants.SectorDistanceBlocks
	}
	return math.Abs(float64(fromPlanet.Distance - toPlanet.Distance))
}

func validateAttackLineIds(formationMap map[string][]models.Formation, attackPointId string) error {
	var lineIds []string
	for key := range formationMap {
		lineIds = append(lineIds, key)
	}
	sort.Strings(lineIds)
	validAttackLineIds := constants.GetAttackLineIds()
	if !reflect.DeepEqual(lineIds, validAttackLineIds) {
		return errors.New("line ids not valid for point id: " + attackPointId)
	}
	return nil
}

func validAttackPointId(attackPointId string) bool {
	for _, validAttackPointId := range constants.GetAttackPointIds() {
		if validAttackPointId == attackPointId {
			return true
		}
	}
	return false
}

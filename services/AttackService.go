package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/schedulers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"reflect"
	"strconv"
	"time"
)

type AttackService struct {
	userRepository          repoModels.UserRepository
	universeRepository      repoModels.UniverseRepository
	missionRepository       repoModels.MissionRepository
	scheduledMissionManager schedulers.ScheduledMissionManager
	shipConstants           map[string]constants.ShipConstants
}

func NewAttackService(
	userRepository repoModels.UserRepository,
	universeRepository repoModels.UniverseRepository,
	missionRepository repoModels.MissionRepository,
	scheduledMissionManager schedulers.ScheduledMissionManager,
	shipConstants map[string]constants.ShipConstants,
) *AttackService {
	return &AttackService{
		userRepository:          userRepository,
		universeRepository:      universeRepository,
		missionRepository:       missionRepository,
		scheduledMissionManager: scheduledMissionManager,
		shipConstants:           shipConstants,
	}
}

func (a *AttackService) Spy(spyRequest controllerModels.SpyRequest) (*controllerModels.MissionResponse, error) {
	var squadSpeed float64
	userData, err := a.userRepository.FindByUsername(spyRequest.Attacker)
	if err != nil {
		return nil, err
	}
	fromPlanetUni, err := a.universeRepository.FindById(spyRequest.FromPlanetId)
	if err != nil {
		return nil, err
	}
	toPlanetUni, err := a.universeRepository.FindById(spyRequest.ToPlanetId)
	if err != nil {
		return nil, err
	}
	if planetUser, ok := userData.OccupiedPlanets[spyRequest.FromPlanetId]; ok {
		availableShips := planetUser.GetAvailableShips()
		scoutMap := map[string]int{}
		for _, formation := range spyRequest.Scouts {
			if availableShips[formation.ShipName] < formation.Quantity {
				return nil, errors.New("error! found insufficient ships for attack formation")
			}
			availableShips[formation.ShipName] -= formation.Quantity
			currentLevel := strconv.Itoa(planetUser.Ships[formation.ShipName].Level)
			speed := a.shipConstants[formation.ShipName].Levels[currentLevel].Speed
			if squadSpeed < float64(speed) {
				squadSpeed = float64(speed)
			}
			scoutMap[formation.ShipName] = formation.Quantity
		}

		blocks := distance(*fromPlanetUni, *toPlanetUni)
		totalSecondsRequired := blocks * squadSpeed
		missionTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired))
		returnTime := time.Now().Add(time.Second * time.Duration(totalSecondsRequired) * 2)

		spyMission, err := a.missionRepository.AddSpyMission(spyRequest.FromPlanetId, spyRequest.ToPlanetId, scoutMap,
			primitive.NewDateTimeFromTime(time.Now()), primitive.NewDateTimeFromTime(missionTime), primitive.NewDateTimeFromTime(returnTime),
		)
		if err != nil {
			return nil, err
		}
		a.scheduledMissionManager.ScheduleSpyMission(*spyMission, missionTime, returnTime)

		response := controllerModels.MissionResponse{MissionTime: missionTime.String(), ReturnTime: returnTime.String()}
		return &response, nil
	}
	return nil, errors.New("error occurred in retrieving planet data")
}

func (a *AttackService) Attack(attackRequest controllerModels.AttackRequest) (*controllerModels.MissionResponse, error) {
	var squadSpeed float64
	userData, err := a.userRepository.FindByUsername(attackRequest.Attacker)
	if err != nil {
		return nil, err
	}
	fromPlanetUni, err := a.universeRepository.FindById(attackRequest.FromPlanetId)
	if err != nil {
		return nil, err
	}
	toPlanetUni, err := a.universeRepository.FindById(attackRequest.ToPlanetId)
	if err != nil {
		return nil, err
	}
	if planetUser, ok := userData.OccupiedPlanets[attackRequest.FromPlanetId]; ok {
		availableShips := planetUser.GetAvailableShips()
		for attackPointId, formationMap := range attackRequest.Formation {
			if !validAttackPointId(attackPointId) {
				return nil, errors.New("error! found invalid attack point id: " + attackPointId)
			}
			err := validateAttackLineIds(formationMap, attackPointId)
			if err != nil {
				return nil, err
			}
			for _, formations := range formationMap {
				for _, formation := range formations {
					if availableShips[formation.ShipName] < formation.Quantity {
						return nil, errors.New("error! found insufficient ships for attack formation")
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

		attackMission, err := a.missionRepository.AddAttackMission(attackRequest.FromPlanetId, attackRequest.ToPlanetId, attackRequest.Formation,
			primitive.NewDateTimeFromTime(time.Now()), primitive.NewDateTimeFromTime(missionTime), primitive.NewDateTimeFromTime(returnTime),
		)
		if err != nil {
			return nil, err
		}
		a.scheduledMissionManager.ScheduleAttackMission(*attackMission, missionTime, returnTime)

		response := controllerModels.MissionResponse{MissionTime: missionTime.String(), ReturnTime: returnTime.String()}
		return &response, nil
	}
	return nil, errors.New("error occurred in retrieving planet data")
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

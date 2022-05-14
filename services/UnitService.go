package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
	"time"
)

type UnitService struct {
	userRepository    repoModels.UserRepository
	missionRepository repoModels.MissionRepository
	militaryConstants map[string]constants.MilitaryConstants
	logger            *constants.LoggingUtils
}

func NewUnitService(
	userRepository repoModels.UserRepository,
	missionRepository repoModels.MissionRepository,
	militaryConstants map[string]constants.MilitaryConstants,
	logLevel string,
) *UnitService {
	return &UnitService{
		userRepository:    userRepository,
		missionRepository: missionRepository,
		militaryConstants: militaryConstants,
		logger:            constants.NewLoggingUtils("UNIT_SERVICE", logLevel),
	}
}

func (u *UnitService) ConstructUnits(username string, planetId string, unitName string, quantity int) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if unitMilitaryConstants, ok := u.militaryConstants[unitName]; ok {
			if unitMilitaryConstants.Type == constants.Defender {
				if planetUser.Defences[unitName].Level <= 0 {
					return errors.New("unit not available for construction")
				}
				unitLevelString := strconv.Itoa(planetUser.Defences[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				requirements, err := u.validateAndGetRequirements(quantity, planetUser, unitLevelConstants)
				if err != nil {
					return err
				}
				if planetUser.Defences[unitName].UnderConstruction.Quantity > 0 {
					err = u.userRepository.AddConstructionDefences(userData.Id, planetId, unitName, float64(quantity), *requirements)
					if err != nil {
						return err
					}
					return nil
				}
				err = u.userRepository.ConstructDefences(userData.Id, planetId, unitName, float64(quantity), *requirements)
				if err != nil {
					return err
				}
				return nil
			} else if unitMilitaryConstants.Type == constants.DefenceShipCarrier {
				if quantity != 1 {
					return errors.New("only single construction of DefenceShipCarrier supported")
				}
				unitLevelConstants := unitMilitaryConstants.Levels["1"]
				requirements, err := u.validateAndGetRequirements(quantity, planetUser, unitLevelConstants)
				if err != nil {
					return err
				}
				id, err := uuid.NewRandom()
				if err != nil {
					return errors.New("error in generating DefenceShipCarrier id")
				}
				err = u.userRepository.ConstructDefenceShipCarrier(userData.Id, planetId, unitName, id.String(), *requirements)
				if err != nil {
					return err
				}
				return nil
			} else {
				if planetUser.Ships[unitName].Level <= 0 {
					return errors.New("unit not available for construction")
				}
				unitLevelString := strconv.Itoa(planetUser.Ships[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				requirements, err := u.validateAndGetRequirements(quantity, planetUser, unitLevelConstants)
				if err != nil {
					return err
				}
				if planetUser.Ships[unitName].UnderConstruction.Quantity > 0 {
					err = u.userRepository.AddConstructionShips(userData.Id, planetId, unitName, float64(quantity), *requirements)
					if err != nil {
						return err
					}
					return nil
				}
				err = u.userRepository.ConstructShips(userData.Id, planetId, unitName, float64(quantity), *requirements)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) CancelUnitsConstruction(username string, planetId string, unitName string) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if unitMilitaryConstants, ok := u.militaryConstants[unitName]; ok {
			if unitMilitaryConstants.Type == constants.Defender {
				if planetUser.Defences[unitName].UnderConstruction.Quantity <= 0 {
					return errors.New("no units under construction")
				}
				unitLevelString := strconv.Itoa(planetUser.Defences[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				cancelReturns := models.Returns{}
				cancelReturns.InitCancelReturns(unitLevelConstants, float64(planetUser.Defences[unitName].UnderConstruction.Quantity))
				err = u.userRepository.CancelDefencesConstruction(userData.Id, planetId, unitName, cancelReturns)
				if err != nil {
					return err
				}
			} else if unitMilitaryConstants.Type == constants.DefenceShipCarrier {
				return errors.New("unit not valid")
			} else {
				if planetUser.Ships[unitName].UnderConstruction.Quantity <= 0 {
					return errors.New("no units under construction")
				}
				unitLevelString := strconv.Itoa(planetUser.Ships[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				cancelReturns := models.Returns{}
				cancelReturns.InitCancelReturns(unitLevelConstants, float64(planetUser.Ships[unitName].UnderConstruction.Quantity))
				err = u.userRepository.CancelShipsConstruction(userData.Id, planetId, unitName, cancelReturns)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) DestructUnits(username string, planetId string, unitName string, quantity int) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if unitMilitaryConstants, ok := u.militaryConstants[unitName]; ok {
			if unitMilitaryConstants.Type == constants.Defender {
				if repoModels.GetIdleDefences(planetUser.Defences[unitName].GuardingShield, planetUser.Defences[unitName].Quantity) < quantity {
					return errors.New("not enough idle units")
				}
				unitLevelString := strconv.Itoa(planetUser.Defences[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				destructionReturns := models.Returns{}
				destructionReturns.InitDestructionReturns(unitLevelConstants)
				err = u.userRepository.DestructDefences(userData.Id, planetId, unitName, float64(quantity), destructionReturns)
				if err != nil {
					return err
				}
			} else if unitMilitaryConstants.Type == constants.DefenceShipCarrier {
				return errors.New("unit not valid")
			} else if unitMilitaryConstants.Type == constants.Scout {
				spyMissions, err1 := u.missionRepository.FindSpyMissionsFromPlanetId(planetId)
				if err1 != nil {
					return err1
				}
				if repoModels.GetAvailableScouts(unitName, spyMissions, planetUser.Ships[unitName].Quantity) < quantity {
					return errors.New("not enough idle units")
				}
				unitLevelString := strconv.Itoa(planetUser.Defences[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				destructionReturns := models.Returns{}
				destructionReturns.InitDestructionReturns(unitLevelConstants)
				err = u.userRepository.DestructShips(userData.Id, planetId, unitName, float64(quantity), destructionReturns)
				if err != nil {
					return err
				}
			} else {
				attackMissions, err2 := u.missionRepository.FindAttackMissionsFromPlanetId(planetId)
				if err2 != nil {
					return err2
				}
				if repoModels.GetAvailableShips(unitName, attackMissions, planetUser.DefenceShipCarriers, planetUser.Ships[unitName].Quantity) < quantity {
					return errors.New("not enough idle units")
				}
				unitLevelString := strconv.Itoa(planetUser.Defences[unitName].Level)
				unitLevelConstants := unitMilitaryConstants.Levels[unitLevelString]
				destructionReturns := models.Returns{}
				destructionReturns.InitDestructionReturns(unitLevelConstants)
				err = u.userRepository.DestructShips(userData.Id, planetId, unitName, float64(quantity), destructionReturns)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) UpgradeDefenceShipCarrier(username string, planetId string, unitId string) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if defenceShipCarrier, ok := planetUser.DefenceShipCarriers[unitId]; ok {
			if defenceShipCarrier.GuardingShield != "" {
				return errors.New("not idle unit")
			}
			if defenceShipCarrier.UnderConstruction.StartTime.Time().After(time.Now()) {
				return errors.New("unit under construction")
			}
			for _, q := range defenceShipCarrier.HostingShips {
				if q > 0 {
					return errors.New("unit has deployed ships")
				}
			}
			requirements, err1 := u.validateAndGetNextLevelRequirements(defenceShipCarrier, planetUser)
			if err1 != nil {
				return err1
			}
			err = u.userRepository.UpgradeDefenceShipCarrier(userData.Id, planetId, unitId, *requirements)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) CancelDefenceShipCarrierUpGradation(username string, planetId string, unitId string) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if defenceShipCarrier, ok := planetUser.DefenceShipCarriers[unitId]; ok {
			if !defenceShipCarrier.UnderConstruction.StartTime.Time().After(time.Now()) {
				return errors.New("unit not under construction")
			}
			unitLevelConstants := u.militaryConstants[defenceShipCarrier.Name].Levels[strconv.Itoa(defenceShipCarrier.Level)]
			cancelReturns := models.Returns{}
			cancelReturns.InitCancelReturns(unitLevelConstants, 1)
			if defenceShipCarrier.Level > 1 {
				err = u.userRepository.CancelDefenceShipCarrierUpGradation(userData.Id, planetId, unitId, cancelReturns)
				if err != nil {
					return err
				}
			} else {
				err = u.userRepository.CancelDefenceShipCarrierConstruction(userData.Id, planetId, unitId, cancelReturns)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) DestructDefenceShipCarrier(username string, planetId string, unitId string) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if defenceShipCarrier, ok := planetUser.DefenceShipCarriers[unitId]; ok {
			if defenceShipCarrier.GuardingShield != "" {
				return errors.New("not idle unit")
			}
			if defenceShipCarrier.UnderConstruction.StartTime.Time().After(time.Now()) {
				return errors.New("unit under construction")
			}
			for _, q := range defenceShipCarrier.HostingShips {
				if q > 0 {
					return errors.New("unit has deployed ships")
				}
			}
			unitLevelConstants := u.militaryConstants[defenceShipCarrier.Name].Levels[strconv.Itoa(defenceShipCarrier.Level)]
			returns := models.Returns{}
			returns.InitDestructionReturns(unitLevelConstants)
			err = u.userRepository.DestructDefenceShipCarrier(userData.Id, planetId, unitId, returns)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitService) validateAndGetRequirements(quantity int, planetUser repoModels.PlanetUser,
	unitLevelConstants map[string]interface{}) (*models.Requirements, error) {

	creationRequirements := models.Requirements{}
	creationRequirements.Init(unitLevelConstants)
	if creationRequirements.Population.Soldiers*float64(quantity) > float64(planetUser.Population.IdleSoldiers) ||
		creationRequirements.Population.Workers*float64(quantity) > float64(planetUser.Population.IdleWorkers) ||
		creationRequirements.Resources.Graphene*float64(quantity) > float64(planetUser.Graphene.Amount) ||
		creationRequirements.Resources.Water*float64(quantity) > float64(planetUser.Water.Amount) ||
		creationRequirements.Resources.Shelio*float64(quantity) > float64(planetUser.Shelio) {
		return nil, errors.New("not enough resources")
	}
	return &creationRequirements, nil
}

func (u *UnitService) validateAndGetNextLevelRequirements(defenceShipCarrier repoModels.DefenceShipCarrier, planetUser repoModels.PlanetUser) (*models.Requirements, error) {
	nextLevelRequirements := models.Requirements{}
	nextLevelRequirements.InitNextLevelRequirements(defenceShipCarrier.Level, u.militaryConstants[defenceShipCarrier.Name])
	if nextLevelRequirements.Population.Soldiers > float64(planetUser.Population.IdleSoldiers) ||
		nextLevelRequirements.Population.Workers > float64(planetUser.Population.IdleWorkers) ||
		nextLevelRequirements.Resources.Graphene > float64(planetUser.Graphene.Amount) ||
		nextLevelRequirements.Resources.Water > float64(planetUser.Water.Amount) ||
		nextLevelRequirements.Resources.Shelio > float64(planetUser.Shelio) {
		return nil, errors.New("not enough resources")
	}
	return &nextLevelRequirements, nil
}

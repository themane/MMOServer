package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type UnitsDeploymentService struct {
	userRepository    repoModels.UserRepository
	missionRepository repoModels.MissionRepository
	militaryConstants map[string]constants.MilitaryConstants
	logger            *constants.LoggingUtils
}

func NewUnitsDeploymentService(
	userRepository repoModels.UserRepository,
	missionRepository repoModels.MissionRepository,
	militaryConstants map[string]constants.MilitaryConstants,
	logLevel string,
) *UnitsDeploymentService {
	return &UnitsDeploymentService{
		userRepository:    userRepository,
		missionRepository: missionRepository,
		militaryConstants: militaryConstants,
		logger:            constants.NewLoggingUtils("UNIT_SERVICE", logLevel),
	}
}

func (u *UnitsDeploymentService) DeployShipsOnDefenceShipCarrier(username string, planetId string, unitId string, ships map[string]int) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if defenceShipCarrier, ok := planetUser.DefenceShipCarriers[unitId]; ok {
			err = u.validateShipsChanges(ships, planetId, defenceShipCarrier, planetUser)
			if err != nil {
				return err
			}
			err = u.userRepository.DeployShipsOnDefenceShipCarrier(userData.Id, planetId, unitId, ships)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

func (u *UnitsDeploymentService) validateShipsChanges(ships map[string]int, planetId string,
	defenceShipCarrier repoModels.DefenceShipCarrier, planetUser repoModels.PlanetUser) error {

	for shipName, quantity := range ships {
		if u.militaryConstants[shipName].Type != constants.Attacker {
			return errors.New("only attackers are allowed to be deployed")
		}
		attackMissions, err := u.missionRepository.FindAttackMissionsFromPlanetId(planetId)
		if err != nil {
			return err
		}
		if repoModels.GetAvailableShips(shipName, attackMissions, planetUser.DefenceShipCarriers, planetUser.Ships[shipName].Quantity)+defenceShipCarrier.HostingShips[shipName] < quantity {
			return errors.New("not enough available ships")
		}
	}
	return nil
}

func (u *UnitsDeploymentService) DeployDefencesOnShield(username string, planetId string, shieldId string, defences map[string]int) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		err = u.validateDefenceChanges(defences, shieldId, planetUser)
		if err != nil {
			return err
		}
		err = u.userRepository.DeployDefencesOnShield(userData.Id, planetId, shieldId, defences)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("planet not occupied")
}

func (u *UnitsDeploymentService) validateDefenceChanges(defences map[string]int, shieldId string,
	planetUser repoModels.PlanetUser) error {

	for defenceName, quantity := range defences {
		if u.militaryConstants[defenceName].Type != constants.Defender {
			return errors.New("not a defender")
		}
		defence := planetUser.Defences[defenceName]
		if repoModels.GetIdleDefences(defence.GuardingShield, defence.Quantity)+defence.GuardingShield[shieldId] < quantity {
			return errors.New("not enough idle defences")
		}
	}
	return nil
}

func (u *UnitsDeploymentService) DeployDefenceShipCarrierOnShield(username string, planetId string, shieldId string, unitId string, deploy bool) error {
	userData, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if planetUser, ok := userData.OccupiedPlanets[planetId]; ok {
		if defenceShipCarrier, ok := planetUser.DefenceShipCarriers[unitId]; ok {
			if deploy == false {
				if defenceShipCarrier.GuardingShield != shieldId {
					return errors.New("already not guarding")
				}
				shieldId = ""
			}
			err = u.userRepository.DeployDefenceShipCarrierOnShield(userData.Id, planetId, unitId, shieldId)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("unit not valid")
	}
	return errors.New("planet not occupied")
}

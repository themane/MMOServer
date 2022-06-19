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
	occupiedPlanet := userData.GetOccupiedPlanet(planetId)
	if occupiedPlanet == nil {
		return errors.New("planet not occupied")
	}
	defenceShipCarrier := occupiedPlanet.GetDefenceShipCarrier(unitId)
	if defenceShipCarrier == nil {
		return errors.New("defence ship carrier id not valid")
	}
	err = u.validateShipsChanges(ships, planetId, *defenceShipCarrier, *occupiedPlanet)
	if err != nil {
		return err
	}
	err = u.userRepository.DeployShipsOnDefenceShipCarrier(userData.Id, planetId, unitId, ships)
	if err != nil {
		return err
	}
	return nil
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
		availableShips := repoModels.GetAvailableShips(shipName, attackMissions, planetUser.DefenceShipCarriers, planetUser.GetShip(shipName).Quantity)
		if availableShips+defenceShipCarrier.HostingShips[shipName] <= quantity {
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
	occupiedPlanet := userData.GetOccupiedPlanet(planetId)
	if occupiedPlanet == nil {
		return errors.New("planet not occupied")
	}
	err = u.validateDefenceChanges(defences, shieldId, *occupiedPlanet)
	if err != nil {
		return err
	}
	err = u.userRepository.DeployDefencesOnShield(userData.Id, planetId, shieldId, defences)
	if err != nil {
		return err
	}
	return nil
}

func (u *UnitsDeploymentService) validateDefenceChanges(defences map[string]int, shieldId string,
	planetUser repoModels.PlanetUser) error {

	for defenceName, quantity := range defences {
		if u.militaryConstants[defenceName].Type != constants.Defender {
			return errors.New("not a defender")
		}
		defence := planetUser.GetDefence(defenceName)
		idleDefences := repoModels.GetIdleDefences(defence.GuardingShield, defence.Quantity)
		if idleDefences+defence.GuardingShield[shieldId] <= quantity {
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
	occupiedPlanet := userData.GetOccupiedPlanet(planetId)
	if occupiedPlanet == nil {
		return errors.New("planet not occupied")
	}
	defenceShipCarrier := occupiedPlanet.GetDefenceShipCarrier(unitId)
	if defenceShipCarrier == nil {
		return errors.New("defence ship carrier id not valid")
	}
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

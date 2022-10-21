package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SpyRequest struct {
	FromPlanetId string             `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string             `json:"to_planet_id" example:"001:002:05"`
	Scouts       []models.Formation `json:"scouts"`
}

func (s *SpyRequest) GetSpyMission(missionTime primitive.DateTime, returnTime primitive.DateTime) (*repoModels.SpyMission, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("error in generating mission id")
	}
	spyMission := repoModels.SpyMission{
		Id:           id.String(),
		FromPlanetId: s.FromPlanetId,
		ToPlanetId:   s.ToPlanetId,
		LaunchTime:   primitive.NewDateTimeFromTime(time.Now()),
		MissionTime:  missionTime,
		ReturnTime:   returnTime,
		State:        constants.DepartureState,
		MissionType:  constants.SpyMission,
	}
	for _, formation := range s.Scouts {
		spyMission.Scouts[formation.ShipName] = formation.Quantity
	}
	return &spyMission, nil
}

type AttackRequest struct {
	FromPlanetId string                                   `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string                                   `json:"to_planet_id" example:"001:002:05"`
	Formation    map[string]map[string][]models.Formation `json:"formation"`
}

func (a *AttackRequest) GetAttackMission(missionTime primitive.DateTime, returnTime primitive.DateTime) (*repoModels.AttackMission, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("error in generating mission id")
	}
	attackMission := repoModels.AttackMission{
		Id:           id.String(),
		FromPlanetId: a.FromPlanetId,
		ToPlanetId:   a.ToPlanetId,
		LaunchTime:   primitive.NewDateTimeFromTime(time.Now()),
		MissionTime:  missionTime,
		ReturnTime:   returnTime,
		State:        constants.DepartureState,
		MissionType:  constants.AttackMission,
	}
	for pointId, pointFormation := range a.Formation {
		for lineId, lineFormation := range pointFormation {
			for _, formation := range lineFormation {
				attackMission.Formation[pointId][lineId][formation.ShipName] = formation.Quantity
			}
		}
	}
	return &attackMission, nil
}

type ActiveMission struct {
	Id          string                                   `json:"_id"`
	ToPlanetId  string                                   `json:"to_planet_id"`
	Formation   map[string]map[string][]models.Formation `json:"formation,omitempty"`
	Scouts      []models.Formation                       `json:"scouts,omitempty"`
	LaunchTime  time.Time                                `json:"launch_time" bson:"launch_time"`
	MissionTime time.Time                                `json:"mission_time" bson:"mission_time"`
	ReturnTime  time.Time                                `json:"return_time" bson:"return_time"`
	MissionType string                                   `json:"mission_type" bson:"mission_type"`
}

func (a *ActiveMission) InitAttackMission(missionData repoModels.AttackMission) {
	a.Id = missionData.Id
	a.ToPlanetId = missionData.ToPlanetId
	for pointId, pointFormation := range missionData.Formation {
		for lineId, lineFormation := range pointFormation {
			for unitName, units := range lineFormation {
				a.Formation[pointId][lineId] = append(a.Formation[pointId][lineId], models.Formation{ShipName: unitName, Quantity: units})
			}
		}
	}
	a.LaunchTime = missionData.LaunchTime.Time()
	a.MissionTime = missionData.MissionTime.Time()
	a.ReturnTime = missionData.ReturnTime.Time()
	a.MissionType = missionData.MissionType
}

func (a *ActiveMission) InitSpyMission(missionData repoModels.SpyMission) {
	a.Id = missionData.Id
	a.ToPlanetId = missionData.ToPlanetId
	for unitName, units := range missionData.Scouts {
		a.Scouts = append(a.Scouts, models.Formation{ShipName: unitName, Quantity: units})
	}
	a.LaunchTime = missionData.LaunchTime.Time()
	a.MissionTime = missionData.MissionTime.Time()
	a.ReturnTime = missionData.ReturnTime.Time()
	a.MissionType = missionData.MissionType
}

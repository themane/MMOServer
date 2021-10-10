package models

import (
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"time"
)

type MissionResponse struct {
	MissionTime string `json:"mission_time"`
	ReturnTime  string `json:"return_time"`
}

type SpyRequest struct {
	Attacker     string             `json:"attacker" example:"devashish"`
	FromPlanetId string             `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string             `json:"to_planet_id" example:"001:002:05"`
	Scouts       []models.Formation `json:"scouts"`
}

type AttackRequest struct {
	Attacker     string                                   `json:"attacker" example:"devashish"`
	FromPlanetId string                                   `json:"from_planet_id" example:"001:002:03"`
	ToPlanetId   string                                   `json:"to_planet_id" example:"001:002:05"`
	Formation    map[string]map[string][]models.Formation `json:"formation"`
}

type ActiveMission struct {
	Id          string                                   `json:"_id"`
	ToPlanetId  string                                   `json:"to_planet_id"`
	Formation   map[string]map[string][]models.Formation `json:"formation"`
	Scouts      map[string]int                           `json:"scouts"`
	LaunchTime  time.Time                                `json:"launch_time" bson:"launch_time"`
	MissionTime time.Time                                `json:"mission_time" bson:"mission_time"`
	ReturnTime  time.Time                                `json:"return_time" bson:"return_time"`
	MissionType string                                   `json:"mission_type" bson:"mission_type"`
}

func (a *ActiveMission) InitAttackMission(missionData repoModels.AttackMission) {
	a.Id = missionData.Id
	a.ToPlanetId = missionData.ToPlanetId
	a.Formation = missionData.Formation
	a.LaunchTime = missionData.LaunchTime
	a.MissionTime = missionData.MissionTime
	a.ReturnTime = missionData.ReturnTime
	a.MissionType = missionData.MissionType
}

func (a *ActiveMission) InitSpyMission(missionData repoModels.SpyMission) {
	a.Id = missionData.Id
	a.ToPlanetId = missionData.ToPlanetId
	a.Scouts = missionData.Scouts
	a.LaunchTime = missionData.LaunchTime
	a.MissionTime = missionData.MissionTime
	a.ReturnTime = missionData.ReturnTime
	a.MissionType = missionData.MissionType
}

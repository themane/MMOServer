package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttackMission struct {
	Id           string                               `json:"_id" bson:"_id"`
	FromPlanetId string                               `json:"from_planet_id" bson:"from_planet_id"`
	ToPlanetId   string                               `json:"to_planet_id" bson:"to_planet_id"`
	Formation    map[string]map[string]map[string]int `json:"formation"`
	Result       AttackResult                         `json:"result" bson:"result"`
	LaunchTime   primitive.DateTime                   `json:"launch_time" bson:"launch_time"`
	MissionTime  primitive.DateTime                   `json:"mission_time" bson:"mission_time"`
	ReturnTime   primitive.DateTime                   `json:"return_time" bson:"return_time"`
	State        string                               `json:"state" bson:"state"`
	MissionType  string                               `json:"mission_type" bson:"mission_type"`
}

type AttackResult struct {
}

type SpyMission struct {
	Id           string             `json:"_id" bson:"_id"`
	FromPlanetId string             `json:"from_planet_id" bson:"from_planet_id"`
	ToPlanetId   string             `json:"to_planet_id" bson:"to_planet_id"`
	Scouts       map[string]int     `json:"scouts" bson:"scouts"`
	Result       SpyResult          `json:"result" bson:"result"`
	LaunchTime   primitive.DateTime `json:"launch_time" bson:"launch_time"`
	MissionTime  primitive.DateTime `json:"mission_time" bson:"mission_time"`
	ReturnTime   primitive.DateTime `json:"return_time" bson:"return_time"`
	State        string             `json:"state" bson:"state"`
	MissionType  string             `json:"mission_type" bson:"mission_type"`
}

type SpyResult struct {
}

type MissionRepository interface {
	FindAttackMissionsFromPlanetId(id string) ([]AttackMission, error)
	FindSpyMissionsFromPlanetId(id string) ([]SpyMission, error)
	FindAttackMissionsToPlanetId(id string) ([]AttackMission, error)
	FindSpyMissionsToPlanetId(id string) ([]SpyMission, error)

	AddAttackMission(mission AttackMission) error
	AddSpyMission(mission SpyMission) error

	UpdateAttackResult(id string, result AttackResult) error
	UpdateSpyResult(id string, result SpyResult) error
	UpdateMissionState(id string, state string) error
}

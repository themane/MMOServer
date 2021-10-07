package models

import (
	"github.com/themane/MMOServer/models"
)

type AttackMission struct {
	Id           string                                   `json:"_id" bson:"_id"`
	FromPlanetId string                                   `json:"from_planet_id" bson:"from_planet_id"`
	ToPlanetId   string                                   `json:"to_planet_id" bson:"to_planet_id"`
	Formation    map[string]map[string][]models.Formation `json:"formation"`
	Result       AttackResult                             `json:"result" bson:"result"`
	State        string                                   `json:"state" bson:"state"`
	MissionType  string                                   `json:"mission_type" bson:"mission_type"`
}

type AttackResult struct {
}

type SpyMission struct {
	Id           string         `json:"_id" bson:"_id"`
	FromPlanetId string         `json:"from_planet_id" bson:"from_planet_id"`
	ToPlanetId   string         `json:"to_planet_id" bson:"to_planet_id"`
	Scouts       map[string]int `json:"scouts" bson:"scouts"`
	Result       SpyResult      `json:"result" bson:"result"`
	State        string         `json:"state" bson:"state"`
	MissionType  string         `json:"mission_type" bson:"mission_type"`
}

type SpyResult struct {
}

type MissionRepository interface {
	FindAttackMissionsFromPlanetId(id string) ([]AttackMission, error)
	FindSpyMissionsFromPlanetId(id string) ([]SpyMission, error)
	FindAttackMissionsToPlanetId(id string) ([]AttackMission, error)
	FindSpyMissionsToPlanetId(id string) ([]SpyMission, error)

	AddAttackMission(fromPlanetId string, toPlanetId string, formation map[string]map[string][]models.Formation) error
	AddSpyMission(fromPlanetId string, toPlanetId string, scouts map[string]int) error

	UpdateAttackResult(id string, result AttackResult) error
	UpdateSpyResult(id string, result SpyResult) error
	UpdateMissionState(id string, state string) error
}

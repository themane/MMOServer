package constants

import (
	"errors"
	"strings"
)

//  Shield Types
const (
	Invulnerable string = "INVULNERABLE"
	Active       string = "Active"
	Broken       string = "BROKEN"
	Disabled     string = "Disabled"
	Unavailable  string = "UNAVAILABLE"
)

//  Planet Types
const (
	Primitive string = "PRIMITIVE"
	Resource  string = "RESOURCE"
	Bot       string = "BOT"
	User      string = "USER"
	Base      string = "BASE"
)

// Resources
const (
	Water    string = "WATER"
	Graphene string = "GRAPHENE"
	Shelio   string = "SHELIO"
)

// Experience Constant Types
const (
	UserExperiences string = "USER"
	ClanExperiences string = "CLAN"
)

// Paltan Roles
const (
	PaltanLeader    string = "LEADER"
	PaltanSubLeader string = "SUB_LEADER"
	PaltanMember    string = "MEMBER"
)

//  Building States
const (
	WorkingState   string = "WORKING"
	UpgradingState string = "UPGRADING"
)

//  Building Types
const (
	WaterMiningPlant        string = "WATER_MINING_PLANT"
	GrapheneMiningPlant     string = "GRAPHENE_MINING_PLANT"
	Shield                  string = "SHIELD"
	Vikram                  string = "VIKRAM"
	PopulationControlCenter string = "POPULATION_CONTROL_CENTER"
	AttackProductionCenter  string = "ATTACK_PRODUCTION_CENTER"
	DefenceProductionCenter string = "DEFENCE_PRODUCTION_CENTER"
	DiamondStorage          string = "DIAMOND_STORAGE"
	WaterPressureTank       string = "WATER_PRESSURE_TANK"
	WarRoom                 string = "WAR_ROOM"
	CommunicationCenter     string = "COMMUNICATION_CENTER"
	ResearchLab             string = "RESEARCH_LAB"
	Library                 string = "LIBRARY"
)

//  Ship Types
const (
	Scout           string = "SCOUT"
	ResourceCarrier string = "RESOURCE CARRIER"
	AircraftCarrier string = "AIRCRAFT CARRIER"
	Attacker        string = "ATTACKER"
)

//  Ship Types
const (
	SystemDistanceBlocks float64 = 150
	SectorDistanceBlocks float64 = 20
)

//  Mission Constants
const (
	AttackMission   string = "ATTACK"
	SpyMission      string = "SPY"
	DepartureState  string = "DEPARTURE"
	InProgressState string = "IN_PROGRESS"
	ReturningState  string = "RETURNING"
	CompletedState  string = "COMPLETED"
)

type ExperienceConstants struct {
	MaxLevel            int                                `json:"max_level"`
	ExperiencesRequired map[string]ExperienceLevelConstant `json:"experiences_required"`
}

type ExperienceLevelConstant struct {
	ExperienceRequired int `json:"experience_required"`
}

type MiningConstants struct {
	MaxLevel int                            `json:"max_level"`
	Levels   map[string]MiningLevelConstant `json:"levels"`
}

type MiningLevelConstant struct {
	MiningRatePerWorker int `json:"mining_rate_per_worker"`
	WorkersMaxLimit     int `json:"workers_max_limit"`
}

type DefenceConstants struct {
	MaxLevel int                             `json:"max_level"`
	Levels   map[string]DefenceLevelConstant `json:"levels"`
}

type DefenceLevelConstant struct {
	RequiredSoldiers int `json:"required_soldiers"`
	HitPoints        int `json:"hit_points"`
	Armor            int `json:"armor"`
	MinAttack        int `json:"min_attack"`
	MaxAttack        int `json:"max_attack"`
	Range            int `json:"range"`
	SingleHitTargets int `json:"single_hit_targets"`
}

type ShipConstants struct {
	MaxLevel int                          `json:"max_level"`
	Type     string                       `json:"type"`
	Levels   map[string]ShipLevelConstant `json:"levels"`
}

type ShipLevelConstant struct {
	RequiredSoldiers int `json:"required_soldiers"`
	HitPoints        int `json:"hit_points"`
	Armor            int `json:"armor"`
	ResourceCapacity int `json:"resource_capacity"`
	WorkerCapacity   int `json:"worker_capacity"`
	MinAttack        int `json:"min_attack"`
	MaxAttack        int `json:"max_attack"`
	Range            int `json:"range"`
	Speed            int `json:"speed"`
}

type UpgradeConstants struct {
	MaxLevel int                             `json:"max_level"`
	Levels   map[string]UpgradeLevelConstant `json:"levels"`
}

type UpgradeLevelConstant struct {
	WaterRequired    int `json:"water_required"`
	GrapheneRequired int `json:"graphene_required"`
	ShelioRequired   int `json:"shelio_required"`
	MinutesRequired  int `json:"minutes_required"`
}

type BuildingConstants struct {
	Levels map[string]map[string]string `json:"levels"`
}

func GetBuildingType(buildingId string) (string, error) {
	if strings.HasPrefix(buildingId, "WMP") {
		return WaterMiningPlant, nil
	}
	if strings.HasPrefix(buildingId, "GMP") {
		return GrapheneMiningPlant, nil
	}
	if strings.HasPrefix(buildingId, "SHLD") {
		return Shield, nil
	}
	if strings.HasPrefix(buildingId, "DSC") {
		return Vikram, nil
	}
	return "", errors.New("error. invalid building id" + buildingId)
}

func GetShieldIds() []string {
	return []string{"SHLD01", "SHLD02", "SHLD03"}
}

func IsShieldId(id string) bool {
	return id == "SHLD01" || id == "SHLD02" || id == "SHLD03"
}

func GetAttackPointIds() []string {
	return []string{"point1", "point2", "point3", "point4", "point5", "point6"}
}

func GetAttackLineIds() []string {
	return []string{"line1", "line2", "line3", "line4"}
}

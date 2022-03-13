package models

import (
	"github.com/themane/MMOServer/constants"
	"math"
	"strconv"
)

type Requirements struct {
	SoldiersRequired float64 `json:"soldiers_required" example:"40"`
	WorkersRequired  float64 `json:"workers_required" example:"20"`
	GrapheneRequired float64 `json:"graphene_required" example:"101"`
	WaterRequired    float64 `json:"water_required" example:"5"`
	ShelioRequired   float64 `json:"shelio_required" example:"0"`
	MinutesRequired  float64 `json:"minutes_required" example:"10"`
}

func (r *Requirements) Init(unitLevelConstants map[string]interface{}) {
	r.SoldiersRequired = unitLevelConstants["soldiers_required"].(float64)
	r.WorkersRequired = unitLevelConstants["workers_required"].(float64)
	r.GrapheneRequired = unitLevelConstants["graphene_required"].(float64)
	r.WaterRequired = unitLevelConstants["water_required"].(float64)
	r.ShelioRequired = unitLevelConstants["shelio_required"].(float64)
	r.MinutesRequired = unitLevelConstants["minutes_required"].(float64)
}

func (r *Requirements) InitNextLevelRequirements(currentLevel int, militaryConstants constants.MilitaryConstants) {
	if currentLevel < militaryConstants.MaxLevel {
		currentLevelString := strconv.Itoa(currentLevel)
		nextLevelString := strconv.Itoa(currentLevel + 1)
		r.SoldiersRequired = militaryConstants.Levels[nextLevelString]["soldiers_required"].(float64) - militaryConstants.Levels[currentLevelString]["soldiers_required"].(float64)
		r.WorkersRequired = militaryConstants.Levels[nextLevelString]["workers_required"].(float64) - militaryConstants.Levels[currentLevelString]["workers_required"].(float64)
		r.GrapheneRequired = militaryConstants.Levels[nextLevelString]["graphene_required"].(float64)
		r.WaterRequired = militaryConstants.Levels[nextLevelString]["water_required"].(float64)
		r.ShelioRequired = militaryConstants.Levels[nextLevelString]["shelio_required"].(float64)
		r.MinutesRequired = militaryConstants.Levels[nextLevelString]["minutes_required"].(float64)
	}
}

type Returns struct {
	SoldiersReturned float64 `json:"soldiers_returned" example:"40"`
	WorkersReturned  float64 `json:"workers_returned" example:"20"`
	GrapheneReturned float64 `json:"graphene_returned" example:"101"`
	WaterReturned    float64 `json:"water_returned" example:"5"`
	ShelioReturned   float64 `json:"shelio_returned" example:"0"`
}

func (r *Returns) InitDestructionReturns(unitLevelConstants map[string]interface{}) {
	r.SoldiersReturned = math.Floor(unitLevelConstants["soldiers_required"].(float64))
	r.WorkersReturned = math.Floor(unitLevelConstants["workers_required"].(float64))
	r.GrapheneReturned = math.Floor(unitLevelConstants["graphene_required"].(float64) / 2)
	r.WaterReturned = math.Floor(unitLevelConstants["water_required"].(float64) / 2)
	r.ShelioReturned = math.Floor(unitLevelConstants["shelio_required"].(float64) / 2)
}

func (r *Returns) InitCancelReturns(unitLevelConstants map[string]interface{}, quantity float64) {
	r.SoldiersReturned = math.Floor(unitLevelConstants["soldiers_required"].(float64)) * quantity
	r.WorkersReturned = math.Floor(unitLevelConstants["workers_required"].(float64)) * quantity
	r.GrapheneReturned = math.Floor(unitLevelConstants["graphene_required"].(float64)) * quantity
	r.WaterReturned = math.Floor(unitLevelConstants["water_required"].(float64)) * quantity
	r.ShelioReturned = math.Floor(unitLevelConstants["shelio_required"].(float64)) * quantity
}

type ShipAttributes struct {
	HitPoints        float64 `json:"hit_points" example:"40"`
	Armor            float64 `json:"armor" example:"2"`
	ResourceCapacity float64 `json:"resource_capacity" example:"40"`
	WorkerCapacity   float64 `json:"worker_capacity" example:"20"`
	MinAttack        float64 `json:"min_attack" example:"5"`
	MaxAttack        float64 `json:"max_attack" example:"7"`
	Range            float64 `json:"range" example:"2"`
	Speed            float64 `json:"speed" example:"600"`
}

func (a *ShipAttributes) Init(shipLevelConstants map[string]interface{}) {
	a.HitPoints = shipLevelConstants["hit_points"].(float64)
	a.Armor = shipLevelConstants["armor"].(float64)
	a.ResourceCapacity = shipLevelConstants["resource_capacity"].(float64)
	a.WorkerCapacity = shipLevelConstants["worker_capacity"].(float64)
	a.MinAttack = shipLevelConstants["min_attack"].(float64)
	a.MaxAttack = shipLevelConstants["max_attack"].(float64)
	a.Range = shipLevelConstants["range"].(float64)
	a.Speed = shipLevelConstants["speed"].(float64)
}

type DefenceAttributes struct {
	HitPoints        float64 `json:"hit_points" example:"40"`
	Armor            float64 `json:"armor" example:"2"`
	MinAttack        float64 `json:"min_attack" example:"5"`
	MaxAttack        float64 `json:"max_attack" example:"7"`
	Range            float64 `json:"range" example:"2"`
	SingleHitTargets float64 `json:"single_hit_targets" example:"1"`
}

func (a *DefenceAttributes) Init(defenceLevelConstant map[string]interface{}) {
	a.HitPoints = defenceLevelConstant["hit_points"].(float64)
	a.Armor = defenceLevelConstant["armor"].(float64)
	a.MinAttack = defenceLevelConstant["min_attack"].(float64)
	a.MaxAttack = defenceLevelConstant["max_attack"].(float64)
	a.Range = defenceLevelConstant["range"].(float64)
}

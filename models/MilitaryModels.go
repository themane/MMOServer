package models

import (
	"github.com/themane/MMOServer/constants"
	"math"
	"strconv"
)

type Requirements struct {
	Population      Population `json:"population_required"`
	Resources       Resources  `json:"resources_required"`
	MinutesRequired float64    `json:"minutes_required" example:"10"`
}

type Population struct {
	Soldiers float64 `json:"soldiers" example:"40"`
	Workers  float64 `json:"workers" example:"20"`
}

type Resources struct {
	Graphene float64 `json:"graphene" example:"101"`
	Water    float64 `json:"water" example:"5"`
	Shelio   float64 `json:"shelio" example:"0"`
}

func (r *Requirements) Init(unitLevelConstants map[string]interface{}) {
	r.Population.Soldiers = unitLevelConstants["soldiers_required"].(float64)
	r.Population.Workers = unitLevelConstants["workers_required"].(float64)
	r.Resources.Graphene = unitLevelConstants["graphene_required"].(float64)
	r.Resources.Water = unitLevelConstants["water_required"].(float64)
	r.Resources.Shelio = unitLevelConstants["shelio_required"].(float64)
	r.MinutesRequired = unitLevelConstants["minutes_required"].(float64)
}

func (r *Requirements) InitNextLevelRequirements(currentLevel int, militaryConstants constants.MilitaryConstants) {
	if currentLevel < militaryConstants.MaxLevel {
		currentLevelString := strconv.Itoa(currentLevel)
		nextLevelString := strconv.Itoa(currentLevel + 1)
		r.Population.Soldiers = militaryConstants.Levels[nextLevelString]["soldiers_required"].(float64) - militaryConstants.Levels[currentLevelString]["soldiers_required"].(float64)
		r.Population.Workers = militaryConstants.Levels[nextLevelString]["workers_required"].(float64) - militaryConstants.Levels[currentLevelString]["workers_required"].(float64)
		r.Resources.Graphene = militaryConstants.Levels[nextLevelString]["graphene_required"].(float64)
		r.Resources.Water = militaryConstants.Levels[nextLevelString]["water_required"].(float64)
		r.Resources.Shelio = militaryConstants.Levels[nextLevelString]["shelio_required"].(float64)
		r.MinutesRequired = militaryConstants.Levels[nextLevelString]["minutes_required"].(float64)
	}
}

type Returns struct {
	Population Population `json:"population_returned"`
	Resources  Resources  `json:"resources_returned"`
}

func (r *Returns) InitDestructionReturns(unitLevelConstants map[string]interface{}) {
	r.Population.Soldiers = math.Floor(unitLevelConstants["soldiers_required"].(float64))
	r.Population.Workers = math.Floor(unitLevelConstants["workers_required"].(float64))
	r.Resources.Graphene = math.Floor(unitLevelConstants["graphene_required"].(float64) / 2)
	r.Resources.Water = math.Floor(unitLevelConstants["water_required"].(float64) / 2)
	r.Resources.Shelio = math.Floor(unitLevelConstants["shelio_required"].(float64) / 2)
}

func (r *Returns) InitCancelReturns(unitLevelConstants map[string]interface{}, quantity float64) {
	r.Population.Soldiers = math.Floor(unitLevelConstants["soldiers_required"].(float64)) * quantity
	r.Population.Workers = math.Floor(unitLevelConstants["workers_required"].(float64)) * quantity
	r.Resources.Graphene = math.Floor(unitLevelConstants["graphene_required"].(float64)) * quantity
	r.Resources.Water = math.Floor(unitLevelConstants["water_required"].(float64)) * quantity
	r.Resources.Shelio = math.Floor(unitLevelConstants["shelio_required"].(float64)) * quantity
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

package models

import (
	"math"
	"strconv"
)

type ResourceType string

const (
	WATER ResourceType = "WATER"
	GRAPHENE
	SHELIO
)

type ResourceConstants struct {
	MaxLevel int                      `json:"max_level"`
	Levels   map[string]LevelConstant `json:"levels"`
}

type LevelConstant struct {
	WaterRequired       int `json:"water_required"`
	GrapheneRequired    int `json:"graphene_required"`
	ShelioRequired      int `json:"shelio_required"`
	MiningRatePerWorker int `json:"mining_rate_per_worker"`
	WorkersMaxLimit     int `json:"workers_max_limit"`
}

type Resources struct {
	Water    Resource `json:"water"`
	Graphene Resource `json:"graphene"`
	Shelio   int      `json:"shelio" example:"23"`
}

type Resource struct {
	MaxLimit float64 `json:"max_limit" example:"100"`
	Amount   int     `json:"amount" example:"23"`
	Reserved int     `json:"reserved" example:"14"`
}

func (r *Resources) Init(planetUser PlanetUser) {
	limit := getMaxLimit(planetUser.Water.Amount, planetUser.Graphene.Amount)
	r.Water = Resource{Amount: planetUser.Water.Amount, Reserved: planetUser.Water.Reserved, MaxLimit: limit}
	r.Graphene = Resource{Amount: planetUser.Graphene.Amount, Reserved: planetUser.Graphene.Reserved, MaxLimit: limit}
	r.Shelio = planetUser.Shelio
}

func getMaxLimit(water int, graphene int) float64 {
	var biggerAmount int
	if water > graphene {
		biggerAmount = water
	} else {
		biggerAmount = graphene
	}
	nDigits := len(strconv.Itoa(biggerAmount))
	return math.Pow10(nDigits)
}

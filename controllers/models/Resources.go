package models

import (
	"github.com/themane/MMOServer/mongoRepository/models"
	"math"
	"strconv"
)

type Resources struct {
	Water    Resource `json:"water"`
	Graphene Resource `json:"graphene"`
	Shelio   int      `json:"shelio" example:"23"`
}

type Resource struct {
	MaxLimit  float64 `json:"max_limit" example:"100"`
	Amount    int     `json:"amount" example:"23"`
	Reserved  int     `json:"reserved" example:"14"`
	Reserving int     `json:"reserving" example:"10"`
	Total     int     `json:"total" example:"47"`
}

func InitResources(planetUser models.PlanetUser) *Resources {
	r := new(Resources)
	limit := getMaxLimit(planetUser.Water.Amount, planetUser.Graphene.Amount)
	r.Water = Resource{
		Amount:    planetUser.Water.Amount,
		Reserved:  planetUser.Water.Reserved,
		Reserving: planetUser.Water.Reserving,
		Total:     planetUser.Water.Amount + planetUser.Water.Reserved + planetUser.Water.Reserving,
		MaxLimit:  limit,
	}
	r.Graphene = Resource{
		Amount:    planetUser.Graphene.Amount,
		Reserved:  planetUser.Graphene.Reserved,
		Reserving: planetUser.Graphene.Reserving,
		Total:     planetUser.Graphene.Amount + planetUser.Graphene.Reserved + planetUser.Graphene.Reserving,
		MaxLimit:  limit,
	}
	r.Shelio = planetUser.Shelio
	return r
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

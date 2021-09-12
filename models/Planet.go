package models

import "fmt"

type UnoccupiedPlanet struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
}

type OccupiedPlanet struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
	Resources    Resources      `json:"resources"`
	Population   Population     `json:"population"`
	Mines        []Mine         `json:"mines"`
	Home         bool           `json:"home" example:"true"`
}

type PlanetPosition struct {
	System int    `json:"system" example:"23"`
	Sector int    `json:"sector" example:"49"`
	Planet int    `json:"planet" example:"7"`
	Id     string `json:"_id" example:"023:049:07"`
}

func (position PlanetPosition) SystemId() string {
	return fmt.Sprintf("%03d", position.System)
}

func (position PlanetPosition) SectorId() string {
	return fmt.Sprintf("%03d:%03d", position.System, position.Sector)
}

func (position PlanetPosition) PlanetId() string {
	return fmt.Sprintf("%03d:%03d:%02d", position.System, position.Sector, position.Planet)
}

func (position PlanetPosition) Clone() PlanetPosition {
	return PlanetPosition{Id: position.PlanetId(), System: position.System, Sector: position.Sector, Planet: position.Planet}
}

type Resources struct {
	Water    Resource `json:"water"`
	Graphene Resource `json:"graphene"`
	Shelio   int      `json:"shelio" example:"23"`
}

type Resource struct {
	GenerationRate int `json:"generation_rate" example:"3"`
	MaxLimit       int `json:"max_limit" example:"100"`
	Amount         int `json:"amount" example:"23"`
	Reserved       int `json:"reserved" example:"14"`
}

type Population struct {
	Total          int                `json:"total" example:"45"`
	GenerationRate int                `json:"generation_rate" example:"3"`
	Unemployed     int                `json:"unemployed" example:"3"`
	Workers        EmployedPopulation `json:"workers"`
	Soldiers       EmployedPopulation `json:"soldiers"`
}

type EmployedPopulation struct {
	Total int `json:"total" example:"21"`
	Idle  int `json:"idle" example:"4"`
}

type Mine struct {
	Id          string      `json:"_id" example:"W101"`
	Type        string      `json:"type" example:"WATER"`
	MaxLimit    int         `json:"max_limit" example:"550"`
	Mined       int         `json:"mined" example:"125"`
	MiningPlant MiningPlant `json:"mining_plant"`
}

type MiningPlant struct {
	BuildingId          string              `json:"building_id" example:"WMP101"`
	BuildingLevel       int                 `json:"building_level" example:"3"`
	Workers             int                 `json:"workers" example:"12"`
	NextLevelAttributes NextLevelAttributes `json:"next_level"`
}

type NextLevelAttributes struct {
	GrapheneRequired           int `json:"graphene_required" example:"101"`
	WaterRequired              int `json:"water_required" example:"5"`
	ShelioRequired             int `json:"shelio_required" example:"0"`
	CurrentMiningRatePerWorker int `json:"current_mining_rate_per_worker" example:"1"`
	NextMiningRatePerWorker    int `json:"next_mining_rate_per_worker" example:"1"`
	CurrentWorkersMaxLimit     int `json:"current_workers_max_limit" example:"40"`
	NextWorkersMaxLimit        int `json:"next_workers_max_limit" example:"65"`
}

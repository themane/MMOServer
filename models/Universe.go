package models

import "fmt"

type Universe struct {
	NumSystems int       `json:"num_systems"`
	Systems    []SystemU `json:"systems"`
}

type SystemU struct {
	System     int       `json:"system"`
	NumSystems int       `json:"num_systems"`
	Sectors    []SectorU `json:"sectors"`
}

type SectorU struct {
	Sector     int       `json:"sector"`
	NumPlanets int       `json:"num_planets"`
	Planets    []PlanetU `json:"planets"`
}

type PlanetU struct {
	Planet       int     `json:"planet"`
	Mines        []MineU `json:"mines"`
	PlanetConfig string  `json:"planet_config"`
}

type MineU struct {
	Id           string `json:"_id"`
	Type         string `json:"type"`
	MaxLimit     int    `json:"max_limit"`
	IncreaseRate int    `json:"increase_rate"`
}

func SystemId(System int) string {
	return fmt.Sprintf("%03d", System)
}

func SectorId(Sector int) string {
	return fmt.Sprintf("%03d", Sector)
}

func PlanetId(System int, Sector int, Planet int) string {
	return fmt.Sprintf("%03d:%03d:%02d", System, Sector, Planet)
}

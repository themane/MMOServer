package models

import "fmt"

type Universe struct {
	NumSystems int                `json:"num_systems"`
	Systems    map[string]SystemU `json:"systems"`
}

type SystemU struct {
	System     int                `json:"system"`
	NumSystems int                `json:"num_systems"`
	Sectors    map[string]SectorU `json:"sectors"`
}

type SectorU struct {
	Sector     int                `json:"sector"`
	NumPlanets int                `json:"num_planets"`
	Planets    map[string]PlanetU `json:"planets"`
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

func SystemIdString(System int) string {
	return fmt.Sprintf("%03d", System)
}

func SystemId(Position PlanetPosition) string {
	return SystemIdString(Position.System)
}

func SectorIdString(Sector int) string {
	return fmt.Sprintf("%03d", Sector)
}
func SectorId(Position PlanetPosition) string {
	return SectorIdString(Position.Sector)
}
func PlanetIdString(System int, Sector int, Planet int) string {
	return fmt.Sprintf("%03d:%03d:%02d", System, Sector, Planet)
}

func PlanetId(Position PlanetPosition) string {
	return PlanetIdString(Position.System, Position.Sector, Position.Planet)
}

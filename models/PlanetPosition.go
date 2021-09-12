package models

import "fmt"

type PlanetPosition struct {
	Id     string `json:"_id" example:"023:049:07"`
	System int    `json:"system" example:"23"`
	Sector int    `json:"sector" example:"49"`
	Planet int    `json:"planet" example:"7"`
}

func (p *PlanetPosition) Init(system int, sector int, planet int) {
	p.System = system
	p.Sector = sector
	p.Planet = planet
	p.Id = PlanetId(system, sector, planet)
}

func (p PlanetPosition) SystemId() string {
	return fmt.Sprintf("%03d", p.System)
}

func (p PlanetPosition) SectorId() string {
	return fmt.Sprintf("%03d:%03d", p.System, p.Sector)
}

func (p PlanetPosition) PlanetId() string {
	return fmt.Sprintf("%03d:%03d:%02d", p.System, p.Sector, p.Planet)
}

func PlanetId(system int, sector int, planet int) string {
	return fmt.Sprintf("%03d:%03d:%02d", system, sector, planet)
}

func (p PlanetPosition) Clone() PlanetPosition {
	return PlanetPosition{Id: p.PlanetId(), System: p.System, Sector: p.Sector, Planet: p.Planet}
}

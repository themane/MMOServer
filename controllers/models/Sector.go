package models

import "github.com/themane/MMOServer/models"

type Sector struct {
	OccupiedPlanets   []OccupiedPlanet   `json:"occupied_planets"`
	UnoccupiedPlanets []UnoccupiedPlanet `json:"unoccupied_planets"`
	Position          SectorPosition     `json:"position"`
}

type SectorPosition struct {
	Id     string `json:"_id" example:"023:049"`
	System int    `json:"system" example:"23"`
	Sector int    `json:"sector" example:"49"`
}

func (sp *SectorPosition) Init(planetPosition models.PlanetPosition) {
	sp.Id = planetPosition.SectorId()
	sp.System = planetPosition.System
	sp.Sector = planetPosition.Sector
}

type VisitSectorRequest struct {
	Username string `json:"username" example:"devashish"`
	Sector   string `json:"sector" example:"005:001"`
}

type TeleportRequest struct {
	Username string `json:"username" example:"devashish"`
	Planet   string `json:"planet" example:"005:001:03"`
}

package models

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

func (sp *SectorPosition) Init(planetPosition PlanetPosition) {
	sp.Id = planetPosition.SectorId()
	sp.System = planetPosition.System
	sp.Sector = planetPosition.Sector
}

package models

type UnoccupiedPlanet struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
}

func (u *UnoccupiedPlanet) Init(planet PlanetUni, sectorPosition SectorPosition) {
	planetId := PlanetId(sectorPosition.System, sectorPosition.Sector, planet.Planet)
	u.Position = PlanetPosition{System: sectorPosition.System, Sector: sectorPosition.Sector, Planet: planet.Planet, Id: planetId}
	u.PlanetConfig = planet.PlanetConfig
}

type OccupiedPlanet struct {
	PlanetConfig string         `json:"planet_config" example:"Planet2.json"`
	Position     PlanetPosition `json:"position"`
	Resources    Resources      `json:"resources"`
	Population   Population     `json:"population"`
	Mines        []Mine         `json:"mines"`
	Home         bool           `json:"home" example:"true"`
}

func (o *OccupiedPlanet) Init(planetUni PlanetUni, planetUser PlanetUser, sectorPosition SectorPosition,
	waterConstants ResourceConstants, grapheneConstants ResourceConstants) {

	planetId := PlanetId(sectorPosition.System, sectorPosition.Sector, planetUni.Planet)
	o.Position = PlanetPosition{System: sectorPosition.System, Sector: sectorPosition.Sector, Planet: planetUni.Planet, Id: planetId}
	o.PlanetConfig = planetUni.PlanetConfig
	o.Resources.Init(planetUser)
	o.Population.Init(planetUser)
	for mineId, mineUser := range planetUser.Mines {
		mine := Mine{}
		mine.Init(planetUni.Mines[mineId], mineUser, waterConstants, grapheneConstants)
		o.Mines = append(o.Mines, mine)
	}
	o.Home = planetUser.Home
}

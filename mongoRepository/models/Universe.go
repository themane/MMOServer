package models

import (
	"github.com/themane/MMOServer/models"
)

type PlanetUni struct {
	Id           string                `json:"_id"`
	Position     models.PlanetPosition `json:"position"`
	Mines        map[string]MineUni    `json:"mines"`
	PlanetConfig string                `json:"planet_config"`
	Occupied     string                `json:"occupied"`
	Distance     int                   `json:"distance"`
}

type MineUni struct {
	Id           string `json:"_id"`
	Type         string `json:"type"`
	MaxLimit     int    `json:"max_limit"`
	IncreaseRate int    `json:"increase_rate"`
}

type UniverseRepository interface {
	GetSector(system int, sector int) (map[string]PlanetUni, error)
	GetPlanet(system int, sector int, planet int) (*PlanetUni, error)
	GetAllOccupiedPlanets(system int) (map[string]PlanetUni, error)
	GetRandomUnoccupiedPlanet(system int) (*PlanetUni, error)
	MarkOccupied(system int, sector int, planet int, userId string) error
}

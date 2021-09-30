package models

import (
	"github.com/themane/MMOServer/models"
)

type PlanetUni struct {
	Id           string                `json:"_id" bson:"_id"`
	Position     models.PlanetPosition `json:"position" bson:"position"`
	Distance     int                   `json:"distance" bson:"distance"`
	Mines        map[string]MineUni    `json:"mines" bson:"mines"`
	PlanetConfig string                `json:"planet_config" bson:"planet_config"`
	Occupied     string                `json:"occupied" bson:"occupied"`
	Workers      int                   `json:"workers" bson:"workers"`
	BasePlanet   bool                  `json:"base_planet" bson:"base_planet"`
}

type MineUni struct {
	Id           string `json:"_id" bson:"_id"`
	Type         string `json:"type" bson:"type"`
	MaxLimit     int    `json:"max_limit" bson:"max_limit"`
	IncreaseRate int    `json:"increase_rate" bson:"increase_rate"`
}

type UniverseRepository interface {
	FindById(id string) (*PlanetUni, error)
	FindByPosition(system int, sector int, planet int) (*PlanetUni, error)

	GetSector(system int, sector int) (map[string]PlanetUni, error)
	GetAllOccupiedPlanets(system int) (map[string]PlanetUni, error)
	GetRandomUnoccupiedPlanet(system int) (*PlanetUni, error)
	MarkOccupied(system int, sector int, planet int, userId string) error
}

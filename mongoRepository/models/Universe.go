package models

import (
	"github.com/themane/MMOServer/models"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type PlanetUni struct {
	Id           string                `json:"_id"`
	Position     models.PlanetPosition `json:"position"`
	Mines        map[string]MineUni    `json:"mines"`
	PlanetConfig string                `json:"planet_config"`
	Occupied     uuid.UUID             `json:"occupied"`
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
	MarkOccupied(system int, sector int, planet int, userId uuid.UUID) error
}

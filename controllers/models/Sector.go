package models

import (
	"github.com/themane/MMOServer/models"
)

type Sector struct {
	OccupiedPlanets   []OccupiedPlanet      `json:"occupied_planets"`
	UnoccupiedPlanets []UnoccupiedPlanet    `json:"unoccupied_planets"`
	Position          models.SectorPosition `json:"position"`
}

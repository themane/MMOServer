package models

import (
	"github.com/themane/MMOServer/models"
	"log"
	"strconv"
	"strings"
)

type Sector struct {
	OccupiedPlanets   []OccupiedPlanet      `json:"occupied_planets"`
	UnoccupiedPlanets []UnoccupiedPlanet    `json:"unoccupied_planets"`
	Position          models.SectorPosition `json:"position"`
}

type VisitSectorRequest struct {
	Username string `json:"username" example:"devashish"`
	Sector   string `json:"sector" example:"005:001"`
}

type TeleportRequest struct {
	Username string `json:"username" example:"devashish"`
	Planet   string `json:"planet" example:"005:001:03"`
}

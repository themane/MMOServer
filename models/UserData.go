package models

type UserData struct {
	Profile         Profile          `json:"profile"`
	OccupiedPlanets []OccupiedPlanet `json:"occupied_planets"`
}

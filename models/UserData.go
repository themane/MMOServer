package models

type UserData struct {
	Profile Profile          `json:"profile"`
	Planets []OccupiedPlanet `json:"occupied_planets"`
}

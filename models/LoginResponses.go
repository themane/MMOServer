package models

type LoginResponse struct {
	Profile Profile  `json:"profile"`
	Sectors []Sector `json:"occupied_sectors"`
}

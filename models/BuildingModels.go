package models

const (
	PRODUCING_STATE = "PRODUCING"
	UPGRADING_STATE = "UPGRADING"
)

type BuildingState struct {
	State                    string `json:"state" example:"PRODUCING"`
	BuildingMinutesPerWorker int    `json:"building_minutes_per_worker" example:"1440"`
	Workers                  int    `json:"workers" example:"6"`
}

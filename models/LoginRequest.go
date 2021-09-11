package models

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	PlanetConfig string     `json:"planet_config"`
	Position     Position   `json:"position"`
	Water        Resource   `json:"water"`
	Graphene     Resource   `json:"graphene"`
	Shelio       Resource   `json:"shelio"`
	Population   Population `json:"population"`
	Mines        []Mine     `json:"mines"`
}

type Position struct {
	System  int    `json:"system"`
	Sector  int    `json:"sector"`
	Planet  int    `json:"planet"`
	Display string `json:"display"`
}

type Resource struct {
	GenerationRate int `json:"generation_rate"`
	Amount         int `json:"amount"`
	Reserved       int `json:"reserved"`
}

type Population struct {
	Total          int                `json:"total"`
	GenerationRate int                `json:"generation_rate"`
	Unemployed     int                `json:"unemployed"`
	Workers        EmployedPopulation `json:"workers"`
	Soldiers       EmployedPopulation `json:"soldiers"`
}

type EmployedPopulation struct {
	Total int `json:"total"`
	Idle  int `json:"idle"`
}

type Mine struct {
	Id          string      `json:"_id"`
	Type        string      `json:"type"`
	MaxLimit    int         `json:"max_limit"`
	Mined       int         `json:"mined"`
	MiningPlant MiningPlant `json:"mining_plant"`
}

type MiningPlant struct {
	BuildingId          string              `json:"building_id"`
	BuildingLevel       int                 `json:"building_level"`
	MiningRate          int                 `json:"mining_rate"`
	Workers             int                 `json:"workers"`
	NextLevelAttributes NextLevelAttributes `json:"next_level"`
}

type NextLevelAttributes struct {
	GrapheneRequired           int `json:"graphene_required"`
	WaterRequired              int `json:"water_required"`
	ShelioRequired             int `json:"shelio_required"`
	CurrentMiningRatePerWorker int `json:"current_mining_rate_per_worker"`
	NextMiningRatePerWorker    int `json:"next_mining_rate_per_worker"`
	CurrentWorkersMaxLimit     int `json:"current_workers_max_limit"`
	NextWorkersMaxLimit        int `json:"next_workers_max_limit"`
}

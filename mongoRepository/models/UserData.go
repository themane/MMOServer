package models

import (
	"github.com/themane/MMOServer/models"
)

type UserData struct {
	Id              string                `json:"_id" bson:"_id" `
	Profile         ProfileUser           `json:"profile" bson:"profile"`
	OccupiedPlanets map[string]PlanetUser `json:"occupied_planets" bson:"occupied_planets"`
}

type ProfileUser struct {
	Username   string `json:"username" bson:"username"`
	Experience int    `json:"experience" bson:"experience"`
	ClanId     string `json:"clan_id" bson:"clan_id"`
}

type PlanetUser struct {
	Water               Resource                      `json:"water" bson:"water"`
	Graphene            Resource                      `json:"graphene" bson:"graphene"`
	Shelio              int                           `json:"shelio" bson:"shelio"`
	Population          Population                    `json:"population" bson:"population"`
	Mines               map[string]MineUser           `json:"mines" bson:"mines"`
	Defences            map[string]Defence            `json:"defences" bson:"defences"`
	DefenceShipCarriers map[string]DefenceShipCarrier `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	Buildings           map[string]Building           `json:"buildings" bson:"buildings"`
	Home                bool                          `json:"home" bson:"home"`
}

type Resource struct {
	Amount   int `json:"amount" bson:"amount"`
	Reserved int `json:"reserved" bson:"reserved"`
}

type Population struct {
	GenerationRate int                       `json:"generation_rate" bson:"generation_rate"`
	Unemployed     int                       `json:"unemployed" bson:"unemployed"`
	Workers        models.EmployedPopulation `json:"workers" bson:"workers"`
	Soldiers       models.EmployedPopulation `json:"soldiers" bson:"soldiers"`
}

type MineUser struct {
	Mined         int    `json:"mined" bson:"mined"`
	MiningPlantId string `json:"mining_plant_id" bson:"mining_plant_id"`
}

type Building struct {
	BuildingLevel            int `json:"building_level" bson:"building_level"`
	Workers                  int `json:"workers" bson:"workers"`
	BuildingMinutesPerWorker int `json:"building_minutes_per_worker" bson:"building_minutes_per_worker"`
}

type Defence struct {
	HitPoints      int            `json:"hit_points" bson:"hit_points"`
	Attack         int            `json:"attack" bson:"attack"`
	Range          int            `json:"range" bson:"range"`
	Target         int            `json:"target" bson:"target"`
	GuardingShield map[string]int `json:"guarding_shield" bson:"guarding_shield"`
}

type DefenceShipCarrier struct {
	HitPoints      int            `json:"hit_points" bson:"hit_points"`
	HostingShips   map[string]int `json:"hosting_ships" bson:"hosting_ships"`
	GuardingShield string         `json:"guarding_shield" bson:"guarding_shield"`
}

type UserRepository interface {
	FindById(id string) (*UserData, error)
	FindByUsername(username string) (*UserData, error)

	AddExperience(id string, experience int) error
	UpdateClanId(id string, clanId string) error

	UpgradeBuildingLevel(id string, planetId string, buildingId string, waterRequired int, grapheneRequired int, shelioRequired int, minutesRequired int) error
	AddResources(id string, planetId string, water int, graphene int, shelio int) error
	UpdateMineResources(id string, planetId string, mineId string, water int, graphene int) error
	UpdateWorkers(id string, planetId string, buildingId string, workers int) error
	AddPopulation(id string, planetId string, population int) error
	RecruitWorkers(id string, planetId string, worker int) error
	RecruitSoldiers(id string, planetId string, soldiers int) error

	ScheduledPopulationIncrease(id string, planetIdGenerationRateMap map[string]int) error
	ScheduledWaterIncrease(id string, planetIdGenerationRateMap map[string]map[string]int) error
	ScheduledGrapheneIncrease(id string, planetIdGenerationRateMap map[string]map[string]int) error
	ScheduledPopulationConsumption(id string, planetIdGenerationRateMap map[string]int) error
}

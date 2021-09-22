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
	Water      ResourceUser            `json:"water" bson:"water"`
	Graphene   ResourceUser            `json:"graphene" bson:"graphene"`
	Shelio     int                     `json:"shelio" bson:"shelio"`
	Population PopulationUser          `json:"population" bson:"population"`
	Mines      map[string]MineUser     `json:"mines" bson:"mines"`
	Buildings  map[string]BuildingUser `json:"buildings" bson:"buildings"`
	Home       bool                    `json:"home" bson:"home"`
}

type ResourceUser struct {
	Amount   int `json:"amount" bson:"amount"`
	Reserved int `json:"reserved" bson:"reserved"`
}

type PopulationUser struct {
	GenerationRate int                       `json:"generation_rate" bson:"generation_rate"`
	Unemployed     int                       `json:"unemployed" bson:"unemployed"`
	Workers        models.EmployedPopulation `json:"workers" bson:"workers"`
	Soldiers       models.EmployedPopulation `json:"soldiers" bson:"soldiers"`
}

type MineUser struct {
	Mined         int    `json:"mined" bson:"mined"`
	MiningPlantId string `json:"mining_plant_id" bson:"mining_plant_id"`
}

type BuildingUser struct {
	BuildingLevel            int `json:"building_level" bson:"building_level"`
	Workers                  int `json:"workers" bson:"workers"`
	BuildingMinutesPerWorker int `json:"building_minutes_per_worker" example:"1440"`
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

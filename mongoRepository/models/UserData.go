package models

import "github.com/themane/MMOServer/models"

type UserData struct {
	Id              string                `json:"_id"`
	Profile         ProfileUser           `json:"profile"`
	OccupiedPlanets map[string]PlanetUser `json:"occupied_planets"`
}

type ProfileUser struct {
	Username   string `json:"username" example:"nehal"`
	Experience int    `json:"experience" example:"153"`
	ClanId     string `json:"clan_id" example:"MindKrackers"`
}

type PlanetUser struct {
	Water      ResourceUser            `json:"water"`
	Graphene   ResourceUser            `json:"graphene"`
	Shelio     int                     `json:"shelio" example:"23"`
	Population PopulationUser          `json:"population"`
	Mines      map[string]MineUser     `json:"mines"`
	Buildings  map[string]BuildingUser `json:"buildings"`
	Home       bool                    `json:"home" example:"true"`
}

type ResourceUser struct {
	Amount   int `json:"amount" example:"23"`
	Reserved int `json:"reserved" example:"14"`
}

type PopulationUser struct {
	GenerationRate int                       `json:"generation_rate" example:"3"`
	Unemployed     int                       `json:"unemployed" example:"3"`
	Workers        models.EmployedPopulation `json:"workers"`
	Soldiers       models.EmployedPopulation `json:"soldiers"`
}

type MineUser struct {
	Mined         int    `json:"mined" example:"125"`
	MiningPlantId string `json:"mining_plant_id"`
}

type BuildingUser struct {
	BuildingLevel int `json:"building_level" example:"3"`
	Workers       int `json:"workers" example:"12"`
}

type UserRepository interface {
	FindById(id string) (*UserData, error)
	FindByUsername(username string) (*UserData, error)

	AddExperience(id string, experience int) error
	UpdateClanId(id string, clanId string) error

	UpgradeBuildingLevel(id string, planetId string, buildingId string, waterRequired int, grapheneRequired int, shelioRequired int) error
	AddResources(id string, planetId string, water int, graphene int, shelio int) error
	UpdateMineResources(id string, planetId string, mineId string, water int, graphene int) error
	UpdateWorkers(id string, planetId string, buildingId string, workers int) error
	AddPopulation(id string, planetId string, population int) error
	RecruitWorkers(id string, planetId string, worker int) error
	RecruitSoldiers(id string, planetId string, soldiers int) error
}

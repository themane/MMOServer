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

	AddExperience(username string, experience int) error
	UpdateClanId(username string, clanId string) error

	UpgradeBuildingLevel(username string, planetId string, buildingId string, waterRequired int, grapheneRequired int, shelioRequired int) error
	AddResources(username string, planetId string, water int, graphene int, shelio int) error
	UpdateMineResources(username string, planetId string, water int, graphene int) error
	UpdateWorkers(username string, planetId string, buildingId string, workers int) error
	AddPopulation(username string, planetId string, population int) error
	RecruitWorkers(username string, planetId string, worker int) error
	RecruitSoldiers(username string, planetId string, soldiers int) error
}

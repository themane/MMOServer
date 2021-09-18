package models

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
	Water      ResourceUser        `json:"water"`
	Graphene   ResourceUser        `json:"graphene"`
	Shelio     int                 `json:"shelio" example:"23"`
	Population PopulationUser      `json:"population"`
	Mines      map[string]MineUser `json:"mines"`
	Home       bool                `json:"home" example:"true"`
}

type ResourceUser struct {
	Amount   int `json:"amount" example:"23"`
	Reserved int `json:"reserved" example:"14"`
}

type PopulationUser struct {
	GenerationRate int                `json:"generation_rate" example:"3"`
	Unemployed     int                `json:"unemployed" example:"3"`
	Workers        EmployedPopulation `json:"workers"`
	Soldiers       EmployedPopulation `json:"soldiers"`
}

type MineUser struct {
	Mined       int `json:"mined" example:"125"`
	MiningPlant struct {
		Id            string `json:"building_id" example:"WMP101"`
		BuildingLevel int    `json:"building_level" example:"3"`
		Workers       int    `json:"workers" example:"12"`
	} `json:"mining_plant"`
}

type UserRepository interface {
	FindById(id string) (*UserData, error)
	FindByUsername(username string) (*UserData, error)
}

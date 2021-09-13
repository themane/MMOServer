package models

type UserData struct {
	Profile         Profile               `json:"profile"`
	OccupiedPlanets map[string]PlanetUser `json:"occupied_planets"`
}

type PlanetUser struct {
	Position   PlanetPosition      `json:"position"`
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
	Id          string          `json:"_id" example:"W101"`
	Type        ResourceType    `json:"type" example:"WATER"`
	Mined       int             `json:"mined" example:"125"`
	MiningPlant MiningPlantUser `json:"mining_plant"`
}

type MiningPlantUser struct {
	BuildingId    string `json:"building_id" example:"WMP101"`
	BuildingLevel int    `json:"building_level" example:"3"`
	Workers       int    `json:"workers" example:"12"`
}

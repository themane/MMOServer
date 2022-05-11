package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type UserData struct {
	Id              string                `json:"_id" bson:"_id" `
	Profile         ProfileUser           `json:"profile" bson:"profile"`
	OccupiedPlanets map[string]PlanetUser `json:"occupied_planets" bson:"occupied_planets"`
}

type ProfileUser struct {
	Username   string `json:"username" bson:"username"`
	Experience int    `json:"experience" bson:"experience"`
	Species    string `json:"species" bson:"species"`
	ClanId     string `json:"clan_id" bson:"clan_id"`
}

type PlanetUser struct {
	Water               Resource                      `json:"water" bson:"water"`
	Graphene            Resource                      `json:"graphene" bson:"graphene"`
	Shelio              int                           `json:"shelio" bson:"shelio"`
	Population          Population                    `json:"population" bson:"population"`
	Mines               map[string]MineUser           `json:"mines" bson:"mines"`
	Ships               map[string]Ship               `json:"ships" bson:"ships"`
	Defences            map[string]Defence            `json:"defences" bson:"defences"`
	DefenceShipCarriers map[string]DefenceShipCarrier `json:"defence_ship_carriers" bson:"defence_ship_carriers"`
	Buildings           map[string]Building           `json:"buildings" bson:"buildings"`
	Researches          map[string]ResearchUser       `json:"researches" bson:"researches"`
	HomePlanet          bool                          `json:"home_planet" bson:"home_planet"`
	BasePlanet          bool                          `json:"base_planet" bson:"base_planet"`
}

func (p *PlanetUser) GetAvailableShip(shipName string) int {
	var defenceShipCarrierDeployed int
	for _, defenceShipCarrier := range p.DefenceShipCarriers {
		defenceShipCarrierDeployed += defenceShipCarrier.HostingShips[shipName]
	}
	return p.Ships[shipName].Quantity - defenceShipCarrierDeployed
}

func (p *PlanetUser) GetAvailableShips() map[string]int {
	response := map[string]int{}
	for shipName, ship := range p.Ships {
		var defenceShipCarrierDeployed int
		for _, defenceShipCarrier := range p.DefenceShipCarriers {
			defenceShipCarrierDeployed += defenceShipCarrier.HostingShips[shipName]
		}
		response[shipName] = ship.Quantity - defenceShipCarrierDeployed
	}
	return response
}

type Resource struct {
	Amount    int `json:"amount" bson:"amount"`
	Reserved  int `json:"reserved" bson:"reserved"`
	Reserving int `json:"reserving" bson:"reserving"`
}

type ResearchUser struct {
	Level                    int `json:"level" bson:"level"`
	ResearchMinutesPerWorker int `json:"research_minutes_per_worker" bson:"research_minutes_per_worker"`
}

type Population struct {
	GenerationRate int `json:"generation_rate" bson:"generation_rate"`
	Unemployed     int `json:"unemployed" bson:"unemployed"`
	IdleWorkers    int `json:"workers" bson:"workers"`
	IdleSoldiers   int `json:"soldiers" bson:"soldiers"`
}

func GetEmployedPopulation(planetUser PlanetUser, militaryConstants map[string]constants.MilitaryConstants) (int, int) {
	totalEmployedWorkers := 0
	totalEmployedSoldiers := 0
	for _, building := range planetUser.Buildings {
		totalEmployedWorkers += building.Workers
		totalEmployedSoldiers += building.Soldiers
	}
	for shipName, ship := range planetUser.Ships {
		if ship.Level > 0 {
			totalEmployedWorkers += ship.Quantity * int(militaryConstants[shipName].Levels[strconv.Itoa(ship.Level)]["workers_required"].(float64))
			totalEmployedSoldiers += ship.Quantity * int(militaryConstants[shipName].Levels[strconv.Itoa(ship.Level)]["soldiers_required"].(float64))
		}
	}
	for defenceName, defence := range planetUser.Defences {
		if defence.Level > 0 {
			totalEmployedWorkers += defence.Quantity * int(militaryConstants[defenceName].Levels[strconv.Itoa(defence.Level)]["workers_required"].(float64))
			totalEmployedSoldiers += defence.Quantity * int(militaryConstants[defenceName].Levels[strconv.Itoa(defence.Level)]["soldiers_required"].(float64))
		}
	}
	for _, defenceShipCarrier := range planetUser.DefenceShipCarriers {
		if defenceShipCarrier.Level > 0 {
			totalEmployedWorkers += int(militaryConstants[defenceShipCarrier.Name].Levels[strconv.Itoa(defenceShipCarrier.Level)]["workers_required"].(float64))
			totalEmployedSoldiers += int(militaryConstants[defenceShipCarrier.Name].Levels[strconv.Itoa(defenceShipCarrier.Level)]["soldiers_required"].(float64))
		}
	}
	return totalEmployedWorkers, totalEmployedSoldiers
}

type MineUser struct {
	Mined int `json:"mined" bson:"mined"`
}

type Building struct {
	BuildingLevel            int `json:"building_level" bson:"building_level"`
	Workers                  int `json:"workers" bson:"workers"`
	Soldiers                 int `json:"soldiers" bson:"soldiers"`
	BuildingMinutesPerWorker int `json:"building_minutes_per_worker" bson:"building_minutes_per_worker"`
}

type Ship struct {
	Level             int               `json:"level" bson:"level"`
	Quantity          int               `json:"quantity" bson:"quantity"`
	UnderConstruction UnderConstruction `json:"under_construction" bson:"under_construction"`
}

func GetAvailableShips(unitName string, attackMissions []AttackMission, defenceShipCarriers map[string]DefenceShipCarrier, totalUnits int) int {
	deployedUnits := 0
	for _, mission := range attackMissions {
		for _, shieldFormation := range mission.Formation {
			for _, pointFormation := range shieldFormation {
				deployedUnits += pointFormation[unitName]
			}
		}
	}
	for _, defenceShipCarrier := range defenceShipCarriers {
		deployedUnits += defenceShipCarrier.HostingShips[unitName]
	}
	return totalUnits - deployedUnits
}

func GetAvailableScouts(unitName string, spyMissions []SpyMission, totalUnits int) int {
	deployedUnits := 0
	for _, mission := range spyMissions {
		deployedUnits += mission.Scouts[unitName]
	}
	return totalUnits - deployedUnits
}

type Defence struct {
	Level             int               `json:"level" bson:"level"`
	Quantity          int               `json:"quantity" bson:"quantity"`
	GuardingShield    map[string]int    `json:"guarding_shield" bson:"guarding_shield"`
	UnderConstruction UnderConstruction `json:"under_construction" bson:"under_construction"`
}

func GetIdleDefences(guardingShield map[string]int, totalUnits int) int {
	deployedUnits := 0
	for _, quantity := range guardingShield {
		deployedUnits += quantity
	}
	return totalUnits - deployedUnits
}

type DefenceShipCarrier struct {
	Name              string            `json:"name" bson:"name"`
	Level             int               `json:"level" bson:"level"`
	HostingShips      map[string]int    `json:"hosting_ships" bson:"hosting_ships"`
	GuardingShield    string            `json:"guarding_shield" bson:"guarding_shield"`
	UnderConstruction UnderConstruction `json:"under_construction" bson:"under_construction"`
}

type UnderConstruction struct {
	StartTime primitive.DateTime `json:"start_time" bson:"start_time"`
	Quantity  int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type UserRepository interface {
	FindById(id string) (*UserData, error)
	FindByUsername(username string) (*UserData, error)

	AddExperience(id string, experience int) error
	UpdateClanId(id string, clanId string) error

	UpgradeBuildingLevel(id string, planetId string, buildingId string, waterRequired int, grapheneRequired int, shelioRequired int, minutesRequired int) error
	CancelUpgradeBuildingLevel(id string, planetId string, buildingId string, waterReturned int, grapheneReturned int, shelioReturned int) error
	UpdateWorkers(id string, planetId string, buildingId string, workers int) error
	UpdateSoldiers(id string, planetId string, buildingId string, soldiers int) error
	UpdatePopulationRate(id string, planetId string, generationRate int) error

	Recruit(id string, planetId string, worker int, soldiers int) error
	KillPopulation(id string, planetId string, unemployed int, workers int, soldiers int) error

	ReserveResources(id string, planetId string, water int, graphene int) error
	ExtractReservedResources(id string, planetId string, water int, graphene int) error

	Research(id string, planetId string, researchName string,
		grapheneRequired float64, waterRequired float64, shelioRequired float64, minutesRequired float64) error
	ResearchUpgrade(id string, planetId string, researchName string,
		grapheneRequired float64, waterRequired float64, shelioRequired float64, minutesRequired float64) error
	CancelResearch(id string, planetId string, researchName string, grapheneReturned int, waterReturned int, shelioReturned int) error

	ConstructShips(id string, planetId string, unitName string, quantity float64, constructionRequirements models.Requirements) error
	CancelShipsConstruction(id string, planetId string, unitName string, cancelReturns models.Returns) error
	DestructShips(id string, planetId string, unitName string, quantity float64, destructionReturns models.Returns) error

	ConstructDefences(id string, planetId string, unitName string, quantity float64, constructionRequirements models.Requirements) error
	CancelDefencesConstruction(id string, planetId string, unitName string, cancelReturns models.Returns) error
	DestructDefences(id string, planetId string, unitName string, quantity float64, destructionReturns models.Returns) error

	ConstructDefenceShipCarrier(id string, planetId string, unitName string, unitId string, constructionRequirements models.Requirements) error
	CancelDefenceShipCarrierConstruction(id string, planetId string, unitId string, cancelReturns models.Returns) error
	UpgradeDefenceShipCarrier(id string, planetId string, unitId string, constructionRequirements models.Requirements) error
	CancelDefenceShipCarrierUpGradation(id string, planetId string, unitId string, cancelReturns models.Returns) error
	DestructDefenceShipCarrier(id string, planetId string, unitId string, destructionReturns models.Returns) error

	DeployShipsOnDefenceShipCarrier(id string, planetId string, unitId string, ships map[string]int) error
	DeployDefencesOnShield(id string, planetId string, shieldId string, defences map[string]int) error
	DeployDefenceShipCarrierOnShield(id string, planetId string, unitId string, shieldId string) error
}

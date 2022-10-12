package military

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type DeployedDefence struct {
	Name  string `json:"name" example:"BOMBER"`
	Level int    `json:"level" example:"1"`
	Units int    `json:"units" example:"5"`
}

type Defence struct {
	Name                 string                   `json:"name" example:"BOMBER"`
	Level                int                      `json:"level" example:"1"`
	TotalUnits           int                      `json:"total_units" example:"5"`
	IdleUnits            int                      `json:"idle_units" example:"5"`
	DefenceAttributes    models.DefenceAttributes `json:"defence_attributes"`
	CreationRequirements models.Requirements      `json:"creation_requirements"`
	DestructionReturns   models.Returns           `json:"destruction_returns"`
	UnderConstruction    *UnderConstruction       `json:"under_construction,omitempty"`
}

func (d *Defence) Init(unitName string, defenceUser repoModels.Defence, defenceConstants map[string]map[string]interface{}) {
	d.Name = unitName
	d.Level = defenceUser.Level
	d.TotalUnits = defenceUser.Quantity

	d.IdleUnits = repoModels.GetIdleDefences(defenceUser.GuardingShield, defenceUser.Quantity)

	currentLevelString := strconv.Itoa(defenceUser.Level)
	d.DefenceAttributes.Init(defenceConstants[currentLevelString])
	d.CreationRequirements.Init(defenceConstants[currentLevelString])
	d.DestructionReturns.InitDestructionReturns(defenceConstants[currentLevelString])
	d.UnderConstruction = InitUnderConstruction(defenceUser.UnderConstruction, defenceConstants[currentLevelString])
}

type DeployedDefenceShipCarrier struct {
	Id            string         `json:"_id" example:"DSC001"`
	Name          string         `json:"name" example:"VIKRAM"`
	Level         int            `json:"level" example:"1"`
	Deployed      bool           `json:"deployed" example:"true"`
	DeployedShips []DeployedShip `json:"deployed_ships"`
}

type DefenceShipCarrier struct {
	Name                 string                   `json:"name" example:"VIKRAM"`
	CreationRequirements models.Requirements      `json:"creation_requirements"`
	Units                []DefenceShipCarrierUnit `json:"units"`
}

type DefenceShipCarrierUnit struct {
	Id                    string                   `json:"_id" example:"DSC001"`
	Level                 int                      `json:"level" example:"1"`
	DeployedShips         []DeployedShip           `json:"deployed_ships"`
	Idle                  bool                     `json:"idle" example:"true"`
	DefenceAttributes     models.DefenceAttributes `json:"defence_attributes"`
	NextLevelRequirements models.Requirements      `json:"next_level_requirements"`
	DestructionReturns    models.Returns           `json:"destruction_returns"`
	UnderConstruction     *UnderConstruction       `json:"under_construction,omitempty"`
}

func InitAllDefenceShipCarriers(planetUser repoModels.PlanetUser,
	defenceConstants map[string]constants.MilitaryConstants) []DefenceShipCarrier {

	var defenceShipCarriers []DefenceShipCarrier
	for unitName, defenceConstant := range defenceConstants {
		if defenceConstant.Type == constants.DefenceShipCarrier {
			d := DefenceShipCarrier{Name: unitName}
			d.CreationRequirements.Init(defenceConstants[unitName].Levels["1"])
			for _, defenceShipCarrierUser := range planetUser.DefenceShipCarriers {
				if defenceShipCarrierUser.Name == unitName {
					currentLevelString := strconv.Itoa(defenceShipCarrierUser.Level)
					u := DefenceShipCarrierUnit{
						Id:            defenceShipCarrierUser.Id,
						Level:         defenceShipCarrierUser.Level,
						DeployedShips: GetDeployedShips(planetUser, defenceShipCarrierUser.HostingShips),
						Idle:          defenceShipCarrierUser.GuardingShield == "",
					}
					u.DefenceAttributes.Init(defenceConstants[unitName].Levels[currentLevelString])
					u.NextLevelRequirements.InitNextLevelRequirements(defenceShipCarrierUser.Level, defenceConstants[unitName])
					u.DestructionReturns.InitDestructionReturns(defenceConstants[unitName].Levels[currentLevelString])
					u.UnderConstruction = InitUnderUpGradation(defenceShipCarrierUser.UnderConstruction, defenceConstants[unitName].Levels[currentLevelString])
					d.Units = append(d.Units, u)
				}
			}
			defenceShipCarriers = append(defenceShipCarriers, d)
		}
	}
	return defenceShipCarriers
}

func GetDeployedShips(planetUser repoModels.PlanetUser, hostingShips map[string]int) []DeployedShip {
	var deployedShips []DeployedShip
	for unitName, units := range hostingShips {
		s := DeployedShip{
			Name:  unitName,
			Level: planetUser.GetShip(unitName).Level,
			Units: units,
		}
		deployedShips = append(deployedShips, s)
	}
	return deployedShips
}

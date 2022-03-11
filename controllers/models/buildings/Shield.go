package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Shield struct {
	Id                          string                                `json:"_id" example:"SHLD101"`
	Level                       int                                   `json:"level" example:"3"`
	BuildingState               BuildingState                         `json:"building_state"`
	Workers                     int                                   `json:"workers" example:"12"`
	BuildingAttributes          ShieldAttributes                      `json:"building_attributes"`
	NextLevelRequirements       NextLevelRequirements                 `json:"next_level_requirements"`
	DeployedDefences            []military.DeployedDefence            `json:"deployed_defences"`
	DeployedDefenceShipCarriers []military.DeployedDefenceShipCarrier `json:"deployed_defence_ship_carriers"`
}

type ShieldAttributes struct {
	HitPoints IntegerBuildingAttributes `json:"hit_points"`
}

func InitAllShields(planetUser models.PlanetUser,
	defenceConstants map[string]constants.DefenceConstants, shieldBuildingUpgradeConstants constants.UpgradeConstants) []Shield {

	var shields []Shield
	shieldIds := constants.GetShieldIds()
	for _, shieldId := range shieldIds {
		s := Shield{}
		s.Id = shieldId
		s.Level = planetUser.Buildings[shieldId].BuildingLevel
		s.BuildingState.Init(planetUser.Buildings[shieldId], shieldBuildingUpgradeConstants)
		s.Workers = planetUser.Buildings[shieldId].Workers
		s.NextLevelRequirements.Init(planetUser.Buildings[shieldId].BuildingLevel, shieldBuildingUpgradeConstants)
		s.BuildingAttributes.Init(planetUser.Buildings[shieldId].BuildingLevel, defenceConstants[constants.Shield])
		for unitName, defenceUser := range planetUser.Defences {
			if deployedDefences, ok := defenceUser.GuardingShield[shieldId]; ok {
				d := military.DeployedDefence{
					Name:  unitName,
					Level: defenceUser.Level,
					Units: deployedDefences,
				}
				s.DeployedDefences = append(s.DeployedDefences, d)
			}
		}
		for defenceShipCarrierId, defenceShipCarrierUser := range planetUser.DefenceShipCarriers {
			if defenceShipCarrierUser.GuardingShield != "" {
				d := military.DeployedDefenceShipCarrier{
					Id:            defenceShipCarrierId,
					Name:          defenceShipCarrierUser.Name,
					Level:         defenceShipCarrierUser.Level,
					DeployedShips: military.GetDeployedShips(planetUser.Ships, defenceShipCarrierUser.HostingShips),
				}
				s.DeployedDefenceShipCarriers = append(s.DeployedDefenceShipCarriers, d)
			}
		}
		shields = append(shields, s)
	}
	return shields
}

func (n *ShieldAttributes) Init(currentLevel int, shieldConstants constants.DefenceConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(shieldConstants.MaxLevel)
	n.HitPoints.Current = shieldConstants.Levels[currentLevelString].HitPoints
	n.HitPoints.Max = shieldConstants.Levels[maxLevelString].HitPoints
	if currentLevel+1 < shieldConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.HitPoints.Next = shieldConstants.Levels[nextLevelString].HitPoints
	}
}

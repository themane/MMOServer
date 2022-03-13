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
	BuildingState               State                                 `json:"building_state"`
	Workers                     int                                   `json:"workers" example:"12"`
	BuildingAttributes          ShieldAttributes                      `json:"building_attributes"`
	NextLevelRequirements       NextLevelRequirements                 `json:"next_level_requirements"`
	DeployedDefences            []military.DeployedDefence            `json:"deployed_defences"`
	DeployedDefenceShipCarriers []military.DeployedDefenceShipCarrier `json:"deployed_defence_ship_carriers"`
}

type ShieldAttributes struct {
	HitPoints FloatBuildingAttributes `json:"hit_points"`
}

func InitAllShields(planetUser models.PlanetUser,
	shieldConstants map[string]map[string]interface{}, shieldBuildingUpgradeConstants constants.UpgradeConstants) []Shield {

	var shields []Shield
	shieldIds := constants.GetShieldIds()
	for shieldId := range shieldIds {
		s := Shield{}
		s.Id = shieldId
		s.Level = planetUser.Buildings[shieldId].BuildingLevel
		s.BuildingState.Init(planetUser.Buildings[shieldId], shieldBuildingUpgradeConstants)
		s.Workers = planetUser.Buildings[shieldId].Workers
		s.NextLevelRequirements.Init(planetUser.Buildings[shieldId].BuildingLevel, shieldBuildingUpgradeConstants)
		s.BuildingAttributes.Init(planetUser.Buildings[shieldId].BuildingLevel, shieldBuildingUpgradeConstants.MaxLevel, shieldConstants)
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
			if defenceShipCarrierUser.GuardingShield == shieldId {
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

func (n *ShieldAttributes) Init(currentLevel int, maxLevel int, shieldConstants map[string]map[string]interface{}) {
	maxLevelString := strconv.Itoa(maxLevel)
	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		n.HitPoints.Current = shieldConstants[currentLevelString]["hit_points"].(float64)
	}
	n.HitPoints.Max = shieldConstants[maxLevelString]["hit_points"].(float64)
	if currentLevel+1 < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.HitPoints.Next = shieldConstants[nextLevelString]["hit_points"].(float64)
	}
}

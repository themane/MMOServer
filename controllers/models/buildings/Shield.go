package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Shield struct {
	Id                          string                                `json:"_id" example:"SHLD101"`
	Name                        string                                `json:"name" example:"shield"`
	Level                       int                                   `json:"level" example:"3"`
	BuildingState               models.State                          `json:"building_state"`
	Workers                     int                                   `json:"workers" example:"12"`
	BuildingAttributes          ShieldAttributes                      `json:"building_attributes"`
	NextLevelRequirements       *models.NextLevelRequirements         `json:"next_level_requirements"`
	DeployedDefences            []military.DeployedDefence            `json:"deployed_defences"`
	DeployedDefenceShipCarriers []military.DeployedDefenceShipCarrier `json:"deployed_defence_ship_carriers"`
}

type ShieldAttributes struct {
	HitPoints       FloatBuildingAttributes `json:"hit_points"`
	WorkersMaxLimit FloatBuildingAttributes `json:"workers_max_limit"`
}

func InitAllShields(planetUser models.PlanetUser,
	shieldConstants map[string]map[string]interface{}, shieldBuildingUpgradeConstants constants.UpgradeConstants) []Shield {

	var shields []Shield
	shieldIds := constants.GetShieldIds()
	for shieldId := range shieldIds {
		s := Shield{}
		s.Id = shieldId
		s.Name = constants.Shield
		shield := planetUser.GetBuilding(shieldId)
		s.Level = shield.BuildingLevel
		s.BuildingState.Init(*shield, shieldBuildingUpgradeConstants)
		s.Workers = shield.Workers
		s.BuildingAttributes.Init(shield.BuildingLevel, shieldBuildingUpgradeConstants.MaxLevel, shieldConstants)
		if s.Level < shieldBuildingUpgradeConstants.MaxLevel {
			s.NextLevelRequirements = &models.NextLevelRequirements{}
			s.NextLevelRequirements.Init(shield.BuildingLevel, shieldBuildingUpgradeConstants)
		}
		for _, defenceUser := range planetUser.Defences {
			if deployedDefences, ok := defenceUser.GuardingShield[shieldId]; ok {
				d := military.DeployedDefence{
					Name:  defenceUser.Name,
					Level: defenceUser.Level,
					Units: deployedDefences,
				}
				s.DeployedDefences = append(s.DeployedDefences, d)
			}
		}
		for _, defenceShipCarrierUser := range planetUser.DefenceShipCarriers {
			if defenceShipCarrierUser.GuardingShield == "" || defenceShipCarrierUser.GuardingShield == shieldId {
				d := military.DeployedDefenceShipCarrier{
					Id:            defenceShipCarrierUser.Id,
					Name:          defenceShipCarrierUser.Name,
					Level:         defenceShipCarrierUser.Level,
					Deployed:      defenceShipCarrierUser.GuardingShield == shieldId,
					DeployedShips: military.GetDeployedShips(planetUser, defenceShipCarrierUser.HostingShips),
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
		n.WorkersMaxLimit.Current = shieldConstants[currentLevelString]["workers_max_limit"].(float64)
	}
	n.HitPoints.Max = shieldConstants[maxLevelString]["hit_points"].(float64)
	n.WorkersMaxLimit.Max = shieldConstants[maxLevelString]["workers_max_limit"].(float64)
	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.HitPoints.Next = shieldConstants[nextLevelString]["hit_points"].(float64)
		n.WorkersMaxLimit.Next = shieldConstants[nextLevelString]["workers_max_limit"].(float64)
	} else {
		n.HitPoints.Next = shieldConstants[maxLevelString]["hit_points"].(float64)
		n.WorkersMaxLimit.Next = shieldConstants[maxLevelString]["workers_max_limit"].(float64)
	}
}

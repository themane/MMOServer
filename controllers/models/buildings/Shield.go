package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/military"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Shield struct {
	Id                    string                `json:"_id" example:"SHLD101"`
	Level                 int                   `json:"level" example:"3"`
	BuildingState         BuildingState         `json:"building_state"`
	Workers               int                   `json:"workers" example:"12"`
	BuildingAttributes    ShieldAttributes      `json:"building_attributes"`
	NextLevelRequirements NextLevelRequirements `json:"next_level_requirements"`
	DeployedDefences      []military.Defence    `json:"deployed_defences"`
}

type ShieldAttributes struct {
	HitPoints IntegerBuildingAttributes `json:"hit_points" `
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
		for defenceType, defenceUser := range planetUser.Defences {
			if deployedDefences, ok := defenceUser.GuardingShield[shieldId]; ok {
				d := military.Defence{}
				d.Init(defenceType, deployedDefences, defenceUser, defenceConstants[defenceType])
				s.DeployedDefences = append(s.DeployedDefences, d)
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

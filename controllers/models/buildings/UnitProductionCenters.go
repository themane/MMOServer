package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type UnitProductionCenter struct {
	BuildingId            string                           `json:"building_id" example:"ATTACK_PRODUCTION_CENTER"`
	Level                 int                              `json:"level" example:"3"`
	Workers               int                              `json:"workers" example:"12"`
	Soldiers              int                              `json:"soldiers" example:"15"`
	BuildingState         repoModels.State                 `json:"building_state"`
	BuildingAttributes    UnitProductionCenterAttributes   `json:"building_attributes"`
	NextLevelRequirements repoModels.NextLevelRequirements `json:"next_level_requirements"`
}
type UnitProductionCenterAttributes struct {
	WorkerBonus             map[string]FloatBuildingAttributes `json:"worker_bonus"`
	SoldierBonus            map[string]FloatBuildingAttributes `json:"soldier_bonus"`
	MinimumWorkersRequired  FloatBuildingAttributes            `json:"minimum_workers_required"`
	MinimumSoldiersRequired FloatBuildingAttributes            `json:"minimum_soldiers_required"`
	WorkersMaxLimit         FloatBuildingAttributes            `json:"workers_max_limit"`
	SoldiersMaxLimit        FloatBuildingAttributes            `json:"soldiers_max_limit"`
}

func InitAttackProductionCenter(planetUser repoModels.PlanetUser,
	attackProductionCenterUpgradeConstants constants.UpgradeConstants,
	attackProductionCenterBuildingConstants map[string]map[string]interface{}) *UnitProductionCenter {

	u := new(UnitProductionCenter)
	u.BuildingId = constants.AttackProductionCenter
	u.Level = planetUser.Buildings[constants.AttackProductionCenter].BuildingLevel
	u.Workers = planetUser.Buildings[constants.AttackProductionCenter].Workers
	u.Soldiers = planetUser.Buildings[constants.AttackProductionCenter].Soldiers
	u.BuildingState.Init(planetUser.Buildings[constants.AttackProductionCenter], attackProductionCenterUpgradeConstants)
	u.NextLevelRequirements.Init(planetUser.Buildings[constants.AttackProductionCenter].BuildingLevel, attackProductionCenterUpgradeConstants)
	u.BuildingAttributes.Init(planetUser.Buildings[constants.AttackProductionCenter].BuildingLevel,
		attackProductionCenterUpgradeConstants.MaxLevel, attackProductionCenterBuildingConstants, constants.GetShipAttributes())
	return u
}

func InitDefenceProductionCenter(planetUser repoModels.PlanetUser,
	defenceProductionCenterUpgradeConstants constants.UpgradeConstants,
	defenceProductionCenterBuildingConstants map[string]map[string]interface{}) *UnitProductionCenter {

	u := new(UnitProductionCenter)
	u.BuildingId = constants.DefenceProductionCenter
	u.Level = planetUser.Buildings[constants.DefenceProductionCenter].BuildingLevel
	u.Workers = planetUser.Buildings[constants.DefenceProductionCenter].Workers
	u.Soldiers = planetUser.Buildings[constants.DefenceProductionCenter].Soldiers
	u.BuildingState.Init(planetUser.Buildings[constants.DefenceProductionCenter], defenceProductionCenterUpgradeConstants)
	u.NextLevelRequirements.Init(planetUser.Buildings[constants.DefenceProductionCenter].BuildingLevel, defenceProductionCenterUpgradeConstants)
	u.BuildingAttributes.Init(planetUser.Buildings[constants.DefenceProductionCenter].BuildingLevel,
		defenceProductionCenterUpgradeConstants.MaxLevel, defenceProductionCenterBuildingConstants, constants.GetDefenceAttributes())
	return u
}

func (a *UnitProductionCenterAttributes) Init(currentLevel int, maxLevel int,
	buildingConstants map[string]map[string]interface{}, bonusAttrs []string) {

	var currentWorkerBonus map[string]interface{}
	var currentSoldierBonus map[string]interface{}
	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		currentWorkerBonus = buildingConstants[currentLevelString]["workers_bonus"].(map[string]interface{})
		currentSoldierBonus = buildingConstants[currentLevelString]["soldiers_bonus"].(map[string]interface{})
		a.MinimumWorkersRequired.Current = buildingConstants[currentLevelString]["workers_required"].(float64)
		a.MinimumSoldiersRequired.Current = buildingConstants[currentLevelString]["soldiers_required"].(float64)
		a.WorkersMaxLimit.Current = buildingConstants[currentLevelString]["workers_max_limit"].(float64)
		a.SoldiersMaxLimit.Current = buildingConstants[currentLevelString]["soldiers_max_limit"].(float64)
	}
	maxLevelString := strconv.Itoa(maxLevel)
	maxWorkerBonus := buildingConstants[maxLevelString]["workers_bonus"].(map[string]interface{})
	maxSoldierBonus := buildingConstants[maxLevelString]["soldiers_bonus"].(map[string]interface{})
	a.MinimumWorkersRequired.Max = buildingConstants[maxLevelString]["workers_required"].(float64)
	a.MinimumSoldiersRequired.Max = buildingConstants[maxLevelString]["soldiers_required"].(float64)
	a.WorkersMaxLimit.Max = buildingConstants[maxLevelString]["workers_max_limit"].(float64)
	a.SoldiersMaxLimit.Max = buildingConstants[maxLevelString]["soldiers_max_limit"].(float64)

	var nextWorkerBonus map[string]interface{}
	var nextSoldierBonus map[string]interface{}
	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		nextWorkerBonus = buildingConstants[nextLevelString]["workers_bonus"].(map[string]interface{})
		nextSoldierBonus = buildingConstants[nextLevelString]["soldiers_bonus"].(map[string]interface{})
		a.MinimumWorkersRequired.Next = buildingConstants[nextLevelString]["workers_required"].(float64)
		a.MinimumSoldiersRequired.Next = buildingConstants[nextLevelString]["soldiers_required"].(float64)
		a.WorkersMaxLimit.Next = buildingConstants[nextLevelString]["workers_max_limit"].(float64)
		a.SoldiersMaxLimit.Next = buildingConstants[nextLevelString]["soldiers_max_limit"].(float64)
	}
	a.WorkerBonus = map[string]FloatBuildingAttributes{}
	a.initBonus(a.WorkerBonus, currentWorkerBonus, nextWorkerBonus, maxWorkerBonus, bonusAttrs)
	a.SoldierBonus = map[string]FloatBuildingAttributes{}
	a.initBonus(a.SoldierBonus, currentSoldierBonus, nextSoldierBonus, maxSoldierBonus, bonusAttrs)
}

func (a *UnitProductionCenterAttributes) initBonus(bonus map[string]FloatBuildingAttributes,
	currentBonus map[string]interface{}, nextBonus map[string]interface{}, maxBonus map[string]interface{},
	attrs []string,
) {
	for _, attr := range attrs {
		if value, ok := currentBonus[attr]; ok {
			attrValues := FloatBuildingAttributes{
				Current: value.(float64),
				Max:     maxBonus[attr].(float64),
			}
			if len(nextBonus) > 0 {
				attrValues.Next = nextBonus[attr].(float64)
			}
			bonus[attr] = attrValues
		}
	}
}

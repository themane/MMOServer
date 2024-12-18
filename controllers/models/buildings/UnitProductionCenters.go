package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type UnitProductionCenter struct {
	BuildingId            string                            `json:"building_id" example:"ATTACK_PRODUCTION_CENTER"`
	Level                 int                               `json:"level" example:"3"`
	Workers               int                               `json:"workers" example:"12"`
	Soldiers              int                               `json:"soldiers" example:"15"`
	BuildingState         repoModels.State                  `json:"building_state"`
	BuildingAttributes    UnitProductionCenterAttributes    `json:"building_attributes"`
	NextLevelRequirements *repoModels.NextLevelRequirements `json:"next_level_requirements"`
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
	if planetUser.GetBuilding(constants.AttackProductionCenter) != nil {
		u.Level = planetUser.GetBuilding(constants.AttackProductionCenter).BuildingLevel
		u.Workers = planetUser.GetBuilding(constants.AttackProductionCenter).Workers
		u.Soldiers = planetUser.GetBuilding(constants.AttackProductionCenter).Soldiers
	}
	u.BuildingState.Init(planetUser.GetBuilding(constants.AttackProductionCenter), attackProductionCenterUpgradeConstants)
	u.BuildingAttributes.Init(u.Level, attackProductionCenterUpgradeConstants.MaxLevel,
		attackProductionCenterBuildingConstants, constants.GetShipAttributes())
	if u.Level < attackProductionCenterUpgradeConstants.MaxLevel {
		u.NextLevelRequirements = &repoModels.NextLevelRequirements{}
		u.NextLevelRequirements.Init(u.Level, attackProductionCenterUpgradeConstants)
	}
	return u
}

func InitDefenceProductionCenter(planetUser repoModels.PlanetUser,
	defenceProductionCenterUpgradeConstants constants.UpgradeConstants,
	defenceProductionCenterBuildingConstants map[string]map[string]interface{}) *UnitProductionCenter {

	u := new(UnitProductionCenter)
	u.BuildingId = constants.DefenceProductionCenter
	if planetUser.GetBuilding(constants.DefenceProductionCenter) != nil {
		u.Level = planetUser.GetBuilding(constants.DefenceProductionCenter).BuildingLevel
		u.Workers = planetUser.GetBuilding(constants.DefenceProductionCenter).Workers
		u.Soldiers = planetUser.GetBuilding(constants.DefenceProductionCenter).Soldiers
	}
	u.BuildingState.Init(planetUser.GetBuilding(constants.DefenceProductionCenter), defenceProductionCenterUpgradeConstants)
	u.BuildingAttributes.Init(u.Level, defenceProductionCenterUpgradeConstants.MaxLevel,
		defenceProductionCenterBuildingConstants, constants.GetDefenceAttributes())
	if u.Level < defenceProductionCenterUpgradeConstants.MaxLevel {
		u.NextLevelRequirements = &repoModels.NextLevelRequirements{}
		u.NextLevelRequirements.Init(u.Level, defenceProductionCenterUpgradeConstants)
	}
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
	} else {
		nextWorkerBonus = buildingConstants[maxLevelString]["workers_bonus"].(map[string]interface{})
		nextSoldierBonus = buildingConstants[maxLevelString]["soldiers_bonus"].(map[string]interface{})
		a.MinimumWorkersRequired.Next = buildingConstants[maxLevelString]["workers_required"].(float64)
		a.MinimumSoldiersRequired.Next = buildingConstants[maxLevelString]["soldiers_required"].(float64)
		a.WorkersMaxLimit.Next = buildingConstants[maxLevelString]["workers_max_limit"].(float64)
		a.SoldiersMaxLimit.Next = buildingConstants[maxLevelString]["soldiers_max_limit"].(float64)
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

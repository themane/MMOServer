package buildings

import (
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type ResourceStorage struct {
	BuildingId            string                           `json:"building_id" example:"WATER_PRESSURE_TANK"`
	Level                 int                              `json:"level" example:"3"`
	Workers               int                              `json:"workers" example:"12"`
	Soldiers              int                              `json:"soldiers" example:"15"`
	BuildingState         repoModels.State                 `json:"building_state"`
	BuildingAttributes    ResourceStorageAttributes        `json:"building_attributes"`
	NextLevelRequirements repoModels.NextLevelRequirements `json:"next_level_requirements"`
}
type ResourceStorageAttributes struct {
	StorageRatePerWorker    FloatBuildingAttributes `json:"storage_rate_per_worker"`
	MaxStoragePerSoldier    FloatBuildingAttributes `json:"max_storage_per_soldier"`
	MinimumWorkersRequired  FloatBuildingAttributes `json:"minimum_workers_required"`
	MinimumSoldiersRequired FloatBuildingAttributes `json:"minimum_soldiers_required"`
	WorkersMaxLimit         FloatBuildingAttributes `json:"workers_max_limit"`
	SoldiersMaxLimit        FloatBuildingAttributes `json:"soldiers_max_limit"`
}

func InitWaterPressureTank(planetUser repoModels.PlanetUser,
	upgradeConstants constants.UpgradeConstants,
	buildingConstants map[string]map[string]interface{}) *ResourceStorage {

	r := new(ResourceStorage)
	r.BuildingId = constants.WaterPressureTank
	waterPressureTank := planetUser.GetBuilding(constants.WaterPressureTank)
	r.Level = waterPressureTank.BuildingLevel
	r.Workers = waterPressureTank.Workers
	r.Soldiers = waterPressureTank.Soldiers
	r.BuildingState.Init(*waterPressureTank, upgradeConstants)
	r.NextLevelRequirements.Init(waterPressureTank.BuildingLevel, upgradeConstants)
	r.BuildingAttributes.Init(waterPressureTank.BuildingLevel, upgradeConstants.MaxLevel, buildingConstants)
	return r
}

func InitDiamondStorage(planetUser repoModels.PlanetUser,
	upgradeConstants constants.UpgradeConstants,
	buildingConstants map[string]map[string]interface{}) *ResourceStorage {

	r := new(ResourceStorage)
	r.BuildingId = constants.DiamondStorage
	diamondStorage := planetUser.GetBuilding(constants.DiamondStorage)
	r.Level = diamondStorage.BuildingLevel
	r.Workers = diamondStorage.Workers
	r.Soldiers = diamondStorage.Soldiers
	r.BuildingState.Init(*diamondStorage, upgradeConstants)
	r.NextLevelRequirements.Init(diamondStorage.BuildingLevel, upgradeConstants)
	r.BuildingAttributes.Init(diamondStorage.BuildingLevel, upgradeConstants.MaxLevel, buildingConstants)
	return r
}

func (a *ResourceStorageAttributes) Init(currentLevel int, maxLevel int, buildingConstants map[string]map[string]interface{}) {
	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		a.StorageRatePerWorker.Current = buildingConstants[currentLevelString]["storage_rate_per_worker"].(float64)
		a.MaxStoragePerSoldier.Current = buildingConstants[currentLevelString]["max_storage_per_soldier"].(float64)
		a.MinimumWorkersRequired.Current = buildingConstants[currentLevelString]["workers_required"].(float64)
		a.MinimumSoldiersRequired.Current = buildingConstants[currentLevelString]["soldiers_required"].(float64)
		a.WorkersMaxLimit.Current = buildingConstants[currentLevelString]["workers_max_limit"].(float64)
		a.SoldiersMaxLimit.Current = buildingConstants[currentLevelString]["soldiers_max_limit"].(float64)
	}
	maxLevelString := strconv.Itoa(maxLevel)
	a.StorageRatePerWorker.Max = buildingConstants[maxLevelString]["storage_rate_per_worker"].(float64)
	a.MaxStoragePerSoldier.Max = buildingConstants[maxLevelString]["max_storage_per_soldier"].(float64)
	a.MinimumWorkersRequired.Max = buildingConstants[maxLevelString]["workers_required"].(float64)
	a.MinimumSoldiersRequired.Max = buildingConstants[maxLevelString]["soldiers_required"].(float64)
	a.WorkersMaxLimit.Max = buildingConstants[maxLevelString]["workers_max_limit"].(float64)
	a.SoldiersMaxLimit.Max = buildingConstants[maxLevelString]["soldiers_max_limit"].(float64)

	if currentLevel < maxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		a.StorageRatePerWorker.Next = buildingConstants[nextLevelString]["storage_rate_per_worker"].(float64)
		a.MaxStoragePerSoldier.Next = buildingConstants[nextLevelString]["max_storage_per_soldier"].(float64)
		a.MinimumWorkersRequired.Next = buildingConstants[nextLevelString]["workers_required"].(float64)
		a.MinimumSoldiersRequired.Next = buildingConstants[nextLevelString]["soldiers_required"].(float64)
		a.WorkersMaxLimit.Next = buildingConstants[nextLevelString]["workers_max_limit"].(float64)
		a.SoldiersMaxLimit.Next = buildingConstants[nextLevelString]["soldiers_max_limit"].(float64)
	}
}

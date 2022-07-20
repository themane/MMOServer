package buildings

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Mine struct {
	Id          string      `json:"_id" example:"W101"`
	Type        string      `json:"type" example:"WATER"`
	MaxLimit    int         `json:"max_limit" example:"550"`
	Mined       int         `json:"mined" example:"125"`
	MiningPlant MiningPlant `json:"mining_plant"`
}

type MiningPlant struct {
	BuildingId            string                        `json:"building_id,omitempty" example:"WMP101"`
	Level                 int                           `json:"level" example:"3"`
	Workers               int                           `json:"workers" example:"12"`
	BuildingState         models.State                  `json:"building_state"`
	BuildingAttributes    MiningAttributes              `json:"building_attributes"`
	NextLevelRequirements *models.NextLevelRequirements `json:"next_level_requirements"`
}

type MiningAttributes struct {
	MiningRatePerWorker IntegerBuildingAttributes `json:"mining_rate_per_worker"`
	WorkersMaxLimit     IntegerBuildingAttributes `json:"workers_max_limit"`
}

func InitAllMines(planetUni models.PlanetUni, planetUser models.PlanetUser,
	waterMiningPlantConstants constants.UpgradeConstants, grapheneMiningPlantConstants constants.UpgradeConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants) []Mine {

	var mines []Mine
	for mineId := range planetUser.Mines {
		m := Mine{}
		m.Init(planetUni.Mines[mineId], planetUser,
			waterMiningPlantConstants, grapheneMiningPlantConstants,
			waterConstants, grapheneConstants)
		mines = append(mines, m)
	}
	return mines
}

func (m *Mine) Init(mineUni models.MineUni, planetUser models.PlanetUser,
	waterMiningPlantConstants constants.UpgradeConstants, grapheneMiningPlantConstants constants.UpgradeConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants) {

	m.Id = mineUni.Id
	m.Type = mineUni.Type
	m.MaxLimit = mineUni.MaxLimit
	m.Mined = planetUser.GetMine(mineUni.Id).Mined
	if mineUni.Type == constants.Water {
		m.MiningPlant.Init(planetUser, mineUni.Id, waterConstants, waterMiningPlantConstants)
	}
	if mineUni.Type == constants.Graphene {
		m.MiningPlant.Init(planetUser, mineUni.Id, grapheneConstants, grapheneMiningPlantConstants)
	}
}

func (m *MiningPlant) Init(planetUser models.PlanetUser, mineId string,
	miningConstants constants.MiningConstants, miningPlantUpgradeConstants constants.UpgradeConstants) {

	m.BuildingId = models.GetMiningPlantId(mineId)
	m.Level = planetUser.GetBuilding(m.BuildingId).BuildingLevel
	m.Workers = planetUser.GetBuilding(m.BuildingId).Workers
	m.BuildingState.Init(*planetUser.GetBuilding(m.BuildingId), miningPlantUpgradeConstants)
	m.BuildingAttributes.Init(m.Level, miningConstants)
	if m.Level < miningPlantUpgradeConstants.MaxLevel {
		m.NextLevelRequirements = &models.NextLevelRequirements{}
		m.NextLevelRequirements.Init(m.Level, miningPlantUpgradeConstants)
	}
}

func (n *MiningAttributes) Init(currentLevel int, miningConstants constants.MiningConstants) {
	if currentLevel > 0 {
		currentLevelString := strconv.Itoa(currentLevel)
		n.WorkersMaxLimit.Current = miningConstants.Levels[currentLevelString].WorkersMaxLimit
		n.MiningRatePerWorker.Current = miningConstants.Levels[currentLevelString].MiningRatePerWorker
	}
	maxLevelString := strconv.Itoa(miningConstants.MaxLevel)
	n.MiningRatePerWorker.Max = miningConstants.Levels[maxLevelString].MiningRatePerWorker
	n.WorkersMaxLimit.Max = miningConstants.Levels[maxLevelString].WorkersMaxLimit
	if currentLevel < miningConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.WorkersMaxLimit.Next = miningConstants.Levels[nextLevelString].WorkersMaxLimit
		n.MiningRatePerWorker.Next = miningConstants.Levels[nextLevelString].MiningRatePerWorker
	}
}

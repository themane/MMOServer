package models

import "strconv"

type Mine struct {
	Id          string      `json:"_id" example:"W101"`
	Type        string      `json:"type" example:"WATER"`
	MaxLimit    int         `json:"max_limit" example:"550"`
	Mined       int         `json:"mined" example:"125"`
	MiningPlant MiningPlant `json:"mining_plant"`
}

type MiningPlant struct {
	BuildingId          string              `json:"building_id,omitempty" example:"WMP101"`
	BuildingLevel       int                 `json:"building_level" example:"3"`
	Workers             int                 `json:"workers" example:"12"`
	NextLevelAttributes NextLevelAttributes `json:"next_level"`
}

type NextLevelAttributes struct {
	GrapheneRequired           int `json:"graphene_required" example:"101"`
	WaterRequired              int `json:"water_required" example:"5"`
	ShelioRequired             int `json:"shelio_required" example:"0"`
	CurrentMiningRatePerWorker int `json:"current_mining_rate_per_worker" example:"1"`
	NextMiningRatePerWorker    int `json:"next_mining_rate_per_worker" example:"1"`
	MaxMiningRatePerWorker     int `json:"max_mining_rate_per_worker" example:"12"`
	CurrentWorkersMaxLimit     int `json:"current_workers_max_limit" example:"40"`
	NextWorkersMaxLimit        int `json:"next_workers_max_limit" example:"65"`
	MaxWorkersMaxLimit         int `json:"max_workers_max_limit" example:"130"`
}

func (m *Mine) Init(mineUni MineUni, mineUser MineUser, waterConstants ResourceConstants, grapheneConstants ResourceConstants) {
	m.Id = mineUni.Id
	m.Type = mineUni.Type
	m.MaxLimit = mineUni.MaxLimit
	m.Mined = mineUser.Mined
	if mineUni.Type == WATER {
		m.MiningPlant.Init(mineUser.MiningPlant, waterConstants)
	}
	if mineUni.Type == GRAPHENE {
		m.MiningPlant.Init(mineUser.MiningPlant, grapheneConstants)
	}
}

func (m *MiningPlant) Init(miningPlantUser MiningPlantUser, resourceConstants ResourceConstants) {
	m.BuildingId = miningPlantUser.BuildingId
	m.BuildingLevel = miningPlantUser.BuildingLevel
	m.Workers = miningPlantUser.Workers
	m.NextLevelAttributes.Init(miningPlantUser.BuildingLevel, resourceConstants)
}

func (n *NextLevelAttributes) Init(currentLevel int, resourceConstants ResourceConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(resourceConstants.MaxLevel)
	n.CurrentWorkersMaxLimit = resourceConstants.Levels[currentLevelString].WorkersMaxLimit
	n.CurrentMiningRatePerWorker = resourceConstants.Levels[currentLevelString].MiningRatePerWorker
	n.MaxMiningRatePerWorker = resourceConstants.Levels[maxLevelString].MiningRatePerWorker
	n.MaxWorkersMaxLimit = resourceConstants.Levels[maxLevelString].WorkersMaxLimit
	if currentLevel+1 < resourceConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.NextWorkersMaxLimit = resourceConstants.Levels[nextLevelString].WorkersMaxLimit
		n.NextMiningRatePerWorker = resourceConstants.Levels[nextLevelString].MiningRatePerWorker
		n.GrapheneRequired = resourceConstants.Levels[nextLevelString].GrapheneRequired
		n.WaterRequired = resourceConstants.Levels[nextLevelString].WaterRequired
		n.ShelioRequired = resourceConstants.Levels[nextLevelString].ShelioRequired
	}
}

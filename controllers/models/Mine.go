package models

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
	BuildingId                string                    `json:"building_id,omitempty" example:"WMP101"`
	Level                     int                       `json:"level" example:"3"`
	Workers                   int                       `json:"workers" example:"12"`
	BuildingState             BuildingState             `json:"building_state"`
	NextLevelMiningAttributes NextLevelMiningAttributes `json:"next_level_attributes"`
	NextLevelRequirements     NextLevelRequirements     `json:"next_level_requirements"`
}

type NextLevelMiningAttributes struct {
	CurrentMiningRatePerWorker int `json:"current_mining_rate_per_worker" example:"1"`
	NextMiningRatePerWorker    int `json:"next_mining_rate_per_worker" example:"1"`
	MaxMiningRatePerWorker     int `json:"max_mining_rate_per_worker" example:"12"`
	CurrentWorkersMaxLimit     int `json:"current_workers_max_limit" example:"40"`
	NextWorkersMaxLimit        int `json:"next_workers_max_limit" example:"65"`
	MaxWorkersMaxLimit         int `json:"max_workers_max_limit" example:"130"`
}

func (m *Mine) Init(mineUni models.MineUni, planetUser models.PlanetUser,
	waterMiningPlantConstants constants.UpgradeConstants, grapheneMiningPlantConstants constants.UpgradeConstants,
	waterConstants constants.MiningConstants, grapheneConstants constants.MiningConstants) {

	m.Id = mineUni.Id
	m.Type = mineUni.Type
	m.MaxLimit = mineUni.MaxLimit
	m.Mined = planetUser.Mines[mineUni.Id].Mined
	if mineUni.Type == constants.Water {
		m.MiningPlant.Init(planetUser, mineUni.Id, waterConstants, waterMiningPlantConstants)
	}
	if mineUni.Type == constants.Graphene {
		m.MiningPlant.Init(planetUser, mineUni.Id, grapheneConstants, grapheneMiningPlantConstants)
	}
}

func (m *MiningPlant) Init(planetUser models.PlanetUser, mineId string,
	miningConstants constants.MiningConstants, miningPlantUpgradeConstants constants.UpgradeConstants) {

	m.BuildingId = planetUser.Mines[mineId].MiningPlantId
	m.Level = planetUser.Buildings[m.BuildingId].BuildingLevel
	m.Workers = planetUser.Buildings[m.BuildingId].Workers
	m.BuildingState.Init(planetUser.Buildings[m.BuildingId], miningPlantUpgradeConstants)
	m.NextLevelMiningAttributes.Init(m.Level, miningConstants)
	m.NextLevelRequirements.Init(m.Level, miningPlantUpgradeConstants)
}

func (n *NextLevelMiningAttributes) Init(currentLevel int, miningConstants constants.MiningConstants) {
	currentLevelString := strconv.Itoa(currentLevel)
	maxLevelString := strconv.Itoa(miningConstants.MaxLevel)
	n.CurrentWorkersMaxLimit = miningConstants.Levels[currentLevelString].WorkersMaxLimit
	n.CurrentMiningRatePerWorker = miningConstants.Levels[currentLevelString].MiningRatePerWorker
	n.MaxMiningRatePerWorker = miningConstants.Levels[maxLevelString].MiningRatePerWorker
	n.MaxWorkersMaxLimit = miningConstants.Levels[maxLevelString].WorkersMaxLimit
	if currentLevel+1 < miningConstants.MaxLevel {
		nextLevelString := strconv.Itoa(currentLevel + 1)
		n.NextWorkersMaxLimit = miningConstants.Levels[nextLevelString].WorkersMaxLimit
		n.NextMiningRatePerWorker = miningConstants.Levels[nextLevelString].MiningRatePerWorker
	}
}

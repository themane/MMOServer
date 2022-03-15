package constants

import (
	"errors"
	"strings"
)

func GetBuildingType(buildingId string) (string, error) {
	if strings.HasPrefix(buildingId, "WMP") {
		return WaterMiningPlant, nil
	}
	if strings.HasPrefix(buildingId, "GMP") {
		return GrapheneMiningPlant, nil
	}
	if strings.HasPrefix(buildingId, "SHLD") {
		return Shield, nil
	}
	if _, ok := GetUpgradableBuildingIds()[buildingId]; ok {
		return buildingId, nil
	}
	return "", errors.New("error. invalid building id: " + buildingId)
}

func GetShieldIds() map[string]struct{} {
	return map[string]struct{}{"SHLD01": {}, "SHLD02": {}, "SHLD03": {}}
}

func GetUpgradableBuildingIds() map[string]struct{} {
	return map[string]struct{}{
		PopulationControlCenter: {},
		AttackProductionCenter:  {}, DefenceProductionCenter: {},
		DiamondStorage: {}, WaterPressureTank: {},
		ResearchLab: {},
	}
}

func GetSoldiersSupportedBuildingIds() map[string]struct{} {
	return map[string]struct{}{
		AttackProductionCenter: {}, DefenceProductionCenter: {},
		DiamondStorage: {}, WaterPressureTank: {},
	}
}

func IsShieldId(id string) bool {
	return id == "SHLD01" || id == "SHLD02" || id == "SHLD03"
}

func GetAttackPointIds() []string {
	return []string{"point1", "point2", "point3", "point4", "point5", "point6"}
}

func GetAttackLineIds() []string {
	return []string{"line1", "line2", "line3", "line4"}
}

func GetShipAttributes() []string {
	return []string{"hit_points", "armor", "resource_capacity", "worker_capacity", "attack", "range", "speed"}
}
func GetDefenceAttributes() []string {
	return []string{"hit_points", "armor", "attack", "range"}
}

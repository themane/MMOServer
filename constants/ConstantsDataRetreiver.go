package constants

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetExperienceConstants() map[string]ExperienceConstants {
	var experienceConstants map[string]ExperienceConstants
	constantsFile, _ := os.Open("resources/ExperienceConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &experienceConstants)
	if err != nil {
		log.Fatal("Error in initializing experience constants: ", err)
		return nil
	}
	return experienceConstants
}

func GetUpgradeConstants() map[string]UpgradeConstants {
	var upgradeConstants map[string]UpgradeConstants
	constantsFile, _ := os.Open("resources/UpgradeConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &upgradeConstants)
	if err != nil {
		log.Fatal("Error in initializing upgrade constants: ", err)
		return nil
	}
	return upgradeConstants
}

func GetBuildingConstants() map[string]map[string]map[string]interface{} {
	var buildingConstants map[string]map[string]map[string]interface{}
	constantsFile, _ := os.Open("resources/BuildingConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &buildingConstants)
	if err != nil {
		log.Fatal("Error in initializing building constants: ", err)
		return nil
	}
	return buildingConstants
}

func GetMiningConstants() map[string]MiningConstants {
	var mineConstants map[string]MiningConstants
	constantsFile, _ := os.Open("resources/MiningConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &mineConstants)
	if err != nil {
		log.Fatal("Error in initializing mining constants: ", err)
		return nil
	}
	return mineConstants
}

func GetMilitaryConstants() map[string]MilitaryConstants {
	var defenceConstants map[string]MilitaryConstants
	constantsFile, _ := os.Open("resources/MilitaryConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &defenceConstants)
	if err != nil {
		log.Fatal("Error in initializing military constants: ", err)
		return nil
	}
	return defenceConstants
}

func GetSpeciesConstants() map[string]SpeciesConstants {
	var speciesConstants map[string]SpeciesConstants
	constantsFile, _ := os.Open("resources/SpeciesConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &speciesConstants)
	if err != nil {
		log.Fatal("Error in initializing species constants: ", err)
		return nil
	}
	return speciesConstants
}

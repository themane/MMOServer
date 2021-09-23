package constants

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetExperienceConstants() ExperienceConstants {
	var experienceConstants ExperienceConstants
	constantsFile, _ := os.Open("resources/ExperienceConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &experienceConstants)
	if err != nil {
		log.Fatal("Error in initializing experience constants: ", err)
		return ExperienceConstants{}
	}
	return experienceConstants
}

func GetBuildingConstants() map[string]BuildingConstants {
	var mineConstants map[string]BuildingConstants
	constantsFile, _ := os.Open("resources/BuildingConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &mineConstants)
	if err != nil {
		log.Fatal("Error in initializing building constants: ", err)
		return nil
	}
	return mineConstants
}

func GetMineConstants() map[string]MiningConstants {
	var mineConstants map[string]MiningConstants
	constantsFile, _ := os.Open("resources/MineConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &mineConstants)
	if err != nil {
		log.Fatal("Error in initializing mining constants: ", err)
		return nil
	}
	return mineConstants
}

func GetDefenceConstants() map[string]DefenceConstants {
	var mineConstants map[string]DefenceConstants
	constantsFile, _ := os.Open("resources/DefenceConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &mineConstants)
	if err != nil {
		log.Fatal("Error in initializing defence constants: ", err)
		return nil
	}
	return mineConstants
}

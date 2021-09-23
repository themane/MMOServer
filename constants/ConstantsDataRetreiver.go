package constants

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetMineConstants() map[string]ResourceConstants {
	var mineConstants map[string]ResourceConstants
	constantsFile, _ := os.Open("resources/MineConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &mineConstants)
	if err != nil {
		log.Fatal("Error in initializing mine constants: ", err)
		return nil
	}
	return mineConstants
}

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

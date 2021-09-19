package constants

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetWaterConstants() ResourceConstants {
	var waterConstants ResourceConstants
	constantsFile, _ := os.Open("resources/WaterConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &waterConstants)
	if err != nil {
		log.Print(err)
		return ResourceConstants{}
	}
	return waterConstants
}

func GetGrapheneConstants() ResourceConstants {
	var grapheneConstants ResourceConstants
	constantsFile, _ := os.Open("resources/GrapheneConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &grapheneConstants)
	if err != nil {
		log.Print(err)
		return ResourceConstants{}
	}
	return grapheneConstants
}

func GetExperienceConstants() ExperienceConstants {
	var experienceConstants ExperienceConstants
	constantsFile, _ := os.Open("resources/ExperienceConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &experienceConstants)
	if err != nil {
		log.Print(err)
		return ExperienceConstants{}
	}
	return experienceConstants
}

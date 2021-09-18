package constants

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetWaterConstants() models.ResourceConstants {
	var waterConstants models.ResourceConstants
	constantsFile, _ := os.Open("resources/WaterConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &waterConstants)
	if err != nil {
		log.Print(err)
		return models.ResourceConstants{}
	}
	return waterConstants
}

func GetGrapheneConstants() models.ResourceConstants {
	var grapheneConstants models.ResourceConstants
	constantsFile, _ := os.Open("resources/GrapheneConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &grapheneConstants)
	if err != nil {
		log.Print(err)
		return models.ResourceConstants{}
	}
	return grapheneConstants
}

func GetExperienceConstants() models.ExperienceConstants {
	var experienceConstants models.ExperienceConstants
	constantsFile, _ := os.Open("resources/ExperienceConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	err := json.Unmarshal(responseByteValue, &experienceConstants)
	if err != nil {
		log.Print(err)
		return models.ExperienceConstants{}
	}
	return experienceConstants
}

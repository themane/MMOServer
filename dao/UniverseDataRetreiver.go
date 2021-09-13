package dao

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetUniverse() models.Universe {
	var universe models.Universe
	universeFile, _ := os.Open("resources/Universe.json")
	responseByteValue, _ := ioutil.ReadAll(universeFile)
	err := json.Unmarshal(responseByteValue, &universe)
	if err != nil {
		log.Print(err)
		return models.Universe{}
	}
	return universe
}

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

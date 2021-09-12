package dao

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"os"
)

func GetUniverse() models.Universe {
	var universe models.Universe
	universeFile, _ := os.Open("resources/Universe.json")
	responseByteValue, _ := ioutil.ReadAll(universeFile)
	json.Unmarshal(responseByteValue, &universe)
	return universe
}

func GetWaterConstants() models.ResourceConstants {
	var waterConstants models.ResourceConstants
	constantsFile, _ := os.Open("resources/WaterConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	json.Unmarshal(responseByteValue, &waterConstants)
	return waterConstants
}

func GetGrapheneConstants() models.ResourceConstants {
	var grapheneConstants models.ResourceConstants
	constantsFile, _ := os.Open("resources/GrapheneConstants.json")
	responseByteValue, _ := ioutil.ReadAll(constantsFile)
	json.Unmarshal(responseByteValue, &grapheneConstants)
	return grapheneConstants
}

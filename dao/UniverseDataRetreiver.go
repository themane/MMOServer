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

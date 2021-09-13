package dao

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetClanData(Name string) models.ClanData {
	var clanData models.ClanData
	switch Name {
	case "MindKrackers":
		jsonFile, _ := os.Open("sample_user_data/mind_krackers.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &clanData)
		if err != nil {
			log.Print(err)
			return models.ClanData{}
		}
	}
	return clanData
}

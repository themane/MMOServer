package mongoRepository

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetNamedClan(Name string) models.ClanData {
	var clanData models.ClanData
	switch Name {
	case "Mind Krackers":
		jsonFile, _ := os.Open("sample_data/clans/mind_krackers.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &clanData)
		if err != nil {
			log.Print(err)
			return models.ClanData{}
		}
	}
	return clanData
}

func GetClan(Id string) models.ClanData {
	var clanData models.ClanData
	switch Id {
	case "6f744cf6-c873-4c80-8fa6-62da47d87090":
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

package dao

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetUserData(Username string) models.UserData {
	var userData models.UserData
	switch Username {
	case "devashish":
		jsonFile, _ := os.Open("sample_data/users/devashish.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			return models.UserData{}
		}
	case "nehal":
		jsonFile, _ := os.Open("sample_data/users/nehal.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			return models.UserData{}
		}
	case "parth":
		jsonFile, _ := os.Open("sample_data/users/parth.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			return models.UserData{}
		}
	case "sneha":
		jsonFile, _ := os.Open("sample_data/users/sneha.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			return models.UserData{}
		}
	case "sweta":
		jsonFile, _ := os.Open("sample_data/users/sweta.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(responseByteValue, &userData)
		if err != nil {
			log.Print(err)
			return models.UserData{}
		}
	}
	return userData
}

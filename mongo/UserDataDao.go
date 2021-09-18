package mongo

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

func GetUserData(username string) models.UserData {
	var userData models.UserData
	switch username {
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

func UpdateUserData(username string, userData models.UserData) {
	switch username {
	case "devashish":
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile("sample_data/users/devashish.json", updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	case "nehal":
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile("sample_data/users/nehal.json", updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	case "parth":
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile("sample_data/users/parth.json", updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	case "sneha":
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile("sample_data/users/sneha.json", updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	case "sweta":
		updatedResponse, _ := json.Marshal(userData)
		err2 := ioutil.WriteFile("sample_data/users/sweta.json", updatedResponse, 0644)
		if err2 != nil {
			log.Print(err2)
		}
	}
}

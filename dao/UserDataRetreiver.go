package dao

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"os"
)

func GetUserData(Username string) models.UserData {
	var userData models.UserData
	switch Username {
	case "devashish":
		jsonFile, _ := os.Open("sample_user_data/devashish.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &userData)
	case "nehal":
		jsonFile, _ := os.Open("sample_user_data/nehal.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &userData)
	case "parth":
		jsonFile, _ := os.Open("sample_user_data/parth.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &userData)
	case "sneha":
		jsonFile, _ := os.Open("sample_user_data/sneha.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &userData)
	case "sweta":
		jsonFile, _ := os.Open("sample_user_data/sweta.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &userData)
	}
	return userData
}

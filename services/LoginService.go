package services

import (
	"encoding/json"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"os"
)

func Login(Username string) models.LoginResponse {
	var response models.LoginResponse
	switch Username {
	case "devashish":
		jsonFile, _ := os.Open("sample_responses/PlanetConfigResponse1.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &response)
	case "nehal":
		jsonFile, _ := os.Open("sample_responses/PlanetConfigResponse2.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &response)
	case "parth":
		jsonFile, _ := os.Open("sample_responses/PlanetConfigResponse3.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &response)
	case "sneha":
		jsonFile, _ := os.Open("sample_responses/PlanetConfigResponse4.json")
		responseByteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(responseByteValue, &response)
	}
	return response
}

package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/models"
	"io/ioutil"
	"log"
	"os"
)

// Login godoc
// @Summary Login API
// @Description Login verification and first load of complete user data
// @Tags Login
// @Accept json
// @Produce json
// @Param username query string true "valid username for login"
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func Login(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.LoginRequest
	json.Unmarshal(body, &request)
	log.Printf("Logged in user: %s", request.Username)

	var response models.LoginResponse
	switch request.Username {
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
	c.JSON(200, response)
	//c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}

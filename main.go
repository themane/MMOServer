package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	r.GET("/ping", ping)
	r.POST("/login", login)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test 3 Pong",
	})
}

func login(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request LoginRequest
	json.Unmarshal(body, &request)
	log.Printf("Logged in user: %s", request.Username)

	var response LoginResponse
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

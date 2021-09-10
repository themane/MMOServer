package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/themane/MMOServer/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title MMO Game Server
// @version 1.0.0
// @description This is the server for new MMO Game based in space.
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://mmo-server-4xcaklgmnq-el.a.run.app
// @BasePath /
// @schemes https
func main() {
	r := gin.Default()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	url := ginSwagger.URL(baseURL + "/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/ping", ping)
	r.POST("/login", login)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

// HealthCheck ping godoc
// @Summary Pings the server
// @Description get the version and check the health of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test 3 Pong",
	})
}

// Login godoc
// @Summary Login API
// @Description login verification and first load of complete user data
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
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

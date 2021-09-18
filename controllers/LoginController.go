package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

type LoginController struct {
	userRepository     models.UserRepository
	universeRepository models.UniverseRepository
}

func NewLoginController(userRepository models.UserRepository, universeRepository models.UniverseRepository) *LoginController {
	return &LoginController{
		userRepository:     userRepository,
		universeRepository: universeRepository,
	}
}

// Login godoc
// @Summary Login API
// @Description Login verification and first load of complete user data
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func (l *LoginController) Login(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.LoginRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Logged in user: %s", request.Username)

	response := services.Login(request.Username)
	c.JSON(200, response)
}

// RefreshPopulation godoc
// @Summary Refresh population API
// @Description Refresh endpoint to quickly refresh population data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Population
// @Router /refresh/population [post]
func (l *LoginController) RefreshPopulation(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Refreshing population data for: %s", request.Username)

	response := services.RefreshPopulation(request.Username, request.PlanetId)
	c.JSON(200, response)
}

// RefreshResources godoc
// @Summary Refresh resources API
// @Description Refresh endpoint to quickly refresh resources data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Resources
// @Router /refresh/resources [post]
func (l *LoginController) RefreshResources(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Refreshing resources data for: %s", request.Username)

	response := services.RefreshResources(request.Username, request.PlanetId)
	c.JSON(200, response)
}

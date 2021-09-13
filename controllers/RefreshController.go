package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

// RefreshPopulationController godoc
// @Summary Refresh population API
// @Description Refresh endpoint to quickly refresh population data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Population
// @Router /refresh/population [post]
func RefreshPopulationController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.RefreshRequest
	json.Unmarshal(body, &request)
	log.Printf("Refreshing population data for: %s", request.Username)

	response := services.RefreshPopulation(request.Username, request.PlanetId)
	c.JSON(200, response)
}

// RefreshResourcesController godoc
// @Summary Refresh resources API
// @Description Refresh endpoint to quickly refresh resources data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Resources
// @Router /refresh/resources [post]
func RefreshResourcesController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.RefreshRequest
	json.Unmarshal(body, &request)
	log.Printf("Refreshing resources data for: %s", request.Username)

	response := services.RefreshResources(request.Username, request.PlanetId)
	c.JSON(200, response)
}

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
// @Tags Refresh
// @Accept json
// @Produce json
// @Param username query string true "valid username for data retrieval"
// @Success 200 {object} models.Population
// @Router /refresh/population [post]
func RefreshPopulationController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.LoginRequest
	json.Unmarshal(body, &request)
	log.Printf("Refreshing population data for: %s", request.Username)

	response := services.RefreshPopulation(request.Username)
	c.JSON(200, response)
}

// RefreshResourcesController godoc
// @Summary Refresh resources API
// @Description Refresh endpoint to quickly refresh resources data with the latest values
// @Tags Refresh
// @Accept json
// @Produce json
// @Param username query string true "valid username for data retrieval"
// @Success 200 {object} models.Resources
// @Router /refresh/resources [post]
func RefreshResourcesController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.LoginRequest
	json.Unmarshal(body, &request)
	log.Printf("Refreshing resources data for: %s", request.Username)

	response := services.RefreshResources(request.Username)
	c.JSON(200, response)
}

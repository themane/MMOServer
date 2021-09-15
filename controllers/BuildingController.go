package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

// UpgradeBuildingController godoc
// @Summary Upgrade Building API
// @Description Upgrade API for new building or upgrading built one
// @Tags Upgrade
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param building_id query string true "building identifier"
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func UpgradeBuildingController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.UpgradeBuildingRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Upgrading: %s, %s, %s", request.Username, request.PlanetId, request.BuildingId)

	response, err2 := services.UpgradeBuilding(request.Username, request.PlanetId, request.BuildingId)
	if len(err2) > 0 {
		c.JSON(500, models.UpdateResponse{Message: response, Error: err2})
		return
	}
	c.JSON(200, models.UpdateResponse{Message: response, Error: err2})
}

package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

type BuildingController struct {
	userRepository repoModels.UserRepository
}

func NewBuildingController(userRepository *repoModels.UserRepository) *BuildingController {
	return &BuildingController{
		userRepository: *userRepository,
	}
}

// UpgradeBuilding godoc
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
func (b *BuildingController) UpgradeBuilding(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.UpgradeBuildingRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Upgrading: %s, %s, %s", request.Username, request.PlanetId, request.BuildingId)

	err = services.UpgradeBuilding(request.Username, request.PlanetId, request.BuildingId, b.userRepository)
	if err != nil {
		c.JSON(500, controllerModels.UpdateResponse{Error: err.Error()})
		return
	}
	c.JSON(200, controllerModels.UpdateResponse{Message: "Successfully upgraded"})
}

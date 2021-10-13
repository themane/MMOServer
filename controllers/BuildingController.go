package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
)

type BuildingController struct {
	buildingService *services.BuildingService
	refreshService  *services.QuickRefreshService
	logger          *constants.LoggingUtils
}

func NewBuildingController(
	userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(userRepository, buildingConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		logger: constants.NewLoggingUtils("BUILDING_CONTROLLER", logLevel),
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
		b.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	b.logger.Printf("Upgrading: %s, %s, %s", request.Username, request.PlanetId, request.BuildingId)

	err = b.buildingService.UpgradeBuilding(request.Username, request.PlanetId, request.BuildingId)
	if err != nil {
		c.JSON(500, controllerModels.UpdateResponse{Error: err.Error()})
		return
	}
	response, err := b.refreshService.RefreshPlanet(request.Username, request.PlanetId)
	if err != nil {
		b.logger.Error("error in gathering planet data for: "+request.PlanetId, err)
		c.JSON(500, "internal server error. contact administrators for more info")
		return
	}
	c.JSON(200, response)
}

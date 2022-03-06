package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
	"strconv"
)

type BuildingController struct {
	buildingService *services.BuildingService
	planetService   *services.PlanetService
	refreshService  *services.QuickRefreshService
	logger          *constants.LoggingUtils
}

func NewBuildingController(
	userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(userRepository, upgradeConstants, logLevel),
		planetService:   services.NewPlanetService(userRepository, buildingConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
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
// @Success 200 {object} models.PlanetResponse
// @Router /upgrade/building [put]
func (b *BuildingController) UpgradeBuilding(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "building_id")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	b.logger.Printf("Upgrading: %s, %s, %s", parsedParams["username"], parsedParams["planet_id"], parsedParams["building_id"])

	err = b.buildingService.UpgradeBuilding(parsedParams["username"], parsedParams["planet_id"], parsedParams["building_id"])
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	response, err := b.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		b.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// UpdateWorkers godoc
// @Summary Update workers for a building
// @Description Update API for updating workers employed on a building
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param building_id query string true "building identifier"
// @Param workers query int true "updated workers count"
// @Success 200 {object} models.PlanetResponse
// @Router /update/workers [put]
func (b *BuildingController) UpdateWorkers(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "building_id", "workers")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	workers, err := strconv.Atoi(parsedParams["workers"])
	if err != nil || workers < 0 {
		c.JSON(400, "invalid worker count")
	}
	b.logger.Printf("Updating workers: %s, %s, %s", parsedParams["username"], parsedParams["planet_id"], parsedParams["building_id"], workers)

	err = b.buildingService.UpdateWorkers(parsedParams["username"], parsedParams["planet_id"], parsedParams["building_id"], workers)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	response, err := b.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		b.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// UpdatePopulationRate godoc
// @Summary Update population generation rate for the planet
// @Description Update API for population generation rate for the planet
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param generation_rate query int true "updated population generation rate"
// @Success 200 {object} models.PlanetResponse
// @Router /update/population-rate [put]
func (b *BuildingController) UpdatePopulationRate(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "generation_rate")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	generationRate, err := strconv.Atoi(parsedParams["generation_rate"])
	if err != nil || generationRate < 0 {
		c.JSON(400, "invalid population generation rate")
	}
	b.logger.Printf("Updating population generation rate: %s, %s, %s", parsedParams["username"], parsedParams["planet_id"], generationRate)

	err = b.planetService.UpdatePopulationRate(parsedParams["username"], parsedParams["planet_id"], generationRate)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	response, err := b.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		b.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

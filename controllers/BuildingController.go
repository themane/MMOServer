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
	unitService     *services.UnitService
	refreshService  *services.QuickRefreshService
	logger          *constants.LoggingUtils
}

func NewBuildingController(
	userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *BuildingController {
	return &BuildingController{
		buildingService: services.NewBuildingService(userRepository, upgradeConstants, logLevel),
		planetService:   services.NewPlanetService(userRepository, buildingConstants, logLevel),
		unitService:     services.NewUnitService(userRepository, missionRepository, militaryConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, militaryConstants, speciesConstants, logLevel),
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
		b.logger.Error("Error in updating", err)
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
		b.logger.Error("Error in updating", err)
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
		b.logger.Error("Error in updating", err)
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

// EmployPopulation godoc
// @Summary Employ population as worker or soldier
// @Description Update API for employ unemployed population of the planet
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param workers query int true "workers to be employed"
// @Param soldiers query int true "soldiers to be employed"
// @Success 200 {object} models.PlanetResponse
// @Router /population/recruit [put]
func (b *BuildingController) EmployPopulation(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "workers", "soldiers")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	workers, err := strconv.Atoi(parsedParams["workers"])
	if err != nil || workers < 0 {
		c.JSON(400, "invalid workers count")
	}
	soldiers, err := strconv.Atoi(parsedParams["soldiers"])
	if err != nil || soldiers < 0 {
		c.JSON(400, "invalid soldiers count")
	}
	b.logger.Printf("Employing population: %s, %s, workers: %s, soldiers: %s", parsedParams["username"], parsedParams["planet_id"], workers, soldiers)

	err = b.planetService.EmployPopulation(parsedParams["username"], parsedParams["planet_id"], workers, soldiers)
	if err != nil {
		b.logger.Error("Error in updating", err)
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

// KillPopulation godoc
// @Summary Kill population to reduce
// @Description Update API for kill population to reduce water usage and avoid famine
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param workers query int true "unemployed population to be killed"
// @Param workers query int true "workers to be killed"
// @Param soldiers query int true "soldiers to be killed"
// @Success 200 {object} models.PlanetResponse
// @Router /population/kill [put]
func (b *BuildingController) KillPopulation(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unemployed", "workers", "soldiers")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	unemployed, err := strconv.Atoi(parsedParams["unemployed"])
	if err != nil || unemployed < 0 {
		c.JSON(400, "invalid unemployed population count")
	}
	workers, err := strconv.Atoi(parsedParams["workers"])
	if err != nil || workers < 0 {
		c.JSON(400, "invalid workers count")
	}
	soldiers, err := strconv.Atoi(parsedParams["soldiers"])
	if err != nil || soldiers < 0 {
		c.JSON(400, "invalid soldiers count")
	}
	b.logger.Printf("Killing population: %s, %s, unemployed: %s, workers: %s, soldiers: %s",
		parsedParams["username"], parsedParams["planet_id"], unemployed, workers, soldiers)

	err = b.planetService.KillPopulation(parsedParams["username"], parsedParams["planet_id"], unemployed, workers, soldiers)
	if err != nil {
		b.logger.Error("Error in updating", err)
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

// ReserveResources godoc
// @Summary Reserve resources
// @Description Update API for start reserving of resources to avoid loot
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param water query int true "water to be reserved"
// @Param graphene query int true "graphene to be reserved"
// @Success 200 {object} models.PlanetResponse
// @Router /resource/reserve [put]
func (b *BuildingController) ReserveResources(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "water", "graphene")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	water, err := strconv.Atoi(parsedParams["water"])
	if err != nil || water < 0 {
		c.JSON(400, "invalid water amount")
	}
	graphene, err := strconv.Atoi(parsedParams["graphene"])
	if err != nil || graphene < 0 {
		c.JSON(400, "invalid graphene amount")
	}
	b.logger.Printf("Reserving resources: %s, %s, water: %s, graphene: %s",
		parsedParams["username"], parsedParams["planet_id"], water, graphene)

	err = b.planetService.ReserveResources(parsedParams["username"], parsedParams["planet_id"], water, graphene)
	if err != nil {
		b.logger.Error("Error in updating", err)
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

// CancelReserveResources godoc
// @Summary Cancel reserving of resources
// @Description Update API for canceling ongoing reserving of resources
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.PlanetResponse
// @Router /resource/reserve/cancel [put]
func (b *BuildingController) CancelReserveResources(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	b.logger.Printf("Canceling ongoing reserving of resources: %s, %s",
		parsedParams["username"], parsedParams["planet_id"])

	err = b.planetService.CancelReserveResources(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		b.logger.Error("Error in updating", err)
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

// ExtractReservedResources godoc
// @Summary Extract reserved resources
// @Description Update API for extracting reserved resources for use
// @Tags Update
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param water query int true "water to be extracted"
// @Param graphene query int true "graphene to be extracted"
// @Success 200 {object} models.PlanetResponse
// @Router /resource/reserve/extract [put]
func (b *BuildingController) ExtractReservedResources(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "water", "graphene")
	if err != nil {
		b.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	water, err := strconv.Atoi(parsedParams["water"])
	if err != nil || water < 0 {
		c.JSON(400, "invalid water amount")
	}
	graphene, err := strconv.Atoi(parsedParams["graphene"])
	if err != nil || graphene < 0 {
		c.JSON(400, "invalid graphene amount")
	}
	b.logger.Printf("Extracting reserved resources: %s, %s, water: %s, graphene: %s",
		parsedParams["username"], parsedParams["planet_id"], water, graphene)

	err = b.planetService.ExtractReservedResources(parsedParams["username"], parsedParams["planet_id"], water, graphene)
	if err != nil {
		b.logger.Error("Error in updating", err)
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

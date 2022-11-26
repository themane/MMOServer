package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/controllers/utils"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
)

type RefreshController struct {
	refreshService *services.QuickRefreshService
	sectorService  *services.SectorService
	apiSecret      string
	logger         *constants.LoggingUtils
}

func NewRefreshController(userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	researchConstants map[string]constants.ResearchConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	apiSecret string,
	logLevel string,
) *RefreshController {
	return &RefreshController{
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, logLevel),
		sectorService: services.NewSectorService(userRepository, universeRepository, missionRepository,
			experienceConstants, upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, logLevel),
		apiSecret: apiSecret,
		logger:    constants.NewLoggingUtils("REFRESH_CONTROLLER", logLevel),
	}
}

// RefreshPlanet godoc
// @Summary Refresh planet API
// @Description Refresh endpoint to quickly refresh complete planet data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.PlanetResponse
// @Router /refresh/planet [get]
func (r *RefreshController) RefreshPlanet(c *gin.Context) {
	username, err := utils.ExtractUsername(c, r.apiSecret)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "planet_id")
	if err != nil {
		r.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	r.logger.Printf("Refreshing planet data for: %s", username)
	response, err := r.refreshService.RefreshPlanet(username, parsedParams["planet_id"])
	if err != nil {
		r.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		r.logger.Printf("data not found for user: %s, planet_id: %s", username, parsedParams["planet_id"])
		c.JSON(204, nil)
		return
	}
	c.JSON(200, response)
}

// RefreshUserPlanet godoc
// @Summary Refresh user planet API
// @Description Refresh endpoint to quickly refresh user related planet data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.UserPlanetResponse
// @Router /refresh/user_planet [get]
func (r *RefreshController) RefreshUserPlanet(c *gin.Context) {
	username, err := utils.ExtractUsername(c, r.apiSecret)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "planet_id")
	if err != nil {
		r.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	r.logger.Printf("Refreshing population data for: %s", username)
	response, err := r.refreshService.RefreshUserPlanet(username, parsedParams["planet_id"])
	if err != nil {
		r.logger.Error("error in gathering population data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		r.logger.Printf("population data not found for user: %s, planet_id: %s", username, parsedParams["planet_id"])
		c.JSON(204, nil)
		return
	}
	c.JSON(200, response)
}

// Visit godoc
// @Summary Visit Sector API
// @Description Endpoint to switch to another sector to visit and check globally available data on planets
// @Tags Sector
// @Accept json
// @Produce json
// @Param sector_id query string true "sector id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/visit [get]
func (r *RefreshController) Visit(c *gin.Context) {
	username, err := utils.ExtractUsername(c, r.apiSecret)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "sector_id")
	if err != nil {
		r.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	r.logger.Printf("Username: %s visiting sector: %s", username, parsedParams["sector_id"])
	response, err := r.sectorService.Visit(username, parsedParams["sector_id"])
	if err != nil {
		r.logger.Error("error in visiting sector: "+parsedParams["sector_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "internal server error. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// Teleport godoc
// @Summary Teleport to Sector API
// @Description Endpoint to teleport to owned planet in another sector
// @Tags Sector
// @Accept json
// @Produce json
// @Param planet_id query string true "planet id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/teleport [get]
func (r *RefreshController) Teleport(c *gin.Context) {
	username, err := utils.ExtractUsername(c, r.apiSecret)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "planet_id")
	if err != nil {
		r.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	r.logger.Printf("Username: %s teleporting to planet: %s", username, parsedParams["planet_id"])
	response, err := r.sectorService.Teleport(username, parsedParams["planet_id"])
	if err != nil {
		r.logger.Error("error in visiting planet: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "internal server error. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

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

type LoginController struct {
	loginService   *services.LoginService
	refreshService *services.QuickRefreshService
	sectorService  *services.SectorService
	logger         *constants.LoggingUtils
}

func NewLoginController(userRepository models.UserRepository,
	clanRepository models.ClanRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *LoginController {
	return &LoginController{
		loginService: services.NewLoginService(userRepository, clanRepository, universeRepository, missionRepository,
			experienceConstants, upgradeConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, speciesConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, speciesConstants, logLevel),
		sectorService: services.NewSectorService(userRepository, universeRepository, missionRepository,
			experienceConstants, upgradeConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, speciesConstants, logLevel),
		logger: constants.NewLoggingUtils("LOGIN_CONTROLLER", logLevel),
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
	var request controllerModels.LoginRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Logged in user: %s", request.Username)

	response, err := l.loginService.Login(request.Username)
	if err != nil {
		l.logger.Error("error in getting user data", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := "User data not found"
		l.logger.Info(msg)
		c.JSON(204, msg)
		return
	}
	c.JSON(200, response)
}

// RefreshPlanet godoc
// @Summary Refresh planet API
// @Description Refresh endpoint to quickly refresh complete planet data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.PlanetResponse
// @Router /refresh/planet [get]
func (l *LoginController) RefreshPlanet(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id")
	if err != nil {
		l.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Refreshing planet data for: %s", parsedParams["username"])
	response, err := l.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		l.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		l.logger.Printf("data not found for user: %s, planet_id: %s", parsedParams["username"], parsedParams["planet_id"])
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
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.UserPlanetResponse
// @Router /refresh/user_planet [get]
func (l *LoginController) RefreshUserPlanet(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id")
	if err != nil {
		l.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Refreshing population data for: %s", parsedParams["username"])
	response, err := l.refreshService.RefreshUserPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		l.logger.Error("error in gathering population data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		l.logger.Printf("population data not found for user: %s, planet_id: %s", parsedParams["username"], parsedParams["planet_id"])
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
// @Param username query string true "user identifier"
// @Param sector_id query string true "sector id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/visit [get]
func (l *LoginController) Visit(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "sector_id")
	if err != nil {
		l.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Username: %s visiting sector: %s", parsedParams["username"], parsedParams["sector_id"])
	response, err := l.sectorService.Visit(parsedParams["username"], parsedParams["sector_id"])
	if err != nil {
		l.logger.Error("error in visiting sector: "+parsedParams["sector_id"], err)
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
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/teleport [get]
func (l *LoginController) Teleport(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id")
	if err != nil {
		l.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Username: %s teleporting to planet: %s", parsedParams["username"], parsedParams["planet_id"])
	response, err := l.sectorService.Teleport(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		l.logger.Error("error in visiting planet: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "internal server error. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

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
	logger         *constants.LoggingUtils
}

func NewLoginController(userRepository models.UserRepository,
	clanRepository models.ClanRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *LoginController {
	return &LoginController{
		loginService: services.NewLoginService(userRepository, clanRepository, universeRepository, missionRepository,
			experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
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
		c.JSON(500, "error in getting user data. contact administrators for more info")
		return
	}
	if response == nil {
		msg := "User data not found"
		l.logger.Info(msg)
		c.JSON(204, msg)
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
// @Success 200 {object} models.OccupiedPlanet
// @Router /refresh/planet [get]
func (l *LoginController) RefreshPlanet(c *gin.Context) {
	values := c.Request.URL.Query()
	username := values["username"][0]
	planetId := values["planet_id"][1]
	l.logger.Printf("Refreshing planet data for: %s", username)
	response, err := l.refreshService.RefreshPlanet(username, planetId)
	if err != nil {
		l.logger.Error("error in gathering planet data for: "+planetId, err)
		c.JSON(500, "internal server error. contact administrators for more info")
		return
	}
	if response == nil {
		l.logger.Printf("data not found for user: %s, planet_id: %s", username, planetId)
		c.JSON(204, "data not found")
	}
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
// @Router /refresh/population [get]
func (l *LoginController) RefreshPopulation(c *gin.Context) {
	values := c.Request.URL.Query()
	username := values["username"][0]
	planetId := values["planet_id"][1]
	l.logger.Printf("Refreshing population data for: %s", username)

	response, err := l.refreshService.RefreshPopulation(username, planetId)
	if err != nil {
		l.logger.Error("error in gathering population data for: "+planetId, err)
		c.JSON(500, "error in getting user data. contact administrators for more info")
		return
	}
	if response == nil {
		l.logger.Printf("population data not found for user: %s, planet_id: %s", username, planetId)
		c.JSON(204, "data not found")
	}
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
// @Router /refresh/resources [get]
func (l *LoginController) RefreshResources(c *gin.Context) {
	values := c.Request.URL.Query()
	username := values["username"][0]
	planetId := values["planet_id"][1]
	l.logger.Printf("Refreshing resources data for: %s", username)

	response, err := l.refreshService.RefreshResources(username, planetId)
	if err != nil {
		l.logger.Error("error in getting resources data for: "+planetId, err)
		c.JSON(500, "error in getting user data. contact administrators for more info")
		return
	}
	if response == nil {
		l.logger.Printf("resources data not found for user: %s, planet_id: %s", username, planetId)
		c.JSON(204, "data not found")
	}
	c.JSON(200, response)
}

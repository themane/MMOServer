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
			buildingConstants, mineConstants, defenceConstants, logLevel),
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

// RefreshPopulation godoc
// @Summary Refresh population API
// @Description Refresh endpoint to quickly refresh population data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Population
// @Router /refresh/population [post]
func (l *LoginController) RefreshPopulation(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Refreshing population data for: %s", request.Username)

	response, err := l.refreshService.RefreshPopulation(request.Username, request.PlanetId)
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

// RefreshResources godoc
// @Summary Refresh resources API
// @Description Refresh endpoint to quickly refresh resources data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} models.Resources
// @Router /refresh/resources [post]
func (l *LoginController) RefreshResources(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Refreshing resources data for: %s", request.Username)

	response, err := l.refreshService.RefreshResources(request.Username, request.PlanetId)
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

// RefreshMine godoc
// @Summary Refresh mine API
// @Description Refresh endpoint to quickly refresh mine data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param mine_id query string true "mine identifier"
// @Success 200 {object} models.Mine
// @Router /refresh/mines [post]
func (l *LoginController) RefreshMine(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RefreshMineRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Refreshing mines data for: %s", request.Username)

	response, err := l.refreshService.RefreshMine(request.Username, request.PlanetId, request.MineId)
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

// RefreshShields godoc
// @Summary Refresh shields API
// @Description Refresh endpoint to quickly refresh shields data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} []models.Shield
// @Router /refresh/shields [post]
func (l *LoginController) RefreshShields(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Refreshing shields data for: %s", request.Username)

	response, err := l.refreshService.RefreshShields(request.Username, request.PlanetId)
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

// RefreshMissions godoc
// @Summary Refresh missions API
// @Description Refresh endpoint to quickly refresh missions data with the latest values
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Success 200 {object} map[string][]models.ActiveMission
// @Router /refresh/shields [post]
func (l *LoginController) RefreshMissions(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RefreshRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		l.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	l.logger.Printf("Refreshing missions data for: %s", request.Username)

	response, err := l.refreshService.RefreshMissions(request.Username, request.PlanetId)
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

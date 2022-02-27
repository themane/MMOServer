package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"net/url"
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
		sectorService: services.NewSectorService(userRepository, universeRepository, missionRepository,
			experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
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
// @Success 200 {object} models.PlanetResponse
// @Router /refresh/planet [get]
func (l *LoginController) RefreshPlanet(c *gin.Context) {
	values := c.Request.URL.Query()
	username, planetId, err := l.getPlanetParams(values)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Refreshing planet data for: %s", *username)
	response, err := l.refreshService.RefreshPlanet(*username, *planetId)
	if err != nil {
		l.logger.Error("error in gathering planet data for: "+*planetId, err)
		c.JSON(500, "internal server error. contact administrators for more info")
		return
	}
	if response == nil {
		l.logger.Printf("data not found for user: %s, planet_id: %s", *username, *planetId)
		c.JSON(204, "data not found")
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
	username, planetId, err := l.getPlanetParams(values)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Refreshing population data for: %s", *username)
	response, err := l.refreshService.RefreshUserPlanet(*username, *planetId)
	if err != nil {
		l.logger.Error("error in gathering population data for: "+*planetId, err)
		c.JSON(500, "error in getting user data. contact administrators for more info")
		return
	}
	if response == nil {
		l.logger.Printf("population data not found for user: %s, planet_id: %s", *username, *planetId)
		c.JSON(204, "data not found")
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
	username, sectorId, err := l.getSectorParams(values)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Username: %s visiting sector: %s", username, sectorId)
	response, err := l.sectorService.Visit(*username, *sectorId)
	if err != nil {
		l.logger.Error("error in visiting sector: "+*sectorId, err)
		c.JSON(500, "internal server error. contact administrators for more info")
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
	username, planetId, err := l.getPlanetParams(values)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	l.logger.Printf("Username: %s teleporting to planet: %s", username, planetId)
	response, err := l.sectorService.Teleport(*username, *planetId)
	if err != nil {
		l.logger.Error("error in visiting planet: "+*planetId, err)
		c.JSON(500, "internal server error. contact administrators for more info")
		return
	}
	c.JSON(200, response)
}

func (l *LoginController) getPlanetParams(values url.Values) (*string, *string, error) {
	if usernames, ok := values["username"]; ok {
		if planetIds, ok := values["planet_id"]; ok {
			if len(usernames) == 1 && len(planetIds) == 1 {
				return &usernames[0], &planetIds[0], nil
			}
		}
	}
	msg := "cannot parse request parameters correctly"
	l.logger.Println(msg)
	return nil, nil, errors.New(msg)
}

func (l *LoginController) getSectorParams(values url.Values) (*string, *string, error) {
	if usernames, ok := values["username"]; ok {
		if sectorIds, ok := values["sector_id"]; ok {
			if len(usernames) == 1 && len(sectorIds) == 1 {
				return &usernames[0], &sectorIds[0], nil
			}
		}
	}
	msg := "cannot parse request parameters correctly"
	l.logger.Println(msg)
	return nil, nil, errors.New(msg)
}

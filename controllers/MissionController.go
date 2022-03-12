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

type MissionController struct {
	missionService *services.MissionService
	refreshService *services.QuickRefreshService
	logger         *constants.LoggingUtils
}

func NewMissionController(userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *MissionController {
	return &MissionController{
		missionService: services.NewMissionService(userRepository, universeRepository, missionRepository, shipConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, defenceConstants, shipConstants, speciesConstants, logLevel),
		logger: constants.NewLoggingUtils("MISSION_CONTROLLER", logLevel),
	}
}

// Spy godoc
// @Summary Spy API
// @Description Endpoint to launch spy mission with available scout ships
// @Tags Attack
// @Accept json
// @Produce json
// @Param attacker query string true "attacker username"
// @Param from_planet_id query string true "spy launch planet identifier"
// @Param to_planet_id query string true "spy destination planet identifier"
// @Param scouts query object true "scout ship details"
// @Success 200 {object} controllerModels.PlanetResponse
// @Router /spy [post]
func (a *MissionController) Spy(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.SpyRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		a.logger.Error("request not parseable", err)
		c.JSON(400, err.Error())
		return
	}
	a.logger.Printf("Launching spy mission from %s to %s", request.FromPlanetId, request.ToPlanetId)

	err = a.missionService.Spy(request)
	if err != nil {
		a.logger.Error("error in launching spy mission", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "internal server error. contact administrators for more info", HttpCode: 500})
		return
	}
	response, err := a.refreshService.RefreshPlanet(request.Username, request.FromPlanetId)
	if err != nil {
		a.logger.Error("error in gathering planet data for: "+request.FromPlanetId, err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// Attack godoc
// @Summary Attack API
// @Description Endpoint to launch attack mission on other planet
// @Tags Attack
// @Accept json
// @Produce json
// @Param attacker query string true "attacker username"
// @Param from_planet_id query string true "spy launch planet identifier"
// @Param to_planet_id query string true "spy destination planet identifier"
// @Param formation query object true "attack ships details"
// @Success 200 {object} controllerModels.PlanetResponse
// @Router /attack [post]
func (a *MissionController) Attack(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.AttackRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		a.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	a.logger.Printf("Launching attack mission from %s to %s", request.FromPlanetId, request.ToPlanetId)

	err = a.missionService.Attack(request)
	if err != nil {
		a.logger.Error("error in launching attack mission", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "internal server error. contact administrators for more info", HttpCode: 500})
		return
	}
	response, err := a.refreshService.RefreshPlanet(request.Username, request.FromPlanetId)
	if err != nil {
		a.logger.Error("error in gathering planet data for: "+request.FromPlanetId, err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

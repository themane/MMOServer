package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/schedulers"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

type AttackController struct {
	attackService *services.AttackService
}

func NewAttackController(userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	scheduledMissionManager schedulers.ScheduledMissionManager,
	shipConstants map[string]constants.ShipConstants,
) *AttackController {
	return &AttackController{
		attackService: services.NewAttackService(userRepository, universeRepository, missionRepository, scheduledMissionManager, shipConstants),
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
// @Success 200 {object} controllerModels.MissionResponse
// @Router /spy [post]
func (a *AttackController) Spy(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.SpyRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Launching spy mission from %s to %s", request.FromPlanetId, request.ToPlanetId)

	response, err := a.attackService.Spy(request)
	if err != nil {
		c.JSON(500, err.Error())
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
// @Success 200 {object} controllerModels.MissionResponse
// @Router /attack [post]
func (a *AttackController) Attack(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.AttackRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Print(err)
		c.JSON(400, "Request not parseable")
		return
	}
	log.Printf("Launching attack mission from %s to %s", request.FromPlanetId, request.ToPlanetId)

	response, err := a.attackService.Attack(request)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, response)
}

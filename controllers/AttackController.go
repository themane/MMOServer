package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/mongoRepository/models"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type AttackController struct {
}

func NewAttackController(userRepository models.UserRepository,
	clanRepository models.ClanRepository,
	universeRepository models.UniverseRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
) *AttackController {
	return &AttackController{}
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
// @Param scouts query []Formation true "scout ship details"
// @Success 200 {object} models.AttackResponse
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

	randomMinutes := rand.Intn(5) + 5
	attackTime := time.Now().Add(time.Minute * time.Duration(randomMinutes))
	returnTime := time.Now().Add(time.Minute * time.Duration(randomMinutes) * 2)
	response := controllerModels.AttackResponse{AttackTime: attackTime.String(), ReturnTime: returnTime.String()}
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
// @Param formation query map[string]map[string][]Formation true "attack ships details"
// @Success 200 {object} models.AttackResponse
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

	randomMinutes := rand.Intn(100) + 13
	attackTime := time.Now().Add(time.Minute * time.Duration(randomMinutes))
	returnTime := time.Now().Add(time.Minute * time.Duration(randomMinutes) * 2)
	response := controllerModels.AttackResponse{AttackTime: attackTime.String(), ReturnTime: returnTime.String()}
	c.JSON(200, response)
}

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

type SectorController struct {
	sectorService *services.SectorService
	logger        *constants.LoggingUtils
}

func NewSectorController(userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	experienceConstants map[string]constants.ExperienceConstants,
	buildingConstants map[string]constants.BuildingConstants,
	mineConstants map[string]constants.MiningConstants,
	defenceConstants map[string]constants.DefenceConstants,
	shipConstants map[string]constants.ShipConstants,
	logLevel string,
) *SectorController {
	return &SectorController{
		sectorService: services.NewSectorService(userRepository, universeRepository, missionRepository, experienceConstants,
			buildingConstants, mineConstants, defenceConstants, shipConstants, logLevel),
		logger: constants.NewLoggingUtils("SECTOR_CONTROLLER", logLevel),
	}
}

// Visit godoc
// @Summary Visit Sector API
// @Description Endpoint to switch to another sector to visit and check globally available data on planets
// @Tags Sector
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param sector query string true "sector id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/visit [get]
func (s *SectorController) Visit(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.VisitSectorRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		s.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	s.logger.Printf("Username: %s visiting sector: %s", request.Username, request.Sector)

	response, err := s.sectorService.Visit(request)
	if err != nil {
		s.logger.Error("error in visiting sector: "+request.Sector, err)
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
// @Param planet query string true "planet id to visit"
// @Success 200 {object} controllerModels.SectorResponse
// @Router /sector/teleport [get]
func (s *SectorController) Teleport(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.TeleportRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		s.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	s.logger.Printf("Username: %s teleporting to planet: %s", request.Username, request.Planet)

	response, err := s.sectorService.Teleport(request)
	if err != nil {
		s.logger.Error("error in visiting planet: "+request.Planet, err)
		c.JSON(500, "internal server error. contact administrators for more info")
		return
	}
	c.JSON(200, response)
}

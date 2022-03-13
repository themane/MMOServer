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
	"strconv"
)

type UnitsController struct {
	unitService            *services.UnitService
	unitsDeploymentService *services.UnitsDeploymentService
	refreshService         *services.QuickRefreshService
	logger                 *constants.LoggingUtils
}

func NewUnitsController(
	userRepository models.UserRepository,
	universeRepository models.UniverseRepository,
	missionRepository models.MissionRepository,
	upgradeConstants map[string]constants.UpgradeConstants,
	buildingConstants map[string]map[string]map[string]interface{},
	mineConstants map[string]constants.MiningConstants,
	militaryConstants map[string]constants.MilitaryConstants,
	speciesConstants map[string]constants.SpeciesConstants,
	logLevel string,
) *UnitsController {
	return &UnitsController{
		unitService:            services.NewUnitService(userRepository, missionRepository, militaryConstants, logLevel),
		unitsDeploymentService: services.NewUnitsDeploymentService(userRepository, missionRepository, militaryConstants, logLevel),
		refreshService: services.NewQuickRefreshService(userRepository, universeRepository, missionRepository,
			upgradeConstants, buildingConstants, mineConstants, militaryConstants, speciesConstants, logLevel),
		logger: constants.NewLoggingUtils("BUILDING_CONTROLLER", logLevel),
	}
}

// ConstructUnits godoc
// @Summary Construct attack or defence unit
// @Description Update API for constructing attack or defence unit
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_name query int true "unit to be constructed"
// @Param quantity query int true "quantity of unit to be constructed"
// @Success 200 {object} models.PlanetResponse
// @Router /unit/construct [put]
func (u *UnitsController) ConstructUnits(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_name", "quantity")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	quantity, err := strconv.Atoi(parsedParams["quantity"])
	if err != nil || quantity < 0 {
		c.JSON(400, "invalid quantity")
	}
	u.logger.Printf("Constructing units: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"], quantity)

	err = u.unitService.ConstructUnits(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"], quantity)
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// CancelUnitsConstruction godoc
// @Summary Cancel construction of attack or defence units
// @Description Update API for canceling construction of attack or defence units
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_name query int true "unit to be canceled"
// @Success 200 {object} models.PlanetResponse
// @Router /unit/cancel [put]
func (u *UnitsController) CancelUnitsConstruction(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_name")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	u.logger.Printf("Canceling construction of units: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"])

	err = u.unitService.CancelUnitsConstruction(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"])
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// DestructUnits godoc
// @Summary Destruct attack or defence unit
// @Description Update API for destructing attack or defence unit for return of population employed and resources
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_name query int true "unit to be destructed"
// @Param quantity query int true "quantity of unit to be destructed"
// @Success 200 {object} models.PlanetResponse
// @Router /unit/destruct [put]
func (u *UnitsController) DestructUnits(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_name", "quantity")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	quantity, err := strconv.Atoi(parsedParams["quantity"])
	if err != nil || quantity < 0 {
		c.JSON(400, "invalid quantity")
	}
	u.logger.Printf("Destructing units: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"], quantity)

	err = u.unitService.DestructUnits(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_name"], quantity)
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// UpgradeDefenceShipCarrier godoc
// @Summary Upgrade defence ship carrier
// @Description Update API for upgrading defence ship carrier
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_id query int true "defence ship carrier id to be upgraded"
// @Success 200 {object} models.PlanetResponse
// @Router /defence_ship_carrier/upgrade [put]
func (u *UnitsController) UpgradeDefenceShipCarrier(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_id")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	u.logger.Printf("Upgrading defence ship carrier: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])

	err = u.unitService.UpgradeDefenceShipCarrier(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// CancelDefenceShipCarrierUpGradation godoc
// @Summary Cancel up-gradation/construction of defence ship carrier
// @Description Update API for cancelling up-gradation/construction of defence ship carrier for return of population employed and resources
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_id query int true "defence ship carrier id to be canceled for up-gradation/construction"
// @Success 200 {object} models.PlanetResponse
// @Router /defence_ship_carrier/cancel [put]
func (u *UnitsController) CancelDefenceShipCarrierUpGradation(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_id")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	u.logger.Printf("Canceling up-gradation/construction of defence ship carrier: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])

	err = u.unitService.CancelDefenceShipCarrierUpGradation(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// DestructDefenceShipCarrier godoc
// @Summary Destruct defence ship carrier
// @Description Update API for destructing defence ship carrier for return of population employed and resources
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_id query int true "defence ship carrier id to be destructed"
// @Success 200 {object} models.PlanetResponse
// @Router /defence_ship_carrier/destruct [put]
func (u *UnitsController) DestructDefenceShipCarrier(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_id")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	u.logger.Printf("Destructing defence ship carrier: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])

	err = u.unitService.DestructDefenceShipCarrier(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"])
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// DeployShipsOnDefenceShipCarrier godoc
// @Summary Deploy ships on defence ship carrier
// @Description Update API for deploying ships as defenders on defence ship carrier
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param unit_id query int true "defence ship carrier id for deployment"
// @Success 200 {object} models.PlanetResponse
// @Router /deploy/defence_ship_carrier/ships [post]
func (u *UnitsController) DeployShipsOnDefenceShipCarrier(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "unit_id")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request map[string]int
	err = json.Unmarshal(body, &request)
	if err != nil {
		u.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	u.logger.Printf("Deploying on defence ship carrier: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"], request)

	err = u.unitsDeploymentService.DeployShipsOnDefenceShipCarrier(parsedParams["username"], parsedParams["planet_id"], parsedParams["unit_id"], request)
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// DeployDefencesOnShield godoc
// @Summary Deploy defences on shield
// @Description Update API for deploying defences as defenders on shield
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param shield_id query int true "shield id for deployment"
// @Success 200 {object} models.PlanetResponse
// @Router /deploy/defence_ship_carrier/ships [post]
func (u *UnitsController) DeployDefencesOnShield(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "shield_id")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request map[string]int
	err = json.Unmarshal(body, &request)
	if err != nil {
		u.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	if _, ok := constants.GetShieldIds()[parsedParams["shield_id"]]; !ok {
		u.logger.Error("Error in parsing params", errors.New("not valid shield"))
		c.JSON(400, "not valid shield")
		return
	}
	u.logger.Printf("Deploying on shield: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["shield_id"], request)

	err = u.unitsDeploymentService.DeployDefencesOnShield(parsedParams["username"], parsedParams["planet_id"], parsedParams["shield_id"], request)
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

// DeployDefenceShipCarrierOnShield godoc
// @Summary Deploy/Remove defence ship carrier on shield
// @Description Update API for deploying or removing defence ship carrier as defenders on shield
// @Tags Unit
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Param planet_id query string true "planet identifier"
// @Param shield_id query string true "shield id for deployment"
// @Param unit_id query string true "defence ship carrier id to be deployed"
// @Param deploy query bool true "to identify deployment vs removal"
// @Success 200 {object} models.PlanetResponse
// @Router /deploy/defence_ship_carrier/ships [put]
func (u *UnitsController) DeployDefenceShipCarrierOnShield(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username", "planet_id", "shield_id", "unit_id", "deploy")
	if err != nil {
		u.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	deploy, err := strconv.ParseBool(parsedParams["deploy"])
	if err != nil {
		c.JSON(400, "invalid deploy value")
	}
	if deploy {
		u.logger.Printf("Deploying on shield: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["shield_id"], parsedParams["unit_id"])
	} else {
		u.logger.Printf("Un-deploying from shield: ", parsedParams["username"], parsedParams["planet_id"], parsedParams["shield_id"], parsedParams["unit_id"])
	}

	err = u.unitsDeploymentService.DeployDefenceShipCarrierOnShield(parsedParams["username"], parsedParams["planet_id"],
		parsedParams["shield_id"], parsedParams["unit_id"], deploy)
	if err != nil {
		u.logger.Error("Error in updating", err)
		c.JSON(500, err.Error())
		return
	}
	response, err := u.refreshService.RefreshPlanet(parsedParams["username"], parsedParams["planet_id"])
	if err != nil {
		u.logger.Error("error in gathering planet data for: "+parsedParams["planet_id"], err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(200, response)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	controllerModels "github.com/themane/MMOServer/controllers/models"
	"github.com/themane/MMOServer/controllers/utils"
	"github.com/themane/MMOServer/mongoRepository"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"time"
)

type RegistrationController struct {
	loginService        *services.LoginService
	registrationService *services.RegistrationService
	apiSecret           string
	logger              *constants.LoggingUtils
}

func NewRegistrationController(userRepository models.UserRepository,
	clanRepository models.ClanRepository,
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
	maxSystems int,
	retries int,
	logLevel string,
) *RegistrationController {
	return &RegistrationController{
		loginService: services.NewLoginService(userRepository, clanRepository, universeRepository, missionRepository,
			experienceConstants, upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, logLevel),
		registrationService: services.NewRegistrationService(userRepository, universeRepository,
			experienceConstants, buildingConstants, mineConstants, militaryConstants, maxSystems, retries, logLevel),
		apiSecret: apiSecret,
		logger:    constants.NewLoggingUtils("REGISTRATION_CONTROLLER", logLevel),
	}
}

// Register godoc
// @Summary Register API
// @Description Registration payload verification and initial assignment of complete user data
// @Tags registration
// @Accept json
// @Produce json
// @Success 200 {object} models.LoginResponse
// @Router /register [post]
func (r *RegistrationController) Register(c *gin.Context) {
	claims, err := utils.ValidateIdToken(c)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	r.logger.Printf("Logged in user: %s", claims["email"])

	body, _ := ioutil.ReadAll(c.Request.Body)
	var request controllerModels.RegistrationRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		r.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	const datePattern = "2006-01-02"
	birthday, err := time.Parse(datePattern, request.Birthday)
	if err != nil {
		r.logger.Error("Error in parsing birthday", err)
		c.JSON(401, err.Error())
		return
	}
	age := today.Sub(birthday).Hours() / 24 / 365
	if age < 13 {
		r.logger.Error("minimum age required to be 13 yrs", err)
		c.JSON(400, "minimum age required to be 13 yrs")
		return
	}
	_, err = time.LoadLocation(request.Location)
	if err != nil {
		r.logger.Error("Error in parsing location", err)
		c.JSON(401, err.Error())
		return
	}

	id := fmt.Sprintf("%v", claims["sub"])
	email := fmt.Sprintf("%v", claims["email"])
	name := fmt.Sprintf("%v", claims["name"])
	err = r.registrationService.Register(request, id, email, name)
	if _, ok := err.(*mongoRepository.NoSuchCombinationError); ok {
		r.logger.Error("error in finding new planet", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in finding new planet. contact administrators for more info", HttpCode: 500})
		return
	}
	if err != nil {
		r.logger.Error("error in getting user registration", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user registration. contact administrators for more info", HttpCode: 500})
		return
	}

	response, err := r.loginService.GoogleLogin(id)
	if err != nil {
		r.logger.Error("error in getting user data", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := "User data not found"
		r.logger.Info(msg)
		c.JSON(204, msg)
		return
	}
	token, err := utils.GenerateToken(response.Profile.Username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Api-Token", token)

	refreshToken, err := utils.GenerateRefreshToken(response.Profile.Username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Refresh-Token", refreshToken)

	c.JSON(200, response)
}

// CheckUsername godoc
// @Summary Check Username API
// @Description Check Username availability
// @Tags registration
// @Accept json
// @Produce json
// @Success 200 {string}
// @Router /check/username [get]
func (r *RegistrationController) CheckUsername(c *gin.Context) {
	values := c.Request.URL.Query()
	parsedParams, err := parseStrings(values, "username")
	if err != nil {
		r.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	r.logger.Printf("Checking username: %s", parsedParams["username"])
	if r.registrationService.UsernameExists(parsedParams["username"]) {
		c.JSON(406, "not available")
		return
	}
	c.JSON(200, "available")
}

// Login godoc
// @Summary Login API
// @Description Login verification and first load of complete user data
// @Tags data retrieval
// @Accept json
// @Produce json
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func (r *RegistrationController) Login(c *gin.Context) {
	claims, err := utils.ValidateIdToken(c)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	r.logger.Printf("Logged in user: %s", claims["email"])
	id := fmt.Sprintf("%v", claims["sub"])
	response, err := r.loginService.GoogleLogin(id)
	if err != nil {
		r.logger.Error("error in getting user data", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := "User data not found"
		r.logger.Info(msg)
		c.JSON(204, msg)
		return
	}
	token, err := utils.GenerateToken(response.Profile.Username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Api-Token", token)

	refreshToken, err := utils.GenerateRefreshToken(response.Profile.Username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Refresh-Token", refreshToken)

	c.JSON(200, response)
}

// RefreshToken godoc
// @Summary Refresh Token API
// @Description Refresh Token
// @Tags data retrieval
// @Accept json
// @Produce json
// @Router /refresh/token [post]
func (r *RegistrationController) RefreshToken(c *gin.Context) {
	username, err := utils.RefreshTokenValid(c, r.apiSecret)
	if err != nil {
		r.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	token, err := utils.GenerateToken(username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Api-Token", token)

	refreshToken, err := utils.GenerateRefreshToken(username, r.apiSecret)
	if err != nil {
		r.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Refresh-Token", refreshToken)
	c.Status(200)
}

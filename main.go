package main

import (
	secretManager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers"
	"github.com/themane/MMOServer/mongoRepository"
	"github.com/themane/MMOServer/mongoRepository/models"
	secretManagerPb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/themane/MMOServer/docs"
)

var once = sync.Once{}
var baseURL string
var mongoDB string
var maxSystems int
var registrationRetries int
var mongoUrlSecretName string
var apiSecretName string
var logLevel string

// @title MMO Game Server
// @version 1.0.0
// @description This is the server for new MMO Game based in space.
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @schemes https
func main() {
	r := gin.Default()

	once.Do(initialize)
	url := ginSwagger.URL(baseURL + "/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	registrationController, refreshController, buildingController, missionController, unitsController := getHandlers()
	//scheduledJobManager.SchedulePlanetUpdates()

	r.GET("/ping", controllers.Ping)

	r.POST("/register", registrationController.Register)
	r.HEAD("/check/username", registrationController.CheckUsername)
	r.POST("/login", registrationController.Login)
	r.GET("/refresh/token", registrationController.RefreshToken)
	r.GET("/refresh/planet", refreshController.RefreshPlanet)
	//r.GET("/refresh/user_planet", loginController.RefreshUserPlanet)

	r.PUT("/upgrade/building", buildingController.UpgradeBuilding)
	r.PUT("/upgrade/building/cancel", buildingController.CancelUpgradeBuilding)
	r.PUT("/update/workers", buildingController.UpdateWorkers)
	r.PUT("/update/soldiers", buildingController.UpdateSoldiers)
	r.PUT("/update/population-growth", buildingController.UpdatePopulationRate)

	r.PUT("/population/recruit", buildingController.EmployPopulation)
	r.PUT("/population/kill", buildingController.KillPopulation)

	r.PUT("/resource/reserve", buildingController.ReserveResources)
	r.PUT("/resource/reserve/cancel", buildingController.CancelReserveResources)
	r.PUT("/resource/reserve/extract", buildingController.ExtractReservedResources)

	r.PUT("/unit/construct", unitsController.ConstructUnits)
	r.PUT("/unit/cancel", unitsController.CancelUnitsConstruction)
	r.PUT("/unit/destruct", unitsController.DestructUnits)
	r.PUT("/defence_ship_carrier/upgrade", unitsController.UpgradeDefenceShipCarrier)
	r.PUT("/defence_ship_carrier/cancel", unitsController.CancelDefenceShipCarrierUpGradation)
	r.PUT("/defence_ship_carrier/destruct", unitsController.DestructDefenceShipCarrier)

	r.POST("/deploy/defence_ship_carrier/ships", unitsController.DeployShipsOnDefenceShipCarrier)
	r.POST("/deploy/shield/defences", unitsController.DeployDefencesOnShield)
	r.PUT("/deploy/shield/defence_ship_carrier", unitsController.DeployDefenceShipCarrierOnShield)

	r.PUT("/research", buildingController.Research)
	r.PUT("/research/cancel", buildingController.CancelResearch)

	r.POST("/spy", missionController.Spy)
	r.POST("/attack", missionController.Attack)

	r.GET("/sector/visit", refreshController.Visit)
	r.GET("/sector/teleport", refreshController.Teleport)

	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func getHandlers() (
	*controllers.RegistrationController,
	*controllers.RefreshController,
	*controllers.BuildingController,
	*controllers.MissionController,
	*controllers.UnitsController,
) {

	log.Println("Initializing handlers")
	mongoURL := accessSecretVersion(mongoUrlSecretName)
	apiSecret := accessSecretVersion(apiSecretName)

	upgradeConstants := constants.GetUpgradeConstants()
	buildingConstants := constants.GetBuildingConstants()
	experienceConstants := constants.GetExperienceConstants()
	mineConstants := constants.GetMiningConstants()
	militaryConstants := constants.GetMilitaryConstants()
	researchConstants := constants.GetResearchConstants()
	speciesConstants := constants.GetSpeciesConstants()

	var userRepository models.UserRepository
	var clanRepository models.ClanRepository
	var universeRepository models.UniverseRepository
	var missionRepository models.MissionRepository
	userRepository = mongoRepository.NewUserRepository(mongoURL, mongoDB, logLevel)
	clanRepository = mongoRepository.NewClanRepository(mongoURL, mongoDB, logLevel)
	universeRepository = mongoRepository.NewUniverseRepository(mongoURL, mongoDB, logLevel)
	missionRepository = mongoRepository.NewMissionRepository(mongoURL, mongoDB, logLevel)
	registrationController := controllers.NewRegistrationController(userRepository, clanRepository, universeRepository, missionRepository,
		experienceConstants, upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants,
		apiSecret, maxSystems, registrationRetries, logLevel)
	refreshController := controllers.NewRefreshController(userRepository, universeRepository, missionRepository,
		experienceConstants, upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, apiSecret, logLevel)
	buildingController := controllers.NewBuildingController(userRepository, universeRepository, missionRepository,
		upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, apiSecret, logLevel)
	missionController := controllers.NewMissionController(userRepository, universeRepository, missionRepository,
		upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, apiSecret, logLevel)
	unitsController := controllers.NewUnitsController(userRepository, universeRepository, missionRepository,
		upgradeConstants, buildingConstants, mineConstants, militaryConstants, researchConstants, speciesConstants, apiSecret, logLevel)

	log.Println("Initialized all handlers")
	return registrationController, refreshController, buildingController, missionController, unitsController
}

func initialize() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	mongoUrlSecretName = os.Getenv("MONGO_SECRET_NAME")
	mongoDB = os.Getenv("MONGO_DB")
	if mongoUrlSecretName == "" || mongoDB == "" {
		log.Fatal("Mongo not configured")
	}

	apiSecretName = os.Getenv("API_SECRET_NAME")
	if apiSecretName == "" {
		apiSecretName = "API_SECRET"
	}

	maxSystemsString := os.Getenv("MAX_SYSTEMS")
	if maxSystemsString == "" {
		maxSystemsString = "10"
	}
	log.Println("USING MAX_SYSTEMS: " + maxSystemsString)
	var err error
	maxSystems, err = strconv.Atoi(maxSystemsString)
	if err != nil {
		log.Fatal(err)
	}

	registrationRetriesString := os.Getenv("REGISTRATION_RETRIES")
	if registrationRetriesString == "" {
		registrationRetriesString = "3"
	}
	log.Println("USING REGISTRATION_RETRIES: " + registrationRetriesString)
	registrationRetries, err = strconv.Atoi(registrationRetriesString)
	if err != nil {
		log.Fatal(err)
	}

	logLevel = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = constants.Info
	}
}

func accessSecretVersion(secretName string) string {
	ctx := context.Background()
	client, err := secretManager.NewClient(ctx)
	if err != nil {
		log.Fatal("Error in initializing client for secret manager: ", err)
		return ""
	}
	defer func(client *secretManager.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal("Error in closing client for secret manager: ", err)
		}
	}(client)
	req := &secretManagerPb.AccessSecretVersionRequest{
		Name: secretName,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatal("Error in calling access API for retrieving secret data: ", err)
		return ""
	}
	return string(result.Payload.Data)
}

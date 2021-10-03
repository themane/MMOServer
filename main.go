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
	"github.com/themane/MMOServer/schedulers"
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
var secretName string

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

	loginController, buildingController, scheduledJobManager := getHandlers()
	scheduledJobManager.SchedulePlanetUpdates()

	r.GET("/ping", controllers.Ping)
	r.POST("/login", loginController.Login)
	r.POST("/refresh/population", loginController.RefreshPopulation)
	r.POST("/refresh/resources", loginController.RefreshResources)
	r.POST("/refresh/mine", loginController.RefreshMine)
	r.POST("/refresh/shields", loginController.RefreshShields)
	r.POST("/upgrade/building", buildingController.UpgradeBuilding)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func getHandlers() (*controllers.LoginController, *controllers.BuildingController, *schedulers.ScheduledJobManager) {
	log.Println("Initializing handlers")
	mongoURL := accessSecretVersion()

	buildingConstants := constants.GetBuildingConstants()
	experienceConstants := constants.GetExperienceConstants()
	mineConstants := constants.GetMiningConstants()
	defenceConstants := constants.GetDefenceConstants()
	shipConstants := constants.GetShipConstants()

	var userRepository models.UserRepository
	var clanRepository models.ClanRepository
	var universeRepository models.UniverseRepository
	userRepository = mongoRepository.NewUserRepository(mongoURL, mongoDB)
	clanRepository = mongoRepository.NewClanRepository(mongoURL, mongoDB)
	universeRepository = mongoRepository.NewUniverseRepository(mongoURL, mongoDB)
	loginController := controllers.NewLoginController(userRepository, clanRepository, universeRepository, experienceConstants, buildingConstants, mineConstants, defenceConstants, shipConstants)
	buildingController := controllers.NewBuildingController(userRepository, buildingConstants)
	scheduledJobManager := schedulers.NewScheduledJobManager(userRepository, universeRepository, mineConstants, maxSystems)
	log.Println("Initialized all handlers")
	return loginController, buildingController, scheduledJobManager
}

func initialize() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "localhost:8080"
	}
	secretName = os.Getenv("SECRET_NAME")
	mongoDB = os.Getenv("MONGO_DB")
	if secretName == "" || mongoDB == "" {
		log.Fatal("Mongo not configured")
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
}

func accessSecretVersion() string {
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

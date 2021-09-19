package main

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers"
	"github.com/themane/MMOServer/mongoRepository"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/schedulers"
	"log"
	"os"
	"strconv"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	_ "github.com/themane/MMOServer/docs"
)

var once = sync.Once{}
var baseURL string
var mongoDB string
var maxSystems int

// @title MMO Game Server
// @version 1.0.0
// @description This is the server for new MMO Game based in space.
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host mmo-server-4xcaklgmnq-el.a.run.app
// @schemes https
func main() {
	r := gin.Default()

	once.Do(initialize)
	url := ginSwagger.URL(baseURL + "/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	loginController, buildingController, scheduledJobManager := getHandlers()

	r.GET("/ping", controllers.Ping)
	r.POST("/login", loginController.Login)
	r.POST("/refresh/population", loginController.RefreshPopulation)
	r.POST("/refresh/resources", loginController.RefreshResources)
	r.POST("/upgrade/building", buildingController.UpgradeBuilding)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
	scheduledJobManager.SchedulePlanetUpdates()
}

func getHandlers() (*controllers.LoginController, *controllers.BuildingController, *schedulers.ScheduledJobManager) {
	mongoURL := accessSecretVersion()
	client, ctx, cancel := mongoRepository.GetConnection(mongoURL)
	waterConstants := constants.GetWaterConstants()
	grapheneConstants := constants.GetGrapheneConstants()
	experienceConstants := constants.GetExperienceConstants()

	var userRepository models.UserRepository
	var clanRepository models.ClanRepository
	var universeRepository models.UniverseRepository
	userRepository = mongoRepository.NewUserRepository(client, ctx, cancel, mongoDB)
	clanRepository = mongoRepository.NewClanRepository(client, ctx, cancel, mongoDB)
	universeRepository = mongoRepository.NewUniverseRepository(client, ctx, cancel, mongoDB)
	loginController := controllers.NewLoginController(&userRepository, &clanRepository, &universeRepository, waterConstants, grapheneConstants, experienceConstants)
	buildingController := controllers.NewBuildingController(&userRepository)
	scheduledJobManager := schedulers.NewScheduledJobManager(&userRepository, &universeRepository, waterConstants, grapheneConstants, maxSystems)
	return loginController, buildingController, scheduledJobManager
}

func initialize() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	mongoDB = os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "mmo-game"
	}
	maxSystemsString := os.Getenv("MAX_SYSTEMS")
	if maxSystemsString == "" {
		maxSystemsString = "10"
	}
	var err error
	maxSystems, err = strconv.Atoi(maxSystemsString)
	if err != nil {
		log.Fatal(err)
	}
}

func accessSecretVersion() string {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/themane/secrets/MONGO_URL/versions/latest",
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(result.Payload.Data)
}

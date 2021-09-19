package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/themane/MMOServer/controllers"
	"github.com/themane/MMOServer/mongoRepository"
	"github.com/themane/MMOServer/mongoRepository/models"
	"github.com/themane/MMOServer/schedulers"
	"log"
	"os"
	"sync"

	_ "github.com/themane/MMOServer/docs"
)

var once = sync.Once{}
var baseURL string
var mongoURL string
var mongoDB string

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

	loginController, buildingController := getControllers()

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

	schedulers.SchedulePlanetUpdates()
}

func getControllers() (*controllers.LoginController, *controllers.BuildingController) {
	client, ctx, cancel := mongoRepository.GetConnection(mongoURL)
	var userRepository models.UserRepository
	var clanRepository models.ClanRepository
	var universeRepository models.UniverseRepository
	userRepository = mongoRepository.NewUserRepository(client, ctx, cancel, mongoDB)
	clanRepository = mongoRepository.NewClanRepository(client, ctx, cancel, mongoDB)
	universeRepository = mongoRepository.NewUniverseRepository(client, ctx, cancel, mongoDB)
	loginController := controllers.NewLoginController(&userRepository, &clanRepository, &universeRepository)
	buildingController := controllers.NewBuildingController(&userRepository)
	return loginController, buildingController
}

func initialize() {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	mongoURL := os.Getenv("MONGO_URL")
	mongoDB := os.Getenv("MONGO_DB")
	if mongoURL == "" || mongoDB == "" {
		log.Fatal("Error in retrieving mongo configurations")
	}
}

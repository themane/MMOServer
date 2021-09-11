package main

import (
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/controllers"
	"log"
	"os"
	"sync"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/themane/MMOServer/docs"
)

var once = sync.Once{}
var baseURL string

// @title MMO Game Server
// @version 1.0.0
// @description This is the server for new MMO Game based in space.
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://mmo-server-4xcaklgmnq-el.a.run.app
// @BasePath /
// @schemes https
func main() {
	r := gin.Default()

	once.Do(GetBaseURL)
	url := ginSwagger.URL(baseURL + "/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/ping", controllers.Ping)
	r.POST("/login", controllers.Login)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func GetBaseURL() {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
}

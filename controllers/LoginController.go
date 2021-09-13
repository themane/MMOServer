package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/services"
	"io/ioutil"
	"log"
)

// LoginController godoc
// @Summary Login API
// @Description Login verification and first load of complete user data
// @Tags data retrieval
// @Accept json
// @Produce json
// @Param username query string true "user identifier"
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func LoginController(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var request models.LoginRequest
	json.Unmarshal(body, &request)
	log.Printf("Logged in user: %s", request.Username)

	response := services.Login(request.Username)
	c.JSON(200, response)
}

package controllers

import "github.com/gin-gonic/gin"

// Ping godoc
// @Summary Pings the server for checking the health of the server
// @Description get the version and check the health of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test 3 Pong",
	})
}

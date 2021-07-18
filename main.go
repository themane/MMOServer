package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/ping", ping)
	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "branch push linked",
	})
}

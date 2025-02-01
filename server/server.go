package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Привет, WB и Димончик-братка!",
		})
	})

	statuses := []int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest, http.StatusInternalServerError}
	goodStatuses := 0
	badStatuses := 0

	// TODO: Try make 70 good/30 bad statuses.
	r.POST("/requests", func(c *gin.Context) {
		randStatus := statuses[rand.Intn(len(statuses))]
		if randStatus == http.StatusOK || randStatus == http.StatusAccepted {
			goodStatuses++
		} else {
			badStatuses++
		}
		c.JSON(randStatus, gin.H{})
	})
	r.Run(":" + port)
}

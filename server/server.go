// v0.1
//package main
//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/joho/godotenv"
//	"log"
//	"math/rand"
//	"net/http"
//	"os"
//)
//
//func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "8080"
//	}
//
//	r := gin.Default()
//	r.GET("/", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Привет, WB и Димончик-братка!",
//		})
//	})
//
//	statuses := []int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest, http.StatusInternalServerError}
//	goodStatuses := 0
//	badStatuses := 0
//
//	// TODO: Try make 70 good/30 bad statuses.
//	r.POST("/requests", func(c *gin.Context) {
//		randStatus := statuses[rand.Intn(len(statuses))]
//		if randStatus == http.StatusOK || randStatus == http.StatusAccepted {
//			goodStatuses++
//		} else {
//			badStatuses++
//		}
//		c.JSON(randStatus, gin.H{})
//	})
//	r.Run(":" + port)
//}

// v0.2
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

type Client struct {
	GoodClientReq uint `json:"Good client request"`
	BadClientReq  uint `json:"Bad client request"`
}

type ServerStatistics struct {
	GoodReqs    uint              `json:"Good request"`
	BadReqs     uint              `json:"Bad request"`
	ClientsReqs map[string]Client `json:"Clients requests"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mu := sync.Mutex{}
	stats := ServerStatistics{
		ClientsReqs: make(map[string]Client),
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Привет, WB и Димончик-братка!",
		})
	})

	statuses := []int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest, http.StatusInternalServerError}

	// TODO: Try make 70 good/30 bad statuses?
	r.POST("/requests", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		UUID := c.Request.Header.Get("UUID")
		randStatus := statuses[rand.Intn(len(statuses))]

		// Инициализации статистики клиента, если она отсутствует
		if _, exists := stats.ClientsReqs[UUID]; !exists {
			stats.ClientsReqs[UUID] = Client{}
		}

		client := stats.ClientsReqs[UUID]

		if randStatus == http.StatusOK || randStatus == http.StatusAccepted {
			stats.GoodReqs++
			client.GoodClientReq++
			stats.ClientsReqs[UUID] = client
		} else {
			stats.BadReqs++
			client.BadClientReq++
			stats.ClientsReqs[UUID] = client
		}

		c.JSON(randStatus, gin.H{})
	})

	r.GET("/statistics", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"statistics": stats,
		})
	})
	r.Run(":" + port)
}

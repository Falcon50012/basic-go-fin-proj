package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
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
		log.Fatal("Ошибка загрузки .env файла")
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
			"message": "Привет, WB!",
		})
	})

	ctx := context.Background()
	limiter := rate.NewLimiter(5, 1)

	r.POST("/requests", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		statuses := []int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest, http.StatusInternalServerError}

		if err = limiter.Wait(ctx); err != nil {
			log.Printf("Количество запросов превышено: %v", err)
		}

		UUID := c.Request.Header.Get("UUID")
		randStatus := statuses[rand.Intn(len(statuses))]
		
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

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/statistics", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"Statistics": stats,
		})
	})
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Ошибка при старте сервера: %v", err)
	}
}

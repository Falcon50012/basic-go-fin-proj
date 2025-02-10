package main

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

const UUID = "e9d5bf5b-077d-46f6-9bb4-122dd9f87077"

func main() {
	client := http.Client{}
	//UUID := uuid.New().String()
	//fmt.Println("UUID:", UUID)

	limiter := rate.NewLimiter(rate.Every(5*time.Second), 1)
	ctx := context.Background()

	for {
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Too many requests per second: %v", err)
			continue
		}

		request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/ping", nil)
		if err != nil {
			log.Printf("Request error: %v", err)
			continue // Важно не выполнять client.Do(nil)
		}

		request.Header.Add("UUID", UUID)

		response, err := client.Do(request)
		if err != nil {
			log.Printf("Ошибка при запросе: %v. Повторная попытка через 5 секунд...", err)
			continue
		}

		log.Printf("СТАТУС ОТВЕТА: %v", response.Status)

		// Закрываем тело ответа, чтобы избежать утечек соединений
		if response.Body != nil {
			response.Body.Close()
		}
	}
}

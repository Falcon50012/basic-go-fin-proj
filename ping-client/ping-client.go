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

	ctx := context.Background()
	limiter := rate.NewLimiter(rate.Every(5*time.Second), 1)

	for {
		if err := limiter.Wait(ctx); err != nil {
			continue
		}

		request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/ping", nil)
		if err != nil {
			log.Printf("Ошибка запроса:: %v", err)
			continue
		}

		request.Header.Add("UUID", UUID)

		response, err := client.Do(request)
		if err != nil {
			log.Printf("Ошибка ответа: %v. Повторная попытка через 5 секунд...", err)
			continue
		}

		log.Printf("СТАТУС ОТВЕТА: %v", response.Status)

		func() {
			err = response.Body.Close()
			if err != nil {
				return
			}
		}()
	}
}

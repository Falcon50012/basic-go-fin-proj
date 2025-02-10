package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"sync"
)

const (
	UUID          = "6e3c3a57-7f91-4fa5-b75a-5a7373c913bf"
	workers       = 2
	reqsPerWorker = 50
)

func main() {
	client := http.Client{}
	//UUID := uuid.New().String()
	//fmt.Println("UUID:", UUID)

	var mu sync.Mutex // Так правильно в отличие от mu := sync.Mutex{} ?
	requestsCounter := 0
	statusesCounter := map[int]int{
		http.StatusOK:                  0,
		http.StatusAccepted:            0,
		http.StatusBadRequest:          0,
		http.StatusInternalServerError: 0,
	}

	limiter := rate.NewLimiter(5, 5)
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 1; j <= reqsPerWorker; j++ {

				if err := limiter.Wait(ctx); err != nil {
					log.Printf("Too many requests per second: %v", err)
					continue
				}

				request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/requests", nil)
				if err != nil {
					log.Printf("Request error: %v", err)
				}

				request.Header.Add("UUID", UUID)

				response, err := client.Do(request)
				if err != nil {
					log.Printf("Response error: %v", err)
				}

				mu.Lock()
				requestsCounter++
				statusesCounter[response.StatusCode]++
				mu.Unlock()

				log.Printf("ВОРКЕР # %v СТАТУС ОТВЕТА: %v", workerID, response.Status)

			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("ОТПРАВЛЕНО ЗАПРОСОВ: %d\n", requestsCounter)
	fmt.Println("РАЗБИВКА ПО СТАТУСАМ:")
	fmt.Printf("StatusOK (200): %d\n", statusesCounter[http.StatusOK])
	fmt.Printf("StatusAccepted (202): %d\n", statusesCounter[http.StatusAccepted])
	fmt.Printf("StatusBadRequest (400): %d\n", statusesCounter[http.StatusBadRequest])
	fmt.Printf("StatusInternalServerError(500): %d\n", statusesCounter[http.StatusInternalServerError])
}

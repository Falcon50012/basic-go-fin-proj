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
	UUID          = "64387b8d-0db6-4c41-addb-6d593507bd89"
	workers       = 2
	reqsPerWorker = 50
)

func main() {
	client := http.Client{}

	var mu sync.Mutex
	requestsCounter := 0
	statusesCounter := map[int]int{
		http.StatusOK:                  0,
		http.StatusAccepted:            0,
		http.StatusBadRequest:          0,
		http.StatusInternalServerError: 0,
	}

	ctx := context.Background()
	limiter := rate.NewLimiter(5, 1)

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go func(workerID int) {
			defer wg.Done()

			for j := 1; j <= reqsPerWorker; j++ {
				if err := limiter.Wait(ctx); err != nil {
					log.Printf("Количество запросов превышено: %v", err)
					continue
				}

				request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/requests", nil)
				if err != nil {
					log.Printf("Ошибка запроса: %v", err)
					continue
				}

				request.Header.Add("UUID", UUID)

				response, err := client.Do(request)
				if err != nil {
					log.Printf("Ошибка ответа: %v", err)
					continue
				}

				func() {
					defer response.Body.Close()

					mu.Lock()
					requestsCounter++
					statusesCounter[response.StatusCode]++
					mu.Unlock()
				}()

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
	fmt.Printf("StatusInternalServerError (500): %d\n", statusesCounter[http.StatusInternalServerError])
}

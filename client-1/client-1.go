package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	//requests      = 100
	workers       = 2
	reqsPerWorker = 5
)

func main() {
	client := http.Client{}

	mu := sync.Mutex{} // Is it worth using one mutex to control multiple entities?

	requestsCounter := 0
	statusesCounter := map[int]int{
		http.StatusOK:                  0,
		http.StatusAccepted:            0,
		http.StatusBadRequest:          0,
		http.StatusInternalServerError: 0,
	}

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 1; j <= reqsPerWorker; j++ {
				request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/requests", nil)
				if err != nil {
					log.Printf("Request error: %v", err)
				}

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

package main

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

const DURATION = 420
const COUNT = 21000
const TEST_URL = ""
const X_API_KEY = ""

func main() {
	fmt.Printf("Starting test api during: %v, request count: %v \n", DURATION, COUNT)

	var statuses []int
	var wg sync.WaitGroup
	wg.Add(COUNT)
	requestChannel := make(chan int)
	done := make(chan bool)

	reqPerSec := int(math.Ceil(COUNT / DURATION)) // 50

	fmt.Printf("Requests per second: %v \n", reqPerSec)

	ticker := time.NewTicker(time.Second)

	tickCount := DURATION

	go func() {
		wg.Wait()
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Printf("====STOP TICKER!\n")

			ticker.Stop()

			failed := 0
			for i, _ := range statuses {
				if statuses[i] != 200 {
					fmt.Println("Failed: ", statuses[i])
					failed++
				}
			}

			fmt.Printf("Total responses: %v \n", len(statuses))
			fmt.Printf("Failed: %v \n", failed)
			return
		case <-ticker.C:
			if tickCount == 0 {
				ticker.Stop()
				fmt.Println("****TICK STOP")
				break
			}
			fmt.Println("****TICK")
			for i := 0; i != reqPerSec; i++ {
				fmt.Println("Make requests...", i)

				go makeRequest(TEST_URL, X_API_KEY, requestChannel, &wg)
			}
			tickCount--
		case status := <-requestChannel:
			fmt.Println("Received status: ", status)
			statuses = append(statuses, status)
		}
	}
}

func makeRequest(url string, apiKey string, requestChannel chan int, wg *sync.WaitGroup) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	req.Header = http.Header{
		"X-Api-Key": []string{apiKey},
	}

	res, err := client.Do(req)
	if err != nil {
		requestChannel <- 500
		wg.Done()
		return
	}

	requestChannel <- res.StatusCode
	wg.Done()

	return
}

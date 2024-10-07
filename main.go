package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ApiFetcher struct {
	url        string
	statusCode string
}

type DataWritter struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func fetchURL(url string, results chan<- ApiFetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	retries := 3
	for i := 0; i < retries; i++ {
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			results <- ApiFetcher{url: url, statusCode: "Error"}
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close()
		results <- ApiFetcher{url: url, statusCode: resp.Status}
		return
	}

	results <- ApiFetcher{url: url, statusCode: "Error"}

}

func csvWritter(url string, results chan<- DataWritter, apiResult ApiFetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	if apiResult.statusCode == "Error" {
		results <- DataWritter{ID: 0, Title: "Error", Completed: false}
		return
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		results <- DataWritter{ID: 0, Title: "Error", Completed: false}
		return
	}
	defer resp.Body.Close()

	var data DataWritter
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		results <- DataWritter{ID: 0, Title: "Error", Completed: false}
		return
	}

	results <- data

}

func main() {
	var wg sync.WaitGroup
	urls := []string{
		"https://jsonplaceholder.typicode.com/todos/1",
		"https://jsonplaceholder.typicode.com/todos/2",
		"https://jsonplaceholder.typicode.com/todos/3",
		"https://jsonplaceholder.typicode.com/todos/4",
		"https://jsonplaceholder.typicode.com/todos/5",
	}

	apiResults := make(chan ApiFetcher, len(urls))
	csvResults := make(chan DataWritter, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, apiResults, &wg)
	}

	go func() {
		wg.Wait()
		close(apiResults)
	}()

	for apiResult := range apiResults {
		wg.Add(1)
		go csvWritter(apiResult.url, csvResults, apiResult, &wg)
	}

	go func() {
		wg.Wait()
		close(csvResults)
	}()

	for csvResult := range csvResults {
		fmt.Printf("ID: %d, Title: %s, Completed: %v\n", csvResult.ID, csvResult.Title, csvResult.Completed)
	}
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// define CLI flags
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	// validate flags
	if *url == "" {
		fmt.Println("a URL é obrigatória")
		return
	}

	results := make(chan int, *requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	// launch goroutines
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < *requests / *concurrency; j++ {
				resp, err := http.Get(*url)
				if err != nil {
					results <- 0
					continue
				}
				results <- resp.StatusCode
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	close(results)

	// calculate results
	elapsedTime := time.Since(startTime)
	successCount := 0
	statusCounts := make(map[int]int)
	for status := range results {
		if status == http.StatusOK {
			successCount++
		}
		statusCounts[status]++
	}

	// report
	fmt.Printf("Tempo total: %v\n", elapsedTime)
	fmt.Printf("Total de requests: %d\n", *requests)
	fmt.Printf("Requests com status 200: %d\n", successCount)
	fmt.Println("Distribuição de status:")
	for status, count := range statusCounts {
		fmt.Printf("Status %d: %d\n", status, count)
	}
}

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/bianavic/fullcycle_desafios/internal"
)

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	StatusCodes   map[int]int
}

var (
	url         string
	requests    int
	concurrency int
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Client      internal.HTTPClient
}

func worker(cfg Config, wg *sync.WaitGroup, codes chan<- int) {
	defer wg.Done()

	for i := 0; i < cfg.Requests; i++ {
		resp, err := cfg.Client.Get(cfg.URL)
		if err != nil {
			codes <- 0
			continue
		}

		codes <- resp.StatusCode
		resp.Body.Close()
	}
}

var rootCmd = &cobra.Command{
	Use:   "fullcycle_desafios stress test",
	Short: "A CLI tool to stress test a web service",
	Long:  `A CLI tool to stress test a web service simulating concurrent HTTP requests`,
	Run: func(cmd *cobra.Command, args []string) {
		client := &internal.RealHTTPClient{Client: http.DefaultClient}
		cfg := Config{
			URL:         url,
			Requests:    requests,
			Concurrency: concurrency,
			Client:      client,
		}

		fmt.Println("Making HTTP requests, please wait...")

		codes := make(chan int, cfg.Requests)
		workerRequests := cfg.Requests / cfg.Concurrency
		extraRequests := cfg.Requests % cfg.Concurrency

		startTime := time.Now()
		var wg sync.WaitGroup

		for i := 0; i < cfg.Concurrency; i++ {
			wg.Add(1)
			r := workerRequests
			if i < extraRequests {
				r++
			}
			go worker(cfg, &wg, codes)
		}

		wg.Wait()
		close(codes)

		report := Report{
			TotalTime:     time.Since(startTime),
			TotalRequests: cfg.Requests,
			StatusCodes:   make(map[int]int),
		}

		for code := range codes {
			report.StatusCodes[code]++
		}

		fmt.Println("+-------------------------------------------------------------+")
		fmt.Println("Stress Test Report:")
		fmt.Printf("Total time: %v\n", report.TotalTime)
		fmt.Printf("Total requests: %d\n", report.TotalRequests)
		fmt.Printf("Requests with status code 200: %d\n", report.StatusCodes[200])
		fmt.Printf("Requests by status code:\n")
		for code, count := range report.StatusCodes {
			fmt.Printf(" → %d %s: %d\n", code, http.StatusText(code), count)
		}
		fmt.Println("+-------------------------------------------------------------+")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "URL of the service to test")
	rootCmd.PersistentFlags().IntVarP(&requests, "requests", "r", 0, "total number of requests")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 0, "number of concurrent requests")

	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("requests")
	rootCmd.MarkPersistentFlagRequired("concurrency")
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
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

func worker(url string, requests int, wg *sync.WaitGroup, codes chan<- int) {
	defer wg.Done()

	for i := 0; i < requests; i++ {
		resp, err := http.Get(url)
		if err != nil {
			codes <- 0
			continue
		}

		codes <- resp.StatusCode

		resp.Body.Close()
	}
}

type RunEFunc func(cmd *cobra.Command, args []string) error

var rootCmd = &cobra.Command{
	Use:   "fullcycle_desafios stress test",
	Short: "A CLI tool to stress test a web service",
	Long:  `A CLI tool to stress test a web service simulating concurrent HTTP requests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Making HTTP requests, please wait...")

		codes := make(chan int, requests)

		workerRequests := requests / concurrency
		extraRequests := requests % concurrency

		startTime := time.Now()

		var wg sync.WaitGroup

		for i := 0; i < concurrency; i++ {
			wg.Add(1)

			r := workerRequests
			if i < extraRequests {
				r++
			}

			go worker(url, r, &wg, codes)
		}

		wg.Wait()

		close(codes)

		report := Report{
			TotalTime:     time.Since(startTime),
			TotalRequests: len(codes),
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

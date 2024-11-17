package main

import (
	"context"
	"fmt"
	"github.com/bianavic/fullcycle_desafios.git/client"
	"github.com/bianavic/fullcycle_desafios.git/server"
	"net/http"
	"time"
)

const (
	serverPort = ":8080"
	timeout    = 300 * time.Millisecond
)

func main() {

	// Start the server
	startServer()

	// Allow the server some time to start before the client makes a request
	time.Sleep(1 * time.Second)

	// Get the exchange rate from the local server
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rate, err := client.GetExchangeRate(ctx)
	if err != nil {
		fmt.Printf("Error getting exchange rate: %v\n", err)
		return
	}

	// Save the rate to a file
	if err := client.SaveToFile(rate); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
		return
	}

	fmt.Println("Exchange rate saved to cotacao.txt")

}

// startServer starts an HTTP server that serves exchange rates
func startServer() {
	http.HandleFunc("/cotacao", server.ExchangeRateHandler)
	fmt.Printf("Server running on http://localhost%s/cotacao\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

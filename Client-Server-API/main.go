package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bianavic/fullcycle_desafios.git/client"
	"github.com/bianavic/fullcycle_desafios.git/server"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

const (
	serverPort              = ":8080"
	timeout                 = 300 * time.Millisecond
	createExchangeRateTable = `CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
)

func main() {

	// Start the server
	startServer()

	// Allow the server some time to start before the client makes a request
	time.Sleep(10 * time.Second)

	// Get the exchange rate from the local server
	ctx, cancel := context.WithTimeout(context.Background(), 3*timeout)
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

func InitDB() {
	// Initialize the SQLite database and create the table
	db, err := sql.Open("sqlite3", "root:root@tcp(localhost:3306)/goexpert_challenge1")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the exchange_rates table if it doesn't exist
	_, err = db.Exec(createExchangeRateTable)
	if err != nil {
		fmt.Printf("Error creating table in SQLite: %v\n", err)
		return
	}
}

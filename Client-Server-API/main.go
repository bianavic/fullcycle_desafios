package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bianavic/fullcycle_desafios.git/client"
	"github.com/bianavic/fullcycle_desafios.git/server"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

const (
	dbTimeout               = 10 * time.Millisecond // Timeout for the database operation (10ms)
	dbFile                  = "exchange_rates.db"   // SQLite database file
	serverPort              = ":8080"
	createExchangeRateTable = `CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
)

func main() {

	// Initialize the database
	InitDB()

	// start the server
	go startServer()

	// allow the server some time to start before the client makes a request
	time.Sleep(10 * time.Second)

	// get the exchange rate from the local server
	ctx, cancel := context.WithTimeout(context.Background(), 3*dbTimeout)
	defer cancel()

	rate, err := client.GetExchangeRate(ctx)
	if err != nil {
		fmt.Printf("error getting exchange rate: %v\n", err)
		return
	}

	// save the rate to a file
	if err := client.SaveToFile(rate); err != nil {
		fmt.Printf("error saving to file: %v\n", err)
		return
	}

	fmt.Println("exchange rate saved to cotacao.txt")
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
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Create the exchange_rates table if it doesn't exist
	_, err = db.Exec(createExchangeRateTable)
	if err != nil {
		log.Fatalf("Error creating table in SQLite: %v", err)
	}
}

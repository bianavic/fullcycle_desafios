package main

import (
	"context"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/handlers"
	"os"
	"time"
)

const apiTimeout = 1 * time.Second

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Add the zip code to the command: go run main.go <zip code>")
		return
	}
	cep := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	resultChan := make(chan handlers.Result)

	go handlers.FetchAddressHandler(ctx, cep, "BrasilAPI", resultChan)
	go handlers.FetchAddressHandler(ctx, cep, "ViaCepAPI", resultChan)

	select {
	case result := <-resultChan:
		if result.Error != nil {
			fmt.Printf("Error fetching address: %v\n", result.Error)
		}
		handlers.PrintResult(result)
	case <-time.After(apiTimeout):
		fmt.Println("Timeout")
	}
}

package main

import (
	"github.com/bianavic/fullcycle_desafios/config"
	"github.com/bianavic/fullcycle_desafios/internal/ratelimit"
	"github.com/bianavic/fullcycle_desafios/internal/storage"
	"log"
	"net/http"
	"time"
)

func main() {
	// load environment variables
	if err := config.LoadConfig(); err != nil {
		log.Fatal("error loading .env file")
	}

	// initialize Redis storage
	redisStorage, err := storage.NewRedis(
		config.GetEnv("REDIS_ADDR", "localhost:6379"),
		config.GetEnv("REDIS_PASSWORD", ""))
	if err != nil {
		log.Fatalf("failed to initialize Redis storage: %v", err)
	}

	// parse rate limit and block time from environment
	rateLimitIP := config.GetIntEnv("RATE_LIMIT_IP", 10)
	rateLimitToken := config.GetIntEnv("RATE_LIMIT_TOKEN", 100)
	blockTime := config.GetDurationEnv("BLOCK_TIME", 60*time.Second)

	// initialize rate limiter
	limiter := ratelimit.NewRateLimiter(redisStorage, rateLimitIP, rateLimitToken, blockTime)

	// create HTTP server with rate limiter middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Stranger!"))
	})
	// TODO update port to 8080
	server := &http.Server{
		Addr:    ":8082",
		Handler: ratelimit.RateLimiterMiddleware(limiter)(mux),
	}

	log.Println("server started on :8082")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

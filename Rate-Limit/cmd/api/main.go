package main

import (
	"log"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/infra/config"
	limiter "github.com/bianavic/fullcycle_desafios/internal/infra/limiter"
	"github.com/bianavic/fullcycle_desafios/internal/infra/middleware"
)

func main() {
	// load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// initialize storage (Redis or in-memory)
	var storageStrategy limiter.StorageStrategy
	if cfg.UseRedis {
		redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
		redisStorage, err := limiter.NewRedis(redisAddr, cfg.RedisPassword)
		if err != nil {
			log.Fatalf("Failed to initialize Redis storage: %v", err)
		}
		storageStrategy = redisStorage
	} else {
		storageStrategy = limiter.NewInMemory()
	}

	// Convert config.TokenConfigs to usecase.TokenConfigs
	tokenConfigs := make(map[string]limiter.TokenConfig)
	for k, v := range cfg.TokenConfigs {
		tokenConfigs[k] = limiter.TokenConfig{
			RateLimit: v.RateLimit,
			BlockTime: v.BlockTime,
		}
	}

	// Initialize rate limiter
	limiter := limiter.NewRateLimiter(storageStrategy, cfg.RateLimitIP, cfg.BlockTime, tokenConfigs)

	// create HTTP server with rate limiter middleware
	http.Handle("/", middleware.RateLimiterMiddleware(limiter, http.HandlerFunc(handler)))

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Stranger!\n"))
}

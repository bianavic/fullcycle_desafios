package main

import (
	"log"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/config"
	"github.com/bianavic/fullcycle_desafios/internal/middleware"
	"github.com/bianavic/fullcycle_desafios/internal/repository/storage"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
)

func main() {
	// load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// initialize storage (Redis or in-memory)
	var storageStrategy storage.StorageStrategy
	if cfg.UseRedis {
		redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
		redisStorage, err := storage.NewRedis(redisAddr, cfg.RedisPassword)
		if err != nil {
			log.Fatalf("Failed to initialize Redis storage: %v", err)
		}
		storageStrategy = redisStorage
	} else {
		storageStrategy = storage.NewInMemory()
	}

	// Convert config.TokenConfigs to usecase.TokenConfigs
	tokenConfigs := make(map[string]usecase.TokenConfig)
	for k, v := range cfg.TokenConfigs {
		tokenConfigs[k] = usecase.TokenConfig{
			RateLimit: v.RateLimit,
			BlockTime: v.BlockTime,
		}
	}

	// Initialize rate limiter
	limiter := usecase.NewRateLimiter(storageStrategy, cfg.RateLimitIP, cfg.BlockTime, tokenConfigs)

	// create HTTP server with rate limiter middleware
	http.Handle("/", middleware.RateLimiterMiddleware(limiter, http.HandlerFunc(handler)))

	log.Println("server started on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Stranger!\n"))
}

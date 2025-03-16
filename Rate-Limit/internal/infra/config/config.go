package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RateLimitIP   int
	BlockTime     time.Duration
	UseRedis      bool
	TokenConfigs  map[string]TokenConfig
}

type TokenConfig struct {
	RateLimit int
	BlockTime time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	rateLimitIP, err := strconv.Atoi(getEnv("RATE_LIMIT_IP", "5"))
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT_IP: %v", err)
	}

	blockTime, err := time.ParseDuration(getEnv("BLOCK_TIME", "60s"))
	if err != nil {
		log.Fatalf("Invalid BLOCK_TIME: %v", err)
	}

	useRedis := getEnv("USE_REDIS", "true") == "true"

	tokenConfigs := map[string]TokenConfig{
		getEnv("TOKEN1", "abc123"): {
			RateLimit: getIntEnv("TOKEN1_LIMIT", 10),
			BlockTime: getDurationEnv("TOKEN1_EXPIRATION", 60*time.Second),
		},
		getEnv("TOKEN2", "def456"): {
			RateLimit: getIntEnv("TOKEN2_LIMIT", 20),
			BlockTime: getDurationEnv("TOKEN2_EXPIRATION", 65*time.Second),
		},
	}

	return &Config{
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RateLimitIP:   rateLimitIP,
		BlockTime:     blockTime,
		UseRedis:      useRedis,
		TokenConfigs:  tokenConfigs,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}

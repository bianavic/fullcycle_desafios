package config

import (
	"os"
)

type Config struct {
	Port        string
	ServiceBURL string
	Env         string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8081"),
		ServiceBURL: getEnv("SERVICE_B_URL", "http://service-b:8082"),
		Env:         getEnv("ENV", "development"),
	}
}

// getEnv helper function with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

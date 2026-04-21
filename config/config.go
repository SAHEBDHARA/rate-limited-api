package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIPort    int
	RedisAddr  string
	RedisPass  string
	RateLimit  int
}

func LoadConfig() Config {
	return Config{
		APIPort:    getEnvAsInt("API_PORT", 8080),
		RedisAddr:  getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPass:  getEnv("REDIS_PASS", ""),
		RateLimit:  getEnvAsInt("RATE_LIMIT", 5),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

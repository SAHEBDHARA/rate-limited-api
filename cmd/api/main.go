package main

import (
	"log"
	"time"

	"rate-limited-api/config"
	"rate-limited-api/internal/adapters/http"
	redisAdapter "rate-limited-api/internal/adapters/redis"
	"rate-limited-api/internal/core/services"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Redis Client using the adapter
	rdb, err := redisAdapter.NewRedisClient(cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		log.Fatalf("Redis initialization failed: %v", err)
	}
	log.Println("Connected to Redis at", cfg.RedisAddr)

	// 3. Initialize Adapters (Secondary/Driven)
	statsRepo := redisAdapter.NewStatsRepository(rdb)
	rateLimiter := redisAdapter.NewRateLimiter(rdb, cfg.RateLimit, time.Minute)

	// 4. Initialize Core Domain Services
	requestService := services.NewRequestService(statsRepo)

	// 5. Initialize and run Driving Adapter (The HTTP Server)
	server := http.NewServer(cfg.APIPort, requestService, rateLimiter)

	log.Printf("Starting API Server on port %d...", cfg.APIPort)
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

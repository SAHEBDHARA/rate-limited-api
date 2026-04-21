package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"rate-limited-api/internal/core/ports"
)

type rateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter creates a new Redis-backed fixed-window rate limiter.
func NewRateLimiter(client *redis.Client, limit int, window time.Duration) ports.RateLimiter {
	return &rateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// Allow implements the RateLimiter interface.
func (r *rateLimiter) Allow(ctx context.Context, userID string) (bool, error) {
	// We use the current Unix time divided by the window length to get a unique key per time slice.
	// e.g. if window is 60s, this integer changes exactly every 60 seconds.
	windowKey := time.Now().Unix() / int64(r.window.Seconds())
	key := fmt.Sprintf("ratelimit:%s:%d", userID, windowKey)
	
	// Atomic increment ensures this is perfectly highly concurrent-safe.
	current, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis incr error: %w", err)
	}

	// Set expiration only on the first request of the window.
	// We pad the expiration by a few seconds to ensure it cleans up safely.
	if current == 1 {
		err = r.client.Expire(ctx, key, r.window+time.Second*5).Err()
		if err != nil {
			return false, fmt.Errorf("redis expire error: %w", err)
		}
	}

	if current > int64(r.limit) {
		return false, nil // Limit exceeded
	}

	return true, nil
}

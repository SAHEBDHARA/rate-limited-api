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

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) ports.RateLimiter {
	return &rateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (r *rateLimiter) Allow(ctx context.Context, userID string) (bool, error) {
	windowKey := time.Now().Unix() / int64(r.window.Seconds())
	key := fmt.Sprintf("ratelimit:%s:%d", userID, windowKey)
	
	current, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis incr error: %w", err)
	}

	if current == 1 {
		err = r.client.Expire(ctx, key, r.window+time.Second*5).Err()
		if err != nil {
			return false, fmt.Errorf("redis expire error: %w", err)
		}
	}

	if current > int64(r.limit) {
		return false, nil 
	}

	return true, nil
}

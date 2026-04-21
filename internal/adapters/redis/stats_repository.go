package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"rate-limited-api/internal/core/domain"
	"rate-limited-api/internal/core/ports"
)

type statsRepository struct {
	client *redis.Client
}

func NewStatsRepository(client *redis.Client) ports.StatsRepository {
	return &statsRepository{
		client: client,
	}
}

func (r *statsRepository) Increment(ctx context.Context, userID string) error {
	key := fmt.Sprintf("stats:%s", userID)
	return r.client.Incr(ctx, key).Err()
}

func (r *statsRepository) GetStats(ctx context.Context, userID string) (domain.UserStats, error) {
	key := fmt.Sprintf("stats:%s", userID)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.UserStats{UserID: userID, TotalRequests: 0}, nil
		}
		return domain.UserStats{}, fmt.Errorf("redis get error: %w", err)
	}

	total, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return domain.UserStats{}, fmt.Errorf("failed to parse stats value: %w", err)
	}

	return domain.UserStats{
		UserID:        userID,
		TotalRequests: total,
	}, nil
}

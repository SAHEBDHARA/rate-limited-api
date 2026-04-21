package ports

import (
	"context"

	"rate-limited-api/internal/core/domain"
)

type StatsRepository interface {
	Increment(ctx context.Context, userID string) error
	GetStats(ctx context.Context, userID string) (domain.UserStats, error)
}

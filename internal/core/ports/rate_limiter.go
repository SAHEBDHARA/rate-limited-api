package ports

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, userID string) (bool, error)
}

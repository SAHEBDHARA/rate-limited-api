package ports

import (
	"context"

	"rate-limited-api/internal/core/domain"
)

type RequestService interface {
	ProcessRequest(ctx context.Context, req domain.Request) error
	GetStats(ctx context.Context, userID string) (domain.UserStats, error)
}

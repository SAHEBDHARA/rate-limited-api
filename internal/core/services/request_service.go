package services

import (
	"context"
	"fmt"

	"rate-limited-api/internal/core/domain"
	"rate-limited-api/internal/core/ports"
)

type requestService struct {
	statsRepo ports.StatsRepository
}

func NewRequestService(statsRepo ports.StatsRepository) ports.RequestService {
	return &requestService{
		statsRepo: statsRepo,
	}
}

func (s *requestService) ProcessRequest(ctx context.Context, req domain.Request) error {
	if req.Payload == "" {
		return fmt.Errorf("payload cannot be empty")
	}

	err := s.statsRepo.Increment(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to increment stats: %w", err)
	}

	return nil
}
func (s *requestService) GetStats(ctx context.Context, userID string) (domain.UserStats, error) {
	if userID == "" {
		return domain.UserStats{}, fmt.Errorf("user_id is required")
	}
	
	return s.statsRepo.GetStats(ctx, userID)
}

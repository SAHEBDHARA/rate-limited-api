package domain

// UserStats represents the statistics for a specific user.
type UserStats struct {
	UserID        string `json:"user_id"`
	TotalRequests int64  `json:"total_requests"`
}

package domain

import (
	"errors"
	"time"
)

type Request struct {
	ID        string
	UserID    string
	Payload   string
	CreatedAt time.Time
}

func NewRequest(id, userID, payload string) (Request, error) {
	if userID == "" {
		return Request{}, errors.New("user_id must not be empty")
	}
	if id == "" {
		return Request{}, errors.New("request id must not be empty")
	}
	return Request{
		ID:        id,
		UserID:    userID,
		Payload:   payload,
		CreatedAt: time.Now(),
	}, nil
}

type Job struct {
	RequestID string
	UserID    string
	Payload   string
	Attempts  int
}

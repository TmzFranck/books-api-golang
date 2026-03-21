package jobs

import (
	"context"
	"time"
)

// Job 
type Job struct {
	ID string
	Type string
	Payload map[string]any
	Attempts int
	MaxRetry int
	CreatedAt time.Time
	Error string
}

// Handler 
type Handler func(ctx context.Context, job *Job) error



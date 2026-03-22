package jobs

import (
	"context"
	"time"
)

// Job represents a job to be processed by the worker pool
type Job struct {
	ID        string
	Type      string
	Payload   map[string]any
	Attempts  int
	MaxRetry  int
	CreatedAt time.Time
	Error     string
}

// Handler is a function that processes a job
type Handler func(ctx context.Context, job *Job) error

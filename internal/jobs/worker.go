package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// WorkerPool managers pool works that process jobs
type WorkerPool struct {
	numWorkers int
	jobQueue   chan *Job
	handlers   map[string]Handler
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	logger     *logrus.Logger
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(numWorkers, queueSize int, logger *logrus.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan *Job, queueSize),
		handlers:   make(map[string]Handler),
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
	}
}

// RegisterHandler registers a handler for a specific job type
func (wp *WorkerPool) RegisterHandler(JobType string, handler Handler) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.handlers[JobType] = handler
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes jobs from the job queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			wp.logger.Infof("Worker %d shutting down", id)
			return
		case job := <-wp.jobQueue:
			wp.processJob(id, job)
		}
	}
}

// processJob processes a single job
func (wp *WorkerPool) processJob(workerID int, job *Job) {
	wp.mu.RLock()
	handler, exists := wp.handlers[job.Type]
	wp.mu.RUnlock()

	if !exists {
		wp.logger.Errorf("Worker %d: no handler for job type %s", workerID, job.Type)
	}

	ctx, cancel := context.WithTimeout(wp.ctx, 30*time.Second)
	defer cancel()

	wp.logger.Infof("Worker %d: processing job %s (type: %s, attempt: %d)", workerID, job.ID, job.Type, job.Attempts+1)

	err := wp.safeExecute(ctx, handler, job)

	if err != nil {
		job.Error = err.Error()
		job.Attempts++

		if job.Attempts < job.MaxRetry {
			go wp.scheduleRetry(job)
		} else {
			wp.logger.Errorf("Worker %d: job %s failed permanently: %v", workerID, job.ID, err)
		}
	} else {
		wp.logger.Infof("Worker %d: job %s completed successfully", workerID, job.ID)
	}
}

// safeExecute executes a job handler safely, recovering from panics
func (wp *WorkerPool) safeExecute(ctx context.Context, handler Handler, job *Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	return handler(ctx, job)
}

// Submit submits a job to the worker pool for processing
func (wp *WorkerPool) Submit(ctx context.Context, job *Job) error {
	if job.ID == "" {
		job.ID = generateID()
	}

	if job.MaxRetry == 0 {
		job.MaxRetry = 3 // default retry count
	}

	job.CreatedAt = time.Now()

	select {
	case <-ctx.Done():
		return fmt.Errorf("submission cancelled: %w", ctx.Err())
	case wp.jobQueue <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

// scheduleRetry schedules a retry for a failed job
func (wp *WorkerPool) scheduleRetry(job *Job) {
	delay := time.Duration(1<<job.Attempts) * time.Second

	wp.logger.Infof("Scheduling retry for job %s in %v", job.ID, delay)

	select {
	case <-time.After(delay):
		select {
		case wp.jobQueue <- job:
			wp.logger.Infof("Job %s requeued for retry", job.ID)
		case <-wp.ctx.Done():
			return
		}
	case <-wp.ctx.Done():
		return
	}
}

// generateID generates a unique job ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Shutdown
func (wp *WorkerPool) Shutdown(timeout time.Duration) error {
	wp.logger.Infof("Initiating graceful schutdown...")

	wp.cancel()

	done := make(chan struct{})

	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		wp.logger.Infof("All workers stopped gracefully")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timed out after %v", timeout)
	}
}

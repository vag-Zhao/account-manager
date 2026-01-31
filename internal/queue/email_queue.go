package queue

import (
	"fmt"
	"sync"
	"time"

	"account-manager/internal/config"
)

// EmailJob represents an email to be sent
type EmailJob struct {
	ID        string
	Subject   string
	Content   string
	Timestamp time.Time
	Retries   int
	MaxRetries int
	ResultChan chan error
}

// EmailQueue manages asynchronous email sending
type EmailQueue struct {
	jobs       chan *EmailJob
	workers    int
	wg         sync.WaitGroup
	sendFunc   func(subject, content string) error
	stopChan   chan struct{}
	isRunning  bool
	mu         sync.Mutex
}

// NewEmailQueue creates a new email queue
func NewEmailQueue(workers int, sendFunc func(subject, content string) error) *EmailQueue {
	cfg := config.Get()
	return &EmailQueue{
		jobs:       make(chan *EmailJob, cfg.Worker.EmailQueueSize),
		workers:    workers,
		sendFunc:   sendFunc,
		stopChan:   make(chan struct{}),
		isRunning:  false,
	}
}

// Start starts the email queue workers
func (q *EmailQueue) Start() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isRunning {
		return
	}

	q.isRunning = true

	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// Stop stops the email queue
func (q *EmailQueue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isRunning {
		return
	}

	close(q.stopChan)
	q.wg.Wait()
	q.isRunning = false
}

// worker processes email jobs
func (q *EmailQueue) worker(id int) {
	defer q.wg.Done()

	for {
		select {
		case <-q.stopChan:
			return
		case job := <-q.jobs:
			q.processJob(job)
		}
	}
}

// processJob sends an email with retry logic
func (q *EmailQueue) processJob(job *EmailJob) {
	var err error

	for attempt := 0; attempt <= job.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2^attempt seconds
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			time.Sleep(backoff)
		}

		// Add timeout protection
		done := make(chan error, 1)
		go func() {
			done <- q.sendFunc(job.Subject, job.Content)
		}()

		select {
		case err = <-done:
			if err == nil {
				// Success
				job.ResultChan <- nil
				return
			}
		case <-time.After(30 * time.Second):
			err = fmt.Errorf("邮件发送超时")
		}

		job.Retries = attempt + 1
	}

	// All retries failed
	job.ResultChan <- fmt.Errorf("邮件发送失败 (重试 %d 次): %v", job.MaxRetries, err)
}

// Enqueue adds an email job to the queue
func (q *EmailQueue) Enqueue(subject, content string) <-chan error {
	job := &EmailJob{
		ID:         fmt.Sprintf("%d", time.Now().UnixNano()),
		Subject:    subject,
		Content:    content,
		Timestamp:  time.Now(),
		MaxRetries: 3,
		ResultChan: make(chan error, 1),
	}

	q.jobs <- job
	return job.ResultChan
}

// GetQueueSize returns the current queue size
func (q *EmailQueue) GetQueueSize() int {
	return len(q.jobs)
}

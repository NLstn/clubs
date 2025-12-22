package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/models"
	"gorm.io/gorm"
)

// JobFunc is a function that can be executed as a scheduled job
type JobFunc func() error

// JobConfig contains the configuration for a scheduled job
type JobConfig struct {
	// Name is a unique identifier for the job in the database
	Name string
	// Description provides a human-readable summary of what the job does
	Description string
	// IntervalMinutes specifies how often the job should run, in minutes.
	// Must be a positive integer value representing the interval between executions.
	IntervalMinutes int
}

// Validate checks if the JobConfig has valid values
func (c JobConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("job name cannot be empty")
	}
	if c.IntervalMinutes <= 0 {
		return fmt.Errorf("interval minutes must be positive, got %d", c.IntervalMinutes)
	}
	return nil
}

// Scheduler manages periodic job execution
type Scheduler struct {
	jobs       map[string]JobFunc
	jobsMutex  sync.RWMutex
	ticker     *time.Ticker
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	tickerInterval time.Duration
}

// NewScheduler creates a new scheduler instance
// tickerInterval determines how often the scheduler checks for jobs to run (default: 1 minute)
func NewScheduler(tickerInterval time.Duration) *Scheduler {
	if tickerInterval == 0 {
		tickerInterval = 1 * time.Minute
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scheduler{
		jobs:           make(map[string]JobFunc),
		ctx:            ctx,
		cancel:         cancel,
		tickerInterval: tickerInterval,
	}
}

// RegisterJobWithSchedule registers a job handler and creates/updates its database record
// This method combines in-memory registration with database initialization in a single operation
func (s *Scheduler) RegisterJobWithSchedule(handlerName string, jobFunc JobFunc, config JobConfig) error {
	// Validate input parameters
	if handlerName == "" {
		return fmt.Errorf("handler name cannot be empty")
	}
	if jobFunc == nil {
		return fmt.Errorf("job function cannot be nil")
	}
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid job config: %w", err)
	}
	
	// Step 1: Create/update database record first (for atomicity)
	var existing models.ScheduledJob
	err := database.Db.Where("name = ?", config.Name).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new job record
		job := models.ScheduledJob{
			Name:            config.Name,
			Description:     config.Description,
			JobHandler:      handlerName,
			IntervalMinutes: config.IntervalMinutes,
			Enabled:         true,
		}
		if err := database.Db.Create(&job).Error; err != nil {
			return fmt.Errorf("failed to create scheduled job '%s': %w", config.Name, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to query scheduled job '%s': %w", config.Name, err)
	} else {
		// Update existing job if handler or interval changed
		updates := make(map[string]interface{})
		if existing.JobHandler != handlerName {
			updates["job_handler"] = handlerName
		}
		if existing.IntervalMinutes != config.IntervalMinutes {
			updates["interval_minutes"] = config.IntervalMinutes
		}
		if existing.Description != config.Description {
			updates["description"] = config.Description
		}
		
		if len(updates) > 0 {
			if err := database.Db.Model(&existing).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update scheduled job '%s': %w", config.Name, err)
			}
		}
	}
	
	// Step 2: Register in-memory only after successful database operation
	s.jobsMutex.Lock()
	s.jobs[handlerName] = jobFunc
	s.jobsMutex.Unlock()
	
	log.Printf("Registered scheduled job: %s", config.Name)
	
	return nil
}

// Start begins the scheduler's execution loop
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.tickerInterval)
	s.wg.Add(1)
	
	go func() {
		defer s.wg.Done()
		log.Println("Scheduler started successfully")
		
		// Run immediately on start to catch any overdue jobs
		s.checkAndRunJobs()
		
		for {
			select {
			case <-s.ticker.C:
				s.checkAndRunJobs()
			case <-s.ctx.Done():
				log.Println("Scheduler stopping...")
				return
			}
		}
	}()
}

// checkAndRunJobs queries the database for jobs that need to run
func (s *Scheduler) checkAndRunJobs() {
	var jobs []models.ScheduledJob
	now := time.Now()
	
	// Find all enabled jobs where next_run_at is in the past or null
	err := database.Db.Where("enabled = ? AND (next_run_at IS NULL OR next_run_at <= ?)", true, now).Find(&jobs).Error
	if err != nil {
		log.Printf("Error querying scheduled jobs: %v", err)
		return
	}
	
	for _, job := range jobs {
		// Check if this job is already running (has a pending execution)
		// Use a transaction to prevent race conditions
		var pendingCount int64
		err := database.Db.Transaction(func(tx *gorm.DB) error {
			// Lock the row to prevent concurrent checks
			var count int64
			if err := tx.Model(&models.JobExecution{}).
				Where("scheduled_job_id = ? AND status = ?", job.ID, models.JobStatusPending).
				Count(&count).Error; err != nil {
				return err
			}
			pendingCount = count
			return nil
		})
		
		if err != nil {
			log.Printf("Error checking pending executions for job %s: %v", job.Name, err)
			continue
		}
		
		if pendingCount > 0 {
			log.Printf("Skipping job %s - already running", job.Name)
			continue
		}
		
		// Run each job in a separate goroutine to avoid blocking
		s.wg.Add(1)
		go s.executeJob(job)
	}
}

// executeJob executes a single job and records the execution
func (s *Scheduler) executeJob(job models.ScheduledJob) {
	defer s.wg.Done()
	
	log.Printf("Starting job: %s (ID: %s)", job.Name, job.ID)
	
	// Create job execution record
	execution := &models.JobExecution{
		ScheduledJobID: job.ID,
		StartedAt:      time.Now(),
		Status:         models.JobStatusPending,
	}
	
	if err := database.Db.Create(execution).Error; err != nil {
		log.Printf("Error creating job execution record for job %s: %v", job.Name, err)
		return
	}
	
	// Get the job function
	s.jobsMutex.RLock()
	jobFunc, exists := s.jobs[job.JobHandler]
	s.jobsMutex.RUnlock()
	
	if !exists {
		errMsg := fmt.Sprintf("Job handler not found: %s", job.JobHandler)
		log.Printf("Error executing job %s: %s", job.Name, errMsg)
		s.markJobFailed(execution, errMsg)
		return
	}
	
	// Execute the job with timeout
	done := make(chan error, 1)
	go func() {
		done <- jobFunc()
	}()
	
	// Wait for job completion with timeout (5 minutes default)
	timeout := 5 * time.Minute
	var jobErr error
	timedOut := false
	
	select {
	case jobErr = <-done:
		if jobErr != nil {
			log.Printf("Job %s failed: %v", job.Name, jobErr)
			s.markJobFailed(execution, jobErr.Error())
		} else {
			log.Printf("Job %s completed successfully", job.Name)
			s.markJobSuccess(execution)
		}
	case <-time.After(timeout):
		timedOut = true
		errMsg := "Job execution timeout"
		log.Printf("Job %s timed out after %v", job.Name, timeout)
		errMsgPtr := &errMsg
		if err := execution.MarkCompleted(database.Db, models.JobStatusTimeout, errMsgPtr); err != nil {
			log.Printf("Error marking job execution as timeout: %v", err)
		}
	}
	
	// Update the job's next run time regardless of outcome
	if err := job.UpdateNextRunTime(database.Db); err != nil {
		log.Printf("Error updating next run time for job %s: %v", job.Name, err)
	}
	
	// If timed out, the job goroutine is still running but we return
	if timedOut {
		log.Printf("Job %s goroutine may still be running after timeout", job.Name)
	}
}

// markJobSuccess marks a job execution as successful
func (s *Scheduler) markJobSuccess(execution *models.JobExecution) {
	if err := execution.MarkCompleted(database.Db, models.JobStatusSuccess, nil); err != nil {
		log.Printf("Error marking job execution as success: %v", err)
	}
}

// markJobFailed marks a job execution as failed
func (s *Scheduler) markJobFailed(execution *models.JobExecution, errorMessage string) {
	if err := execution.MarkCompleted(database.Db, models.JobStatusFailed, &errorMessage); err != nil {
		log.Printf("Error marking job execution as failed: %v", err)
	}
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	
	// Stop the ticker
	if s.ticker != nil {
		s.ticker.Stop()
	}
	
	// Cancel the context to signal goroutines to stop
	s.cancel()
	
	// Wait for all jobs to complete (with timeout)
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		log.Println("Scheduler stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Scheduler stopped with timeout (some jobs may still be running)")
	}
}

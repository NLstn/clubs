package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduledJob represents a job that should be executed periodically
type ScheduledJob struct {
	ID              string     `json:"ID" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string     `json:"Name" gorm:"type:varchar(255);unique;not null"`
	Description     string     `json:"Description" gorm:"type:text"`
	JobHandler      string     `json:"JobHandler" gorm:"type:varchar(255);not null"` // function identifier
	Enabled         bool       `json:"Enabled" gorm:"default:true"`
	IntervalMinutes int        `json:"IntervalMinutes" gorm:"not null"`
	LastRunAt       *time.Time `json:"LastRunAt"`
	NextRunAt       *time.Time `json:"NextRunAt" gorm:"index:idx_scheduled_jobs_next_run_at"`
	CreatedAt       time.Time  `json:"CreatedAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"UpdatedAt" gorm:"autoUpdateTime"`
}

// JobExecution represents a single execution of a scheduled job
type JobExecution struct {
	ID             string        `json:"ID" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScheduledJobID string        `json:"ScheduledJobID" gorm:"type:uuid;not null;index:idx_job_executions_scheduled_job_id"`
	ScheduledJob   *ScheduledJob `json:"ScheduledJob,omitempty" gorm:"foreignKey:ScheduledJobID;constraint:OnDelete:CASCADE"`
	StartedAt      time.Time     `json:"StartedAt" gorm:"not null"`
	CompletedAt    *time.Time    `json:"CompletedAt"`
	DurationMs     *int          `json:"DurationMs"`
	Status         string        `json:"Status" gorm:"type:varchar(50);not null"` // 'pending', 'success', 'failed', 'timeout'
	ErrorMessage   *string       `json:"ErrorMessage" gorm:"type:text"`
	CreatedAt      time.Time     `json:"CreatedAt" gorm:"autoCreateTime"`
}

// JobStatus constants
const (
	JobStatusPending = "pending"
	JobStatusSuccess = "success"
	JobStatusFailed  = "failed"
	JobStatusTimeout = "timeout"
)

// BeforeCreate hook for ScheduledJob
func (j *ScheduledJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	// Set NextRunAt if not already set
	if j.NextRunAt == nil {
		now := time.Now()
		j.NextRunAt = &now
	}
	return nil
}

// BeforeCreate hook for JobExecution
func (e *JobExecution) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

// MarkCompleted marks a job execution as completed with success or failure
func (e *JobExecution) MarkCompleted(tx *gorm.DB, status string, errorMessage *string) error {
	now := time.Now()
	e.CompletedAt = &now
	e.Status = status
	e.ErrorMessage = errorMessage
	
	// Calculate duration in milliseconds
	durationMs := int(now.Sub(e.StartedAt).Milliseconds())
	e.DurationMs = &durationMs
	
	return tx.Save(e).Error
}

// UpdateNextRunTime updates the next run time for a scheduled job
func (j *ScheduledJob) UpdateNextRunTime(tx *gorm.DB) error {
	now := time.Now()
	j.LastRunAt = &now
	nextRun := now.Add(time.Duration(j.IntervalMinutes) * time.Minute)
	j.NextRunAt = &nextRun
	return tx.Save(j).Error
}

package core

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/handlers"
	"github.com/stretchr/testify/assert"
)

func TestScheduledJobBeforeCreate(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	t.Run("auto-generates ID", func(t *testing.T) {
		job := &ScheduledJob{
			Name:            "test-job-id",
			JobHandler:      "test_handler",
			IntervalMinutes: 60,
		}
		err := handlers.GetDB().Create(job).Error
		assert.NoError(t, err)
		assert.NotEmpty(t, job.ID)
	})

	t.Run("sets default NextRunAt", func(t *testing.T) {
		job := &ScheduledJob{
			Name:            "test-job-next-run",
			JobHandler:      "test_handler",
			IntervalMinutes: 60,
		}
		err := handlers.GetDB().Create(job).Error
		assert.NoError(t, err)
		assert.NotNil(t, job.NextRunAt)
		assert.True(t, job.NextRunAt.Before(time.Now().Add(1*time.Minute)))
	})
}

func TestScheduledJobUpdateNextRunTime(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	job := &ScheduledJob{
		Name:            "test-job-update",
		JobHandler:      "test_handler",
		IntervalMinutes: 60,
	}
	err := handlers.GetDB().Create(job).Error
	assert.NoError(t, err)

	originalNextRunAt := job.NextRunAt

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Update next run time
	err = job.UpdateNextRunTime(handlers.GetDB())
	assert.NoError(t, err)

	// Verify LastRunAt is set
	assert.NotNil(t, job.LastRunAt)
	assert.True(t, job.LastRunAt.After(time.Now().Add(-1*time.Minute)))

	// Verify NextRunAt is updated (should be LastRunAt + IntervalMinutes)
	assert.NotNil(t, job.NextRunAt)
	if originalNextRunAt != nil {
		assert.True(t, job.NextRunAt.After(*originalNextRunAt))
	}

	// NextRunAt should be approximately 60 minutes from LastRunAt
	expectedNextRun := job.LastRunAt.Add(time.Duration(job.IntervalMinutes) * time.Minute)
	timeDiff := job.NextRunAt.Sub(expectedNextRun).Abs()
	assert.Less(t, timeDiff, 1*time.Second)
}

func TestJobExecutionBeforeCreate(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create a scheduled job first
	job := &ScheduledJob{
		Name:            "test-job-for-execution",
		JobHandler:      "test_handler",
		IntervalMinutes: 60,
	}
	err := handlers.GetDB().Create(job).Error
	assert.NoError(t, err)

	execution := &JobExecution{
		ScheduledJobID: job.ID,
		StartedAt:      time.Now(),
		Status:         JobStatusPending,
	}
	err = handlers.GetDB().Create(execution).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, execution.ID)
}

func TestJobExecutionMarkCompleted(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create a scheduled job first
	job := &ScheduledJob{
		Name:            "test-job-for-completion",
		JobHandler:      "test_handler",
		IntervalMinutes: 60,
	}
	err := handlers.GetDB().Create(job).Error
	assert.NoError(t, err)

	t.Run("mark as success", func(t *testing.T) {
		execution := &JobExecution{
			ScheduledJobID: job.ID,
			StartedAt:      time.Now(),
			Status:         JobStatusPending,
		}
		err := handlers.GetDB().Create(execution).Error
		assert.NoError(t, err)

		// Wait a bit to ensure duration is measurable
		time.Sleep(10 * time.Millisecond)

		err = execution.MarkCompleted(handlers.GetDB(), JobStatusSuccess, nil)
		assert.NoError(t, err)

		assert.Equal(t, JobStatusSuccess, execution.Status)
		assert.NotNil(t, execution.CompletedAt)
		assert.NotNil(t, execution.DurationMs)
		assert.Greater(t, *execution.DurationMs, 0)
		assert.Nil(t, execution.ErrorMessage)
	})

	t.Run("mark as failed with error", func(t *testing.T) {
		execution := &JobExecution{
			ScheduledJobID: job.ID,
			StartedAt:      time.Now(),
			Status:         JobStatusPending,
		}
		err := handlers.GetDB().Create(execution).Error
		assert.NoError(t, err)

		// Wait a bit to ensure duration is measurable
		time.Sleep(10 * time.Millisecond)

		errorMsg := "test error message"
		err = execution.MarkCompleted(handlers.GetDB(), JobStatusFailed, &errorMsg)
		assert.NoError(t, err)

		assert.Equal(t, JobStatusFailed, execution.Status)
		assert.NotNil(t, execution.CompletedAt)
		assert.NotNil(t, execution.DurationMs)
		assert.Greater(t, *execution.DurationMs, 0)
		assert.NotNil(t, execution.ErrorMessage)
		assert.Equal(t, "test error message", *execution.ErrorMessage)
	})
}

func TestJobStatusConstants(t *testing.T) {
	assert.Equal(t, "pending", JobStatusPending)
	assert.Equal(t, "success", JobStatusSuccess)
	assert.Equal(t, "failed", JobStatusFailed)
	assert.Equal(t, "timeout", JobStatusTimeout)
}

func TestScheduledJobEnabledDefault(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	job := &ScheduledJob{
		Name:            "test-job-enabled-default",
		JobHandler:      "test_handler",
		IntervalMinutes: 60,
	}
	err := handlers.GetDB().Create(job).Error
	assert.NoError(t, err)

	// Reload from database to check default value
	var reloaded ScheduledJob
	err = handlers.GetDB().Where("id = ?", job.ID).First(&reloaded).Error
	assert.NoError(t, err)
	assert.True(t, reloaded.Enabled, "Enabled should default to true")
}

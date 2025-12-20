package scheduler_test

import (
	"testing"
	"time"

	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/scheduler"
	"github.com/stretchr/testify/assert"
)

// TestSchedulerIntegration tests the complete scheduler flow including
// OAuth state cleanup job
func TestSchedulerIntegration(t *testing.T) {
	handlers.SetupTestDB(t)
	defer handlers.TeardownTestDB(t)

	// Create scheduler
	s := scheduler.NewScheduler(500 * time.Millisecond)

	// Register the OAuth cleanup job
	s.RegisterJob("cleanup_oauth_states", models.CleanupExpiredOAuthStates)

	// Initialize default jobs
	err := scheduler.InitializeDefaultJobs(database.Db)
	assert.NoError(t, err)

	// Create some OAuth states (expired and valid)
	err = models.CreateOAuthState("valid-state-1", "verifier-1")
	assert.NoError(t, err)

	err = models.CreateOAuthState("valid-state-2", "verifier-2")
	assert.NoError(t, err)

	// Create expired states directly
	expiredState1 := &models.OAuthState{
		State:        "expired-state-1",
		CodeVerifier: "expired-verifier-1",
		ExpiresAt:    time.Now().Add(-1 * time.Hour),
	}
	err = handlers.GetDB().Create(expiredState1).Error
	assert.NoError(t, err)

	expiredState2 := &models.OAuthState{
		State:        "expired-state-2",
		CodeVerifier: "expired-verifier-2",
		ExpiresAt:    time.Now().Add(-2 * time.Hour),
	}
	err = handlers.GetDB().Create(expiredState2).Error
	assert.NoError(t, err)

	// Verify we have 4 states
	var countBefore int64
	handlers.GetDB().Model(&models.OAuthState{}).Count(&countBefore)
	assert.Equal(t, int64(4), countBefore)

	// Start the scheduler
	s.Start()

	// Wait for the job to execute (check every 500ms, job should run within 2 seconds)
	time.Sleep(3 * time.Second)

	// Stop the scheduler
	s.Stop()

	// Verify expired states were cleaned up
	var countAfter int64
	handlers.GetDB().Model(&models.OAuthState{}).Count(&countAfter)
	assert.Equal(t, int64(2), countAfter, "Should have 2 valid states remaining")

	// Verify the job execution was recorded
	var job models.ScheduledJob
	err = database.Db.Where("name = ?", "oauth_state_cleanup").First(&job).Error
	assert.NoError(t, err)

	var executions []models.JobExecution
	err = database.Db.Where("scheduled_job_id = ?", job.ID).Order("started_at DESC").Find(&executions).Error
	assert.NoError(t, err)
	assert.Greater(t, len(executions), 0, "Should have at least one job execution")

	// Verify the execution was successful
	if len(executions) > 0 {
		assert.Equal(t, models.JobStatusSuccess, executions[0].Status)
		assert.NotNil(t, executions[0].CompletedAt)
		assert.NotNil(t, executions[0].DurationMs)
	}

	// Verify LastRunAt and NextRunAt were updated
	err = database.Db.Where("name = ?", "oauth_state_cleanup").First(&job).Error
	assert.NoError(t, err)
	assert.NotNil(t, job.LastRunAt)
	assert.NotNil(t, job.NextRunAt)
	assert.True(t, job.NextRunAt.After(*job.LastRunAt))
}

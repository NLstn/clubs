<div align="center">
  <img src="../assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # Job Scheduler Documentation
  
  **Database-backed background task scheduling**
</div>

---

# Job Scheduler Documentation

## Overview

The Clubs application includes a database-backed job scheduling mechanism for running periodic background tasks. This system is designed to be simple, maintainable, and extensible while providing full visibility into job execution history.

## Architecture

The scheduler consists of three main components:

### 1. Database Models

#### ScheduledJob
Stores configuration for jobs that should run periodically:
- **ID**: Unique identifier (UUID)
- **Name**: Unique name for the job
- **Description**: Human-readable description
- **JobHandler**: String identifier for the job function
- **Enabled**: Boolean flag to enable/disable job execution
- **IntervalMinutes**: How often the job should run (in minutes)
- **LastRunAt**: Timestamp of last execution
- **NextRunAt**: Timestamp of next scheduled execution
- **CreatedAt/UpdatedAt**: Audit timestamps

#### JobExecution
Tracks each execution of a scheduled job:
- **ID**: Unique identifier (UUID)
- **ScheduledJobID**: Foreign key to ScheduledJob
- **StartedAt**: When execution began
- **CompletedAt**: When execution finished
- **DurationMs**: Execution duration in milliseconds
- **Status**: One of: `pending`, `success`, `failed`, `timeout`
- **ErrorMessage**: Error details (if failed)
- **CreatedAt**: Creation timestamp

#### OAuthState
Stores OAuth state parameters for CSRF protection:
- **ID**: Unique identifier (UUID)
- **State**: OAuth state parameter
- **CodeVerifier**: PKCE code verifier
- **ExpiresAt**: When the state expires
- **CreatedAt**: Creation timestamp

### 2. Scheduler Package

The `scheduler` package (`Backend/scheduler/`) provides the core scheduling logic:

- **Job Registry**: Register job handler functions by name
- **Periodic Execution**: Uses `time.Ticker` to check for jobs that need to run
- **Concurrent Execution**: Jobs run in separate goroutines
- **Error Handling**: Captures and logs job failures
- **Graceful Shutdown**: Allows jobs to complete on application shutdown
- **Timeout Protection**: Jobs timeout after 5 minutes (configurable)

### 3. Job Handlers

Job handlers are simple Go functions with the signature:
```go
func() error
```

They are registered with the scheduler at application startup.

## Current Jobs

### OAuth State Cleanup
- **Name**: `oauth_state_cleanup`
- **Handler**: `cleanup_oauth_states`
- **Interval**: 60 minutes (1 hour)
- **Description**: Removes expired OAuth state records from the database
- **Function**: `models.CleanupExpiredOAuthStates()`

This job prevents the `oauth_states` table from accumulating expired records indefinitely.

## How to Add a New Job

### Step 1: Create the Job Handler Function

Define your job logic as a function that returns an error:

```go
// Backend/models/your_model.go
func CleanupExpiredSessions() error {
    result := database.Db.Where("expires_at < ?", time.Now()).Delete(&Session{})
    if result.Error != nil {
        return fmt.Errorf("failed to cleanup sessions: %w", result.Error)
    }
    log.Printf("Cleaned up %d expired sessions", result.RowsAffected)
    return nil
}
```

### Step 2: Register the Job Handler

In `Backend/main.go`, register your job handler with the scheduler:

```go
// Register job handlers
jobScheduler.RegisterJob("cleanup_oauth_states", models.CleanupExpiredOAuthStates)
jobScheduler.RegisterJob("cleanup_sessions", models.CleanupExpiredSessions)  // NEW
```

### Step 3: Add a Default Job Configuration

In `Backend/scheduler/scheduler.go`, add your job to `InitializeDefaultJobs()`:

```go
func InitializeDefaultJobs(db *gorm.DB) error {
    // ... existing oauth cleanup job ...

    // Session cleanup job
    var sessionCleanupJob models.ScheduledJob
    err = db.Where("name = ?", "session_cleanup").First(&sessionCleanupJob).Error
    if err == gorm.ErrRecordNotFound {
        sessionCleanupJob = models.ScheduledJob{
            Name:            "session_cleanup",
            Description:     "Removes expired user sessions from the database",
            JobHandler:      "cleanup_sessions",
            Enabled:         true,
            IntervalMinutes: 1440, // Run daily (24 hours)
        }
        if err := db.Create(&sessionCleanupJob).Error; err != nil {
            return fmt.Errorf("failed to create session_cleanup job: %w", err)
        }
        log.Println("Created default job: session_cleanup")
    } else if err != nil {
        return fmt.Errorf("failed to query session_cleanup job: %w", err)
    }

    return nil
}
```

### Step 4: Add Tests

Create tests for your job in `Backend/models/your_model_test.go`:

```go
func TestCleanupExpiredSessions(t *testing.T) {
    handlers.SetupTestDB(t)
    defer handlers.TeardownTestDB(t)

    // Create test data...
    
    // Run cleanup
    err := models.CleanupExpiredSessions()
    assert.NoError(t, err)

    // Verify results...
}
```

## Configuration

### Environment Variables

The scheduler uses the same database configuration as the rest of the application:
- `DATABASE_URL`: Database host
- `DATABASE_PORT`: Database port
- `DATABASE_USER`: Database user
- `DATABASE_USER_PASSWORD`: Database password
- `DATABASE_NAME`: Database name

### Adjusting Job Intervals

Job intervals can be adjusted by updating the `scheduled_jobs` table directly:

```sql
UPDATE scheduled_jobs 
SET interval_minutes = 30 
WHERE name = 'oauth_state_cleanup';
```

Changes take effect on the next scheduler tick (default: 1 minute).

### Enabling/Disabling Jobs

Jobs can be enabled or disabled without code changes:

```sql
-- Disable a job
UPDATE scheduled_jobs 
SET enabled = false 
WHERE name = 'oauth_state_cleanup';

-- Enable a job
UPDATE scheduled_jobs 
SET enabled = true 
WHERE name = 'oauth_state_cleanup';
```

## Monitoring

### View Job Configuration

```sql
SELECT id, name, enabled, interval_minutes, last_run_at, next_run_at
FROM scheduled_jobs
ORDER BY name;
```

### View Recent Job Executions

```sql
SELECT 
    sj.name,
    je.started_at,
    je.completed_at,
    je.duration_ms,
    je.status,
    je.error_message
FROM job_executions je
JOIN scheduled_jobs sj ON sj.id = je.scheduled_job_id
ORDER BY je.started_at DESC
LIMIT 50;
```

### Check for Failed Jobs

```sql
SELECT 
    sj.name,
    je.started_at,
    je.error_message
FROM job_executions je
JOIN scheduled_jobs sj ON sj.id = je.scheduled_job_id
WHERE je.status = 'failed'
ORDER BY je.started_at DESC;
```

### Job Success Rate

**Note:** This query uses PostgreSQL-specific syntax (INTERVAL). For SQLite or other databases, adjust the date arithmetic accordingly.

```sql
-- PostgreSQL
SELECT 
    sj.name,
    COUNT(*) as total_executions,
    SUM(CASE WHEN je.status = 'success' THEN 1 ELSE 0 END) as successful,
    SUM(CASE WHEN je.status = 'failed' THEN 1 ELSE 0 END) as failed,
    AVG(je.duration_ms) as avg_duration_ms
FROM job_executions je
JOIN scheduled_jobs sj ON sj.id = je.scheduled_job_id
WHERE je.completed_at > NOW() - INTERVAL '7 days'
GROUP BY sj.name;

-- SQLite alternative
SELECT 
    sj.name,
    COUNT(*) as total_executions,
    SUM(CASE WHEN je.status = 'success' THEN 1 ELSE 0 END) as successful,
    SUM(CASE WHEN je.status = 'failed' THEN 1 ELSE 0 END) as failed,
    AVG(je.duration_ms) as avg_duration_ms
FROM job_executions je
JOIN scheduled_jobs sj ON sj.id = je.scheduled_job_id
WHERE je.completed_at > datetime('now', '-7 days')
GROUP BY sj.name;
```

## Application Startup and Shutdown

### Startup
1. Database is initialized and migrated
2. Scheduler is created with 1-minute tick interval
3. Job handlers are registered
4. Default jobs are created in the database (if they don't exist)
5. Scheduler starts running in a background goroutine
6. HTTP server starts

### Shutdown
1. Application receives SIGINT or SIGTERM signal
2. HTTP server stops accepting new requests
3. Scheduler's `Stop()` method is called
4. Scheduler stops checking for new jobs
5. Waits up to 30 seconds for running jobs to complete
6. Application exits

## Error Handling

### Job Execution Errors
- Errors returned by job handlers are logged
- Job execution record is marked as `failed` with error message
- Job's `NextRunAt` is still updated (job will retry on next interval)

### Timeout
- Jobs that run longer than 5 minutes are marked as `timeout`
- The job goroutine may still be running (no force-kill)

### Database Errors
- Errors querying scheduled jobs are logged but don't stop the scheduler
- Errors creating execution records are logged but don't prevent job execution

## Best Practices

1. **Keep Jobs Idempotent**: Jobs should be safe to run multiple times with the same result
2. **Handle Errors Gracefully**: Return descriptive error messages for debugging
3. **Log Progress**: Use `log.Printf()` to log important steps
4. **Test Thoroughly**: Write comprehensive tests for job logic
5. **Monitor Regularly**: Check job execution history for failures
6. **Set Appropriate Intervals**: Balance between frequency and database load
7. **Add Indexes**: If querying large tables, ensure proper indexes exist

## Limitations and Future Enhancements

### Current Limitations
- Single instance only (no distributed locking)
- Fixed 5-minute timeout for all jobs
- No job prioritization
- No job dependencies
- Manual configuration changes require database access

### Possible Future Enhancements
- Add distributed locking for multi-instance deployments
- Configurable per-job timeouts
- Job priority system
- Job dependency chains
- Admin UI for job management
- Metrics and alerting integration
- Retry policies with exponential backoff
- Cron expression support for complex schedules

## Troubleshooting

### Job Not Running
1. Check if job is enabled: `SELECT enabled FROM scheduled_jobs WHERE name = 'job_name'`
2. Check next run time: `SELECT next_run_at FROM scheduled_jobs WHERE name = 'job_name'`
3. Check scheduler logs for errors
4. Verify job handler is registered in `main.go`

### Job Failing Consistently
1. Check error messages: `SELECT error_message FROM job_executions WHERE scheduled_job_id = 'job_id' AND status = 'failed' ORDER BY started_at DESC LIMIT 5`
2. Test job function directly in a test
3. Check database connectivity and permissions
4. Verify job logic handles edge cases

### Job Timing Out
1. Check job execution duration: `SELECT duration_ms FROM job_executions WHERE scheduled_job_id = 'job_id' ORDER BY started_at DESC LIMIT 5`
2. Optimize job logic or break into smaller jobs
3. Consider increasing timeout in `scheduler.go` (requires code change)

## Security Considerations

- Job execution records may contain sensitive error messages
- Ensure database access is restricted
- OAuth state cleanup protects against state fixation attacks
- Job handlers should validate all inputs
- Jobs with external API calls should use secure credentials

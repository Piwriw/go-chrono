package chrono

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Triggered before the job starts
var defaultBeforeJobRuns = func(jobID uuid.UUID, jobName string) {
	slog.Info("Job is about to start", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
}

// Triggered before the job starts (can be skipped, returns error)
var defaultBeforeJobRunsSkipIfBeforeFuncErrors = func(jobID uuid.UUID, jobName string) error {
	slog.Info("Job is about to start (skippable)", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
	return nil
}

// Triggered after the job completes successfully
var defaultAfterJobRuns = func(jobID uuid.UUID, jobName string) {
	slog.Info("Job completed successfully", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
}

// Triggered after the job completes with error
var defaultAfterJobRunsWithError = func(jobID uuid.UUID, jobName string, err error) {
	slog.Error("Job completed with error", "jobID", jobID, "jobName", jobName, "error", err, "execTime", time.Now().Format(time.DateTime))
}

// Triggered after the job panics
var defaultAfterJobRunsWithPanic = func(jobID uuid.UUID, jobName string, recoverData any) {
	slog.Error("Job panicked during execution", "jobID", jobID, "jobName", jobName, "panic", recoverData, "execTime", time.Now().Format(time.DateTime))
}

// Triggered when job lock fails
var defaultAfterLockError = func(jobID uuid.UUID, jobName string, err error) {
	slog.Error("Job lock failed", "jobID", jobID, "jobName", jobName, "error", err, "execTime", time.Now().Format(time.DateTime))
}

// EmptyWatchFunc Empty monitor event handler
var EmptyWatchFunc = func(event JobWatchInterface) {}

// Empty after job runs with error hook
var EmptyAfterJobRunsWithError = defaultAfterJobRunsWithError

// Empty after job runs with panic hook
var EmptyAfterJobRunsWithPanic = defaultAfterJobRunsWithPanic

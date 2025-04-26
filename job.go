package chrono

import (
	"log/slog"

	"github.com/google/uuid"
)

var (
	defaultBeforeJobRuns = func(jobID uuid.UUID, jobName string) {
		slog.Info("chrono:BeforeJobRuns", "jobID", jobID, "jobName", jobName)
	}
	defaultBeforeJobRunsSkipIfBeforeFuncErrors = func(jobID uuid.UUID, jobName string) error {
		slog.Info("chrono:BeforeJobRunsSkipIfBeforeFuncErrors", "jobID", jobID, "jobName", jobName)
		return nil
	}
	defaultAfterJobRuns = func(jobID uuid.UUID, jobName string) {
		slog.Info("chrono:BeforeJobRunsAfterFuncErrors", "jobID", jobID, "jobName", jobName)
	}
	defaultAfterJobRunsWithError = func(jobID uuid.UUID, jobName string, err error) {
		slog.Error("chrono:AfterJobRunsWithError", "jobID", jobID, "jobName", jobName, "err", err)
	}
	defaultAfterJobRunsWithPanic = func(jobID uuid.UUID, jobName string, recoverData any) {
		slog.Error("chrono:AfterJobRunsWithPanic", "jobID", jobID, "jobName", jobName, "recoverData", recoverData)
	}
	defaultAfterLockError = func(jobID uuid.UUID, jobName string, err error) {
		slog.Error("chrono:AfterLockError", "jobID", jobID, "jobName", jobName, "err", err)
	}
	EmptyWatchFunc             = func(event MonitorJobSpec) {}
	EmptyAfterJobRunsWithError = func(jobID uuid.UUID, jobName string, err error) {
		slog.Error("chrono:AfterJobRunsWithError", "jobID", jobID, "jobName", jobName, "err", err)
	}
	EmptyAfterJobRunsWithPanic = func(jobID uuid.UUID, jobName string, recoverData any) {
		slog.Error("¨chrono:AfterJobRunsWithPanic", "jobID", jobID, "jobName", jobName, "recoverData", recoverData)
	}
)

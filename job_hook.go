package chrono

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/retry"
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

// jobOptions holds options for job execution.
// jobOptions 保存任务执行的选项。
type jobOptions struct {
	// tags are the job tags.
	// tags 是任务标签。
	tags []string
	// retryConfig is the retry configuration.
	// retryConfig 是重试配置。
	retryConfig *retry.RetryConfig
	// retryEnabled indicates if retry is enabled for this job.
	// retryEnabled 表示是否为此任务启用重试。
	retryEnabled bool
}

// wrapTaskWithRetry wraps a task function with retry logic.
// wrapTaskWithRetry 用重试逻辑包装任务函数。
//
// Parameters:
//
//	taskFunc - The task function to wrap / 要包装的任务函数
//	jobID   - The job ID / 任务 ID
//	jobName - The job name / 任务名称
//	opts    - The job options / 任务选项
//
// Returns:
//
//	func() error - The wrapped task function / 包装后的任务函数
func wrapTaskWithRetry(taskFunc func() error, jobID uuid.UUID, jobName string, opts *jobOptions) func() error {
	// If no retry configuration, return the original task function
	// 如果没有重试配置，返回原始任务函数
	if opts == nil || opts.retryConfig == nil || !opts.retryEnabled {
		return taskFunc
	}

	// Wrap the task with retry executor
	// 用重试执行器包装任务
	return func() error {
		executor := retry.NewRetryExecutor(opts.retryConfig, jobID, jobName)
		return executor.ExecuteWithRetry(context.TODO(), taskFunc)
	}
}

// newJobOptions creates a new jobOptions instance.
// newJobOptions 创建一个新的 jobOptions 实例。
//
// Returns:
//
//	*jobOptions - The new job options / 新的任务选项
func newJobOptions() *jobOptions {
	return &jobOptions{
		tags:         make([]string, 0),
		retryEnabled: false,
	}
}

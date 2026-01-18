package common

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/retry"
)

// JobWatchInterface defines the interface for jobs event watching.
// JobWatchInterface 定义了任务事件监听的接口。
type JobWatchInterface interface {
	// GetJobID gets the jobs ID.
	// 获取任务 ID。
	//
	// Returns:
	//	string - The jobs ID / 任务 ID
	GetJobID() string
	// GetJobName gets the jobs name.
	// 获取任务名称。
	//
	// Returns:
	//	string - The jobs name / 任务名称
	GetJobName() string
	// GetStartTime gets the jobs start time.
	// 获取任务开始时间。
	//
	// Returns:
	//	*time.Time - The start time / 开始时间
	GetStartTime() *time.Time
	// GetEndTime gets the jobs end time.
	// 获取任务结束时间。
	//
	// Returns:
	//	*time.Time - The end time / 结束时间
	GetEndTime() *time.Time
	// GetStatus gets the jobs status.
	// 获取任务状态。
	//
	// Returns:
	//	int - The jobs status / 任务状态
	GetStatus() int
	// GetTags gets the jobs tags.
	// 获取任务标签。
	//
	// Returns:
	//	[]string - The jobs tags / 任务标签
	GetTags() []string
	// Error gets the jobs error.
	// 获取任务错误。
	//
	// Returns:
	//	error - The jobs error / 任务错误
	Error() error
	// GetCurrentEvent gets the current jobs event.
	// 获取当前任务事件。
	//
	// Returns:
	//	*JobEvent - The current jobs event / 当前任务事件
	GetCurrentEvent() *JobEvent
}

// JobEvent represents a jobs execution event.
// JobEvent 表示一个任务执行事件。
type JobEvent struct {
	// EventID is the event ID.
	// EventID 是事件 ID。
	EventID string
	// StartTime is the start time.
	// StartTime 是开始时间。
	StartTime *time.Time
	// EndTime is the end time.
	// EndTime 是结束时间。
	EndTime *time.Time
	// Status is the jobs status.
	// Status 是任务状态。
	Status int
	// Err is the jobs error.
	// Err 是任务错误。
	Err error
	// RetryCount is the number of retry attempts.
	// RetryCount 重试次数。
	RetryCount int `json:"retry_count,omitempty"`
	// IsRetry indicates if this is a retry event.
	// IsRetry 表示是否为重试事件。
	IsRetry bool `json:"is_retry,omitempty"`
	// OriginalEventID is the original event ID for retries.
	// OriginalEventID 原始事件 ID（用于重试）。
	OriginalEventID string `json:"original_event_id,omitempty"`
}

// DefaultBeforeJobRuns is triggered before the jobs starts.
// DefaultBeforeJobRuns 在任务开始前触发。
var DefaultBeforeJobRuns = func(jobID uuid.UUID, jobName string) {
	slog.Info("Job is about to start", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
}

// DefaultBeforeJobRunsSkipIfBeforeFuncErrors is triggered before the jobs starts (can be skipped, returns error).
// DefaultBeforeJobRunsSkipIfBeforeFuncErrors 在任务开始前触发（可跳过，返回错误）。
var DefaultBeforeJobRunsSkipIfBeforeFuncErrors = func(jobID uuid.UUID, jobName string) error {
	slog.Info("Job is about to start (skippable)", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
	return nil
}

// DefaultAfterJobRuns is triggered after the jobs completes successfully.
// DefaultAfterJobRuns 在任务成功完成后触发。
var DefaultAfterJobRuns = func(jobID uuid.UUID, jobName string) {
	slog.Info("Job completed successfully", "jobID", jobID, "jobName", jobName, "execTime", time.Now().Format(time.DateTime))
}

// DefaultAfterJobRunsWithError is triggered after the jobs completes with error.
// DefaultAfterJobRunsWithError 在任务出错完成后触发。
var DefaultAfterJobRunsWithError = func(jobID uuid.UUID, jobName string, err error) {
	slog.Error("Job completed with error", "jobID", jobID, "jobName", jobName, "error", err, "execTime", time.Now().Format(time.DateTime))
}

// DefaultAfterJobRunsWithPanic is triggered after the jobs panics.
// DefaultAfterJobRunsWithPanic 在任务 panic 后触发。
var DefaultAfterJobRunsWithPanic = func(jobID uuid.UUID, jobName string, recoverData any) {
	slog.Error("Job panicked during execution", "jobID", jobID, "jobName", jobName, "panic", recoverData, "execTime", time.Now().Format(time.DateTime))
}

// DefaultAfterLockError is triggered when jobs lock fails.
// DefaultAfterLockError 在任务加锁失败时触发。
var DefaultAfterLockError = func(jobID uuid.UUID, jobName string, err error) {
	slog.Error("Job lock failed", "jobID", jobID, "jobName", jobName, "error", err, "execTime", time.Now().Format(time.DateTime))
}

// EmptyWatchFunc is an empty monitor event handler.
// EmptyWatchFunc 是一个空的监控事件处理器。
var EmptyWatchFunc = func(event JobWatchInterface) {}

// EmptyAfterJobRunsWithError is an empty after jobs runs with error hook.
// EmptyAfterJobRunsWithError 是一个空的任务运行后错误钩子。
var EmptyAfterJobRunsWithError = DefaultAfterJobRunsWithError

// EmptyAfterJobRunsWithPanic is an empty after jobs runs with panic hook.
// EmptyAfterJobRunsWithPanic 是一个空的任务运行后 panic 钩子。
var EmptyAfterJobRunsWithPanic = DefaultAfterJobRunsWithPanic

// JobOptions holds options for jobs execution.
// JobOptions 保存任务执行的选项。
type JobOptions struct {
	// tags are the jobs tags.
	// tags 是任务标签。
	tags []string
	// retryConfig is the retry configuration.
	// retryConfig 是重试配置。
	retryConfig *retry.RetryConfig
	// retryEnabled indicates if retry is enabled for this jobs.
	// retryEnabled 表示是否为此任务启用重试。
	retryEnabled bool
}

// NewJobOptions creates a new JobOptions instance.
// NewJobOptions 创建一个新的 JobOptions 实例。
//
// Returns:
//
//	*JobOptions - The new jobs options / 新的任务选项
func NewJobOptions() *JobOptions {
	return &JobOptions{
		tags:         make([]string, 0),
		retryEnabled: false,
	}
}

// SetRetryConfig sets the retry configuration for the jobs options.
// SetRetryConfig 设置任务选项的重试配置。
//
// Parameters:
//
//	config - The retry configuration / 重试配置
func (o *JobOptions) SetRetryConfig(config *retry.RetryConfig) {
	o.retryConfig = config
	o.retryEnabled = true
}

// GetRetryConfig gets the retry configuration from the jobs options.
// GetRetryConfig 获取任务选项的重试配置。
//
// Returns:
//
//	*retry.RetryConfig - The retry configuration / 重试配置
func (o *JobOptions) GetRetryConfig() *retry.RetryConfig {
	return o.retryConfig
}

// IsRetryEnabled checks if retry is enabled for the jobs options.
// IsRetryEnabled 检查任务选项是否启用了重试。
//
// Returns:
//
//	bool - True if retry is enabled, false otherwise / 如果启用重试返回 true，否则返回 false
func (o *JobOptions) IsRetryEnabled() bool {
	return o.retryEnabled
}

// WrapTaskWithRetry wraps a task function with retry logic.
// WrapTaskWithRetry 用重试逻辑包装任务函数。
//
// Parameters:
//
//	taskFunc - The task function to wrap / 要包装的任务函数
//	jobID   - The jobs ID / 任务 ID
//	jobName - The jobs name / 任务名称
//	opts    - The jobs options / 任务选项
//
// Returns:
//
//	func() error - The wrapped task function / 包装后的任务函数
func WrapTaskWithRetry(taskFunc func() error, jobID uuid.UUID, jobName string, opts *JobOptions) func() error {
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

// Job defines the interface for a jobs.
// Job 定义了任务的接口。
type Job interface {
	// ID returns the jobs ID.
	// ID 返回任务 ID。
	ID() uuid.UUID
	// Name returns the jobs name.
	// Name 返回任务名称。
	Name() string
}

// SchedulerInterface defines the interface for a scheduler.
// SchedulerInterface 定义了调度器的接口。
type SchedulerInterface interface {
	// GetJobs returns all jobs.
	// GetJobs 返回所有任务。
	//
	// Returns:
	//	[]Job - The list of jobs / 任务列表
	GetJobs() ([]Job, error)
	// GetJobLastAndNextByID gets the last and next run time for a jobs.
	// GetJobLastAndNextByID 获取任务的最后和下次运行时间。
	//
	// Parameters:
	//
	//	jobID - The jobs ID / 任务 ID
	//
	// Returns:
	//	*time.Time - Last run time / 最后运行时间
	//	*time.Time - Next run time / 下次运行时间
	//	error - Error if the operation fails / 操作失败时返回错误
	GetJobLastAndNextByID(jobID string) (*time.Time, *time.Time, error)
	// Enable checks if a specific option is enabled.
	// Enable 用于查询某个选项是否启用。
	//
	// Parameters:
	//
	//	optionName - The option name / 选项名称
	//
	// Returns:
	//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
	Enable(optionName string) bool
	// GetAlias gets the alias for a jobs.
	// GetAlias 获取任务的别名。
	//
	// Parameters:
	//
	//	jobID - The jobs ID / 任务 ID
	//
	// Returns:
	//	string - The alias / 别名
	//	error - Error if the operation fails / 操作失败时返回错误
	GetAlias(jobID string) (string, error)

	// AddCronJob adds a cron jobs.
	// AddCronJob 添加一个 cron 任务。
	AddCronJob(job any) (gocron.Job, error)
	// AddCronJobs adds multiple cron jobs.
	// AddCronJobs 添加多个 cron 任务。
	AddCronJobs(jobs ...any) ([]gocron.Job, error)
	// AddOnceJob adds a once jobs.
	// AddOnceJob 添加一个单次任务。
	AddOnceJob(job any) (gocron.Job, error)
	// AddOnceJobs adds multiple once jobs.
	// AddOnceJobs 添加多个单次任务。
	AddOnceJobs(jobs ...any) ([]gocron.Job, error)
	// AddIntervalJob adds an interval jobs.
	// AddIntervalJob 添加一个间隔任务。
	AddIntervalJob(job any) (gocron.Job, error)
	// AddIntervalJobs adds multiple interval jobs.
	// AddIntervalJobs 添加多个间隔任务。
	AddIntervalJobs(jobs ...any) ([]gocron.Job, error)
	// AddDailyJob adds a daily jobs.
	// AddDailyJob 添加一个每日任务。
	AddDailyJob(job any) (gocron.Job, error)
	// AddDailyJobs adds multiple daily jobs.
	// AddDailyJobs 添加多个每日任务。
	AddDailyJobs(jobs ...any) ([]gocron.Job, error)
	// AddWeeklyJob adds a weekly jobs.
	// AddWeeklyJob 添加一个每周任务。
	AddWeeklyJob(job any) (gocron.Job, error)
	// AddWeeklyJobs adds multiple weekly jobs.
	// AddWeeklyJobs 添加多个每周任务。
	AddWeeklyJobs(jobs ...any) ([]gocron.Job, error)
	// AddMonthlyJob adds a monthly jobs.
	// AddMonthlyJob 添加一个每月任务。
	AddMonthlyJob(job any) (gocron.Job, error)
	// AddMonthlyJobs adds multiple monthly jobs.
	// AddMonthlyJobs 添加多个每月任务。
	AddMonthlyJobs(jobs ...any) ([]gocron.Job, error)
	// RemoveJob removes a jobs by ID.
	// RemoveJob 通过 ID 删除任务。
	RemoveJob(jobID string) error
	// RemoveJobByAlias removes a jobs by alias.
	// RemoveJobByAlias 通过别名删除任务。
	RemoveJobByAlias(alias string) error
	// RemoveJobByName removes a jobs by name.
	// RemoveJobByName 通过名称删除任务。
	RemoveJobByName(name string) error
	// GetJobByID gets a jobs by ID.
	// GetJobByID 通过 ID 获取任务。
	GetJobByID(jobID string) (gocron.Job, error)
	// GetJobByAlias gets a jobs by alias.
	// GetJobByAlias 通过别名获取任务。
	GetJobByAlias(alias string) (gocron.Job, error)
	// GetJobByName gets a jobs by name.
	// GetJobByName 通过名称获取任务。
	GetJobByName(name string) (gocron.Job, error)
}

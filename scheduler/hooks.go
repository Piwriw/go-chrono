package scheduler

import (
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/monitor"
	"github.com/piwriw/go-chrono/pkg/format"
)

// Re-export types and functions from pkg for backward compatibility
// 从 pkg 包重新导出类型和函数以保持向后兼容

// Small exported aliases for testing convenience
// 为测试方便提供的小写导出别名

// defaultBeforeJobRuns is triggered before the jobs starts.
// defaultBeforeJobRuns 在任务开始前触发。
var defaultBeforeJobRuns = common.DefaultBeforeJobRuns

// defaultBeforeJobRunsSkipIfBeforeFuncErrors is triggered before the jobs starts (can be skipped, returns error).
// defaultBeforeJobRunsSkipIfBeforeFuncErrors 在任务开始前触发（可跳过，返回错误）。
var defaultBeforeJobRunsSkipIfBeforeFuncErrors = common.DefaultBeforeJobRunsSkipIfBeforeFuncErrors

// defaultAfterJobRuns is triggered after the jobs completes successfully.
// defaultAfterJobRuns 在任务成功完成后触发。
var defaultAfterJobRuns = common.DefaultAfterJobRuns

// defaultAfterJobRunsWithError is triggered after the jobs completes with error.
// defaultAfterJobRunsWithError 在任务出错完成后触发。
var defaultAfterJobRunsWithError = common.DefaultAfterJobRunsWithError

// defaultAfterJobRunsWithPanic is triggered after the jobs panics.
// defaultAfterJobRunsWithPanic 在任务 panic 后触发。
var defaultAfterJobRunsWithPanic = common.DefaultAfterJobRunsWithPanic

// defaultAfterLockError is triggered when jobs lock fails.
// defaultAfterLockError 在任务加锁失败时触发。
var defaultAfterLockError = common.DefaultAfterLockError

// DefaultBeforeJobRuns is triggered before the jobs starts.
// DefaultBeforeJobRuns 在任务开始前触发。
var DefaultBeforeJobRuns = common.DefaultBeforeJobRuns

// DefaultBeforeJobRunsSkipIfBeforeFuncErrors is triggered before the jobs starts (can be skipped, returns error).
// DefaultBeforeJobRunsSkipIfBeforeFuncErrors 在任务开始前触发（可跳过，返回错误）。
var DefaultBeforeJobRunsSkipIfBeforeFuncErrors = common.DefaultBeforeJobRunsSkipIfBeforeFuncErrors

// DefaultAfterJobRuns is triggered after the jobs completes successfully.
// DefaultAfterJobRuns 在任务成功完成后触发。
var DefaultAfterJobRuns = common.DefaultAfterJobRuns

// DefaultAfterJobRunsWithError is triggered after the jobs completes with error.
// DefaultAfterJobRunsWithError 在任务出错完成后触发。
var DefaultAfterJobRunsWithError = common.DefaultAfterJobRunsWithError

// DefaultAfterJobRunsWithPanic is triggered after the jobs panics.
// DefaultAfterJobRunsWithPanic 在任务 panic 后触发。
var DefaultAfterJobRunsWithPanic = common.DefaultAfterJobRunsWithPanic

// DefaultAfterLockError is triggered when jobs lock fails.
// DefaultAfterLockError 在任务加锁失败时触发。
var DefaultAfterLockError = common.DefaultAfterLockError

// EmptyWatchFunc is an empty monitor event handler.
// EmptyWatchFunc 是一个空的监控事件处理器。
var EmptyWatchFunc = common.EmptyWatchFunc

// EmptyWatchFuncMonitor is an empty watch function that accepts monitor.JobWatchInterface.
// EmptyWatchFuncMonitor 是接受 monitor.JobWatchInterface 的空监听函数。
var EmptyWatchFuncMonitor = func(event monitor.JobWatchInterface) {}

// EmptyAfterJobRunsWithError is an empty after jobs runs with error hook.
// EmptyAfterJobRunsWithError 是一个空的任务运行后错误钩子。
var EmptyAfterJobRunsWithError = common.EmptyAfterJobRunsWithError

// EmptyAfterJobRunsWithPanic is an empty after jobs runs with panic hook.
// EmptyAfterJobRunsWithPanic 是一个空的任务运行后 panic 钩子。
var EmptyAfterJobRunsWithPanic = common.EmptyAfterJobRunsWithPanic

// JobWatchInterface is an alias for pkg.JobWatchInterface.
// JobWatchInterface 是 pkg.JobWatchInterface 的别名。
type JobWatchInterface = common.JobWatchInterface

// JobEvent is an alias for pkg.JobEvent.
// JobEvent 是 pkg.JobEvent 的别名。
type JobEvent = common.JobEvent

// JobOptions is an alias for pkg.JobOptions.
// JobOptions 是 pkg.JobOptions 的别名。
type JobOptions = common.JobOptions

// NewJobOptions creates a new JobOptions instance.
// NewJobOptions 创建一个新的 JobOptions 实例。
//
// Returns:
//
//	*JobOptions - The new jobs options / 新的任务选项
func NewJobOptions() *JobOptions {
	return common.NewJobOptions()
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
func WrapTaskWithRetry(taskFunc func() error, jobID string, jobName string, opts *JobOptions) func() error {
	// Convert string jobID to uuid.UUID
	// 将字符串 jobID 转换为 uuid.UUID
	parsedID, err := format.ParseJobID(jobID)
	if err != nil {
		return taskFunc // Return original function if ID is invalid
	}
	return common.WrapTaskWithRetry(taskFunc, parsedID, jobName, opts)
}

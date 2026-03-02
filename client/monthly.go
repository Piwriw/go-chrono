package client

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/retry"
)

type MonthJobClientInterface interface {
	// AtTime sets the specific days and time for the monthly jobs to run.
	// 设置月度任务运行的特定日期和时间。
	//
	// Parameters:
	//	days   - The days of the month when the jobs should run / 任务应该运行的月份日期
	//	hour   - The hour of the day (0-23) / 一天中的小时（0-23）
	//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
	//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AtTime(days []int, hour, minute, second int) MonthJobClientInterface

	// Alias sets an alias for the jobs.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) MonthJobClientInterface

	// JobID sets the unique identifier for the jobs.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The jobs ID string / 任务ID字符串
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) MonthJobClientInterface

	// Name sets the name for the jobs.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The jobs name string / 任务名称字符串
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) MonthJobClientInterface

	// Tags adds tags to the jobs for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tags(tags ...string) MonthJobClientInterface

	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) MonthJobClientInterface

	// Watch sets a watcher function to monitor jobs events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for jobs events / 任务事件的监视函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event common.JobWatchInterface)) MonthJobClientInterface

	// DefaultHooks enables default event hooks for the jobs.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() MonthJobClientInterface

	// BeforeJobRuns sets a callback function to be executed before the jobs runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface

	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) MonthJobClientInterface

	// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface

	// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface

	// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) MonthJobClientInterface

	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface

	// WithRetry sets the retry configuration for the jobs.
	// 设置任务的重试配置。
	//
	// Parameters:
	//	maxRetries - Maximum number of retries / 最大重试次数
	//	policy     - Retry policy / 重试策略
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetry(maxRetries int, policy retry.RetryPolicy) MonthJobClientInterface

	// WithRetryConfig sets the complete retry configuration.
	// 设置完整的重试配置。
	//
	// Parameters:
	//	config - Retry configuration / 重试配置
	// Returns:
	//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetryConfig(config *retry.RetryConfig) MonthJobClientInterface

	// Add adds the configured monthly jobs to the scheduler.
	// 将配置的月度任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added jobs / 添加的任务
	//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)

	// BatchAdd adds multiple monthly jobs to the scheduler.
	// 批量添加多个月度任务到调度器。
	//
	// Parameters:
	//	monthJobs - Variable number of monthly jobs to add / 可变数量的要添加的月度任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(monthJobs ...*jobs.MonthJob) ([]gocron.Job, error)

	// Remove removes the monthly jobs from the scheduler.
	// 从调度器中删除月度任务。
	//
	// Returns:
	//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
	Remove() error

	// Get retrieves the monthly jobs from the scheduler.
	// 从调度器中获取月度任务。
	//
	// Returns:
	//	gocron.Job - The retrieved jobs / 获取的任务
	//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// MonthJobClient is a client for managing monthly jobs with a scheduler.
// MonthJobClient 是一个用于通过调度器管理月度任务的客户端。
type MonthJobClient struct {
	// scheduler is the scheduler instance.
	// scheduler 是调度器实例。
	scheduler common.SchedulerInterface
	// job is the monthly job being managed.
	// job 是正在管理的月度任务。
	job *jobs.MonthJob
}

// var _ MonthJobClientInterface = (*MonthJobClient)(nil) ensures that MonthJobClient implements MonthJobClientInterface.
// var _ MonthJobClientInterface = (*MonthJobClient)(nil) 确保 MonthJobClient 实现了 MonthJobClientInterface。
var _ MonthJobClientInterface = (*MonthJobClient)(nil)

// AtTime sets the specific days and time for the monthly jobs to run.
// AtTime 设置月度任务运行的特定日期和时间。
//
// Parameters:
//
//	days   - The days of the month when the jobs should run / 任务应该运行的月份日期
//	hour   - The hour of the day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) AtTime(days []int, hour, minute, second int) MonthJobClientInterface {
	m.job.AtTime(days, hour, minute, second)
	return m
}

// Alias sets the alias for the monthly jobs.
// Alias 设置月度任务的别名。
//
// Parameters:
//
//	alias - The alias string for the jobs / 任务的别名字符串
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) Alias(alias string) MonthJobClientInterface {
	m.job.Alias(alias)
	return m
}

// JobID sets the unique identifier for the monthly jobs.
// JobID 设置月度任务的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) JobID(id string) MonthJobClientInterface {
	m.job.JobID(id)
	return m
}

// Name sets the name for the monthly jobs.
// Name 设置月度任务的名称。
//
// Parameters:
//
//	name - The name string for the jobs / 任务的名称字符串
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) Name(name string) MonthJobClientInterface {
	m.job.Names(name)
	return m
}

// Tags adds tags to the monthly jobs.
// Tags 添加标签到月度任务。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) Tags(tags ...string) MonthJobClientInterface {
	m.job.Tags(tags...)
	return m
}

// Task sets the task function and its parameters for the monthly jobs.
// Task 设置月度任务的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) Task(task any, parameters ...any) MonthJobClientInterface {
	m.job.Task(task, parameters...)
	return m
}

// Watch sets a watcher function for jobs events.
// Watch 设置任务事件的监听函数。
//
// Parameters:
//
//	watch - The watch function for jobs events / 任务事件的监听函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) Watch(watch func(event common.JobWatchInterface)) MonthJobClientInterface {
	m.job.Watch(watch)
	return m
}

// DefaultHooks enables default event hooks for the monthly jobs.
// DefaultHooks 为月度任务启用默认的事件钩子。
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) DefaultHooks() MonthJobClientInterface {
	m.job.DefaultHooks()
	return m
}

// BeforeJobRuns sets a callback function to be executed before the jobs runs.
// BeforeJobRuns 设置在任务运行前执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface {
	m.job.BeforeJobRuns(eventListenerFunc)
	return m
}

// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 设置一个回调函数，如果返回错误则跳过任务执行。
//
// Parameters:
//
//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) MonthJobClientInterface {
	m.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return m
}

// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
// AfterJobRuns 设置在任务成功运行后执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface {
	m.job.AfterJobRuns(eventListenerFunc)
	return m
}

// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
// AfterJobRunsWithError 设置在任务运行出错时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface {
	m.job.AfterJobRunsWithError(eventListenerFunc)
	return m
}

// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
// AfterJobRunsWithPanic 设置在任务发生 panic 时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) MonthJobClientInterface {
	m.job.AfterJobRunsWithPanic(eventListenerFunc)
	return m
}

// AfterLockError sets a callback function to be executed when a lock error occurs.
// AfterLockError 设置在发生锁错误时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface {
	m.job.AfterLockError(eventListenerFunc)
	return m
}

// Add adds the configured monthly jobs to the scheduler.
// Add 将配置的月度任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
func (m *MonthJobClient) Add() (gocron.Job, error) {
	if m.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if m.job == nil {
		return nil, common.ErrMonthJobNil
	}
	return m.scheduler.AddMonthlyJob(m.job)
}

// BatchAdd adds multiple monthly jobs to the scheduler.
// BatchAdd 批量添加多个月度任务到调度器。
//
// Parameters:
//
//	monthlyJobs - Variable number of monthly jobs to add / 可变数量的要添加的月度任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (m *MonthJobClient) BatchAdd(monthlyJobs ...*jobs.MonthJob) ([]gocron.Job, error) {
	if m.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if m.job == nil {
		return nil, common.ErrMonthJobNil
	}
	jobs := make([]any, len(monthlyJobs))
	for i, j := range monthlyJobs {
		jobs[i] = j
	}
	return m.scheduler.AddMonthlyJobs(jobs...)
}

// Remove removes the monthly jobs from the scheduler.
// Remove 从调度器中删除月度任务。
//
// Returns:
//
//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
func (m *MonthJobClient) Remove() error {
	if m.scheduler == nil {
		return common.ErrScheduleNil
	}
	if m.job == nil {
		return common.ErrMonthJobNil
	}
	if m.job.ID != "" {
		return m.scheduler.RemoveJob(m.job.ID)
	}
	if m.job.Ali != "" {
		return m.scheduler.RemoveJobByAlias(m.job.Ali)
	}
	if m.job.Name != "" {
		return m.scheduler.RemoveJobByName(m.job.Name)
	}
	return common.ErrJobNotFound
}

// Get retrieves the monthly jobs from the scheduler.
// Get 从调度器中获取月度任务。
//
// Returns:
//
//	gocron.Job - The retrieved jobs / 获取的任务
//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
func (m *MonthJobClient) Get() (gocron.Job, error) {
	if m.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if m.job == nil {
		return nil, common.ErrMonthJobNil
	}
	if m.job.ID != "" {
		return m.scheduler.GetJobByID(m.job.ID)
	}
	if m.job.Ali != "" {
		return m.scheduler.GetJobByAlias(m.job.Ali)
	}
	if m.job.Name != "" {
		return m.scheduler.GetJobByName(m.job.Name)
	}
	return nil, common.ErrJobNotFound
}

// WithRetry sets the retry configuration for the jobs.
// 设置任务的重试配置。
//
// Parameters:
//
//	maxRetries - Maximum number of retries / 最大重试次数
//	policy     - Retry policy / 重试策略
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) WithRetry(maxRetries int, policy retry.RetryPolicy) MonthJobClientInterface {
	if m.job == nil {
		return m
	}
	if m.job.JobOptions == nil {
		m.job.JobOptions = common.NewJobOptions()
	}
	m.job.JobOptions.SetRetryConfig(&retry.RetryConfig{
		MaxRetries: maxRetries,
		Policy:     policy,
	})
	return m
}

// WithRetryConfig sets the complete retry configuration.
// 设置完整的重试配置。
//
// Parameters:
//
//	config - Retry configuration / 重试配置
//
// Returns:
//
//	MonthJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (m *MonthJobClient) WithRetryConfig(config *retry.RetryConfig) MonthJobClientInterface {
	if m.job == nil {
		return m
	}
	if m.job.JobOptions == nil {
		m.job.JobOptions = common.NewJobOptions()
	}
	m.job.JobOptions.SetRetryConfig(config)
	return m
}

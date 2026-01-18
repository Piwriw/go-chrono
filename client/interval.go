package client

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/retry"
)

// IntervalJobClientInterface defines the interface for configuring and managing interval jobs.
// 定义用于配置和管理间隔任务的接口。
type IntervalJobClientInterface interface {
	// Interval sets the time interval between jobs executions.
	// 设置任务执行之间的时间间隔。
	//
	// Parameters:
	//	interval - The duration between jobs executions / 任务执行之间的时间间隔
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Interval(interval time.Duration) IntervalJobClientInterface

	// Alias sets an alias for the jobs.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) IntervalJobClientInterface

	// JobID sets the unique identifier for the jobs.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The jobs ID string / 任务ID字符串
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) IntervalJobClientInterface

	// Name sets the name for the jobs.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The jobs name string / 任务名称字符串
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) IntervalJobClientInterface

	// Tag adds tags to the jobs for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tag(tags ...string) IntervalJobClientInterface

	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) IntervalJobClientInterface

	// Watch sets a watcher function to monitor jobs events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for jobs events / 任务事件的监视函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event common.JobWatchInterface)) IntervalJobClientInterface

	// DefaultHooks enables default event hooks for the jobs.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() IntervalJobClientInterface

	// BeforeJobRuns sets a callback function to be executed before the jobs runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface

	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) IntervalJobClientInterface

	// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface

	// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface

	// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) IntervalJobClientInterface

	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface

	// WithRetry sets the retry configuration for the jobs.
	// 设置任务的重试配置。
	//
	// Parameters:
	//	maxRetries - Maximum number of retries / 最大重试次数
	//	policy     - Retry policy / 重试策略
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetry(maxRetries int, policy retry.RetryPolicy) IntervalJobClientInterface

	// WithRetryConfig sets the complete retry configuration.
	// 设置完整的重试配置。
	//
	// Parameters:
	//	config - Retry configuration / 重试配置
	// Returns:
	//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetryConfig(config *retry.RetryConfig) IntervalJobClientInterface

	// Add adds the configured interval jobs to the scheduler.
	// 将配置的间隔任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added jobs / 添加的任务
	//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)

	// BatchAdd adds multiple interval jobs to the scheduler.
	// 批量添加多个间隔任务到调度器。
	//
	// Parameters:
	//	intervalJobs - Variable number of interval jobs to add / 可变数量的要添加的间隔任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(intervalJobs ...*jobs.IntervalJob) ([]gocron.Job, error)

	// Remove removes the interval jobs from the scheduler.
	// 从调度器中删除间隔任务。
	//
	// Returns:
	//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
	Remove() error

	// Get retrieves the interval jobs from the scheduler.
	// 从调度器中获取间隔任务。
	//
	// Returns:
	//	gocron.Job - The retrieved jobs / 获取的任务
	//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// IntervalJobClient is a client for managing interval jobs with a scheduler.
// IntervalJobClient 是一个用于通过调度器管理间隔任务的客户端。
type IntervalJobClient struct {
	// scheduler is the scheduler instance.
	// scheduler 是调度器实例。
	scheduler common.SchedulerInterface
	// job is the interval job being managed.
	// job 是正在管理的间隔任务。
	job *jobs.IntervalJob
}

// var _ IntervalJobClientInterface = (*IntervalJobClient)(nil) ensures that IntervalJobClient implements IntervalJobClientInterface.
// var _ IntervalJobClientInterface = (*IntervalJobClient)(nil) 确保 IntervalJobClient 实现了 IntervalJobClientInterface。
var _ IntervalJobClientInterface = (*IntervalJobClient)(nil)

// Interval sets the time interval for the interval jobs.
// Interval 设置间隔任务的时间间隔。
//
// Parameters:
//
//	interval - The time duration between jobs executions / 任务执行之间的时间间隔
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Interval(interval time.Duration) IntervalJobClientInterface {
	c.job.IntervalTime(interval)
	return c
}

// Alias sets the alias for the interval jobs.
// Alias 设置间隔任务的别名。
//
// Parameters:
//
//	alias - The alias string for the jobs / 任务的别名字符串
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Alias(alias string) IntervalJobClientInterface {
	c.job.Alias(alias)
	return c
}

// JobID sets the unique identifier for the interval jobs.
// JobID 设置间隔任务的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) JobID(id string) IntervalJobClientInterface {
	c.job.JobID(id)
	return c
}

// Name sets the name for the interval jobs.
// Name 设置间隔任务的名称。
//
// Parameters:
//
//	name - The name string for the jobs / 任务的名称字符串
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Name(name string) IntervalJobClientInterface {
	c.job.Names(name)
	return c
}

// Tag adds tags to the interval jobs.
// Tag 添加标签到间隔任务。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Tag(tags ...string) IntervalJobClientInterface {
	c.job.Tag(tags...)
	return c
}

// Task sets the task function and its parameters for the interval jobs.
// Task 设置间隔任务的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Task(task any, parameters ...any) IntervalJobClientInterface {
	c.job.Task(task, parameters...)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) Watch(watch func(event common.JobWatchInterface)) IntervalJobClientInterface {
	c.job.Watch(watch)
	return c
}

// DefaultHooks enables default event hooks for the interval jobs.
// DefaultHooks 为间隔任务启用默认的事件钩子。
//
// Returns:
//
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) DefaultHooks() IntervalJobClientInterface {
	c.job.DefaultHooks()
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface {
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) IntervalJobClientInterface {
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface {
	c.job.AfterJobRuns(eventListenerFunc)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface {
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) IntervalJobClientInterface {
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface {
	c.job.AfterLockError(eventListenerFunc)
	return c
}

// Add adds the configured interval jobs to the scheduler.
// Add 将配置的间隔任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
func (c *IntervalJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if c.job == nil {
		return nil, common.ErrIntervalJobNil
	}
	return c.scheduler.AddIntervalJob(c.job)
}

// BatchAdd adds multiple interval jobs to the scheduler.
// BatchAdd 批量添加多个间隔任务到调度器。
//
// Parameters:
//
//	intervalJobs - Variable number of interval jobs to add / 可变数量的要添加的间隔任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (c *IntervalJobClient) BatchAdd(intervalJobs ...*jobs.IntervalJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if c.job == nil {
		return nil, common.ErrIntervalJobNil
	}
	jobs := make([]any, len(intervalJobs))
	for i, j := range intervalJobs {
		jobs[i] = j
	}
	return c.scheduler.AddIntervalJobs(jobs...)
}

// Remove removes the interval jobs from the scheduler.
// Remove 从调度器中删除间隔任务。
//
// Returns:
//
//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
func (c *IntervalJobClient) Remove() error {
	if c.scheduler == nil {
		return common.ErrScheduleNil
	}
	if c.job == nil {
		return common.ErrIntervalJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.RemoveJob(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.RemoveJob(c.job.Name)
	}
	return common.ErrJobNotFound
}

// Get retrieves the interval jobs from the scheduler.
// Get 从调度器中获取间隔任务。
//
// Returns:
//
//	gocron.Job - The retrieved jobs / 获取的任务
//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
func (c *IntervalJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if c.job == nil {
		return nil, common.ErrIntervalJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.GetJobByID(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.GetJobByID(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.GetJobByName(c.job.Name)
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) WithRetry(maxRetries int, policy retry.RetryPolicy) IntervalJobClientInterface {
	if c.job == nil {
		return c
	}
	if c.job.JobOptions == nil {
		c.job.JobOptions = common.NewJobOptions()
	}
	c.job.JobOptions.SetRetryConfig(&retry.RetryConfig{
		MaxRetries: maxRetries,
		Policy:     policy,
	})
	return c
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
//	IntervalJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *IntervalJobClient) WithRetryConfig(config *retry.RetryConfig) IntervalJobClientInterface {
	if c.job == nil {
		return c
	}
	if c.job.JobOptions == nil {
		c.job.JobOptions = common.NewJobOptions()
	}
	c.job.JobOptions.SetRetryConfig(config)
	return c
}

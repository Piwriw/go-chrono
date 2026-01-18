package client

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/retry"
)

// OnceJobClientInterface defines the interface for configuring and managing one-time jobs.
// 定义用于配置和管理一次性任务的接口。
type OnceJobClientInterface interface {
	// AtTimes sets the specific times when the one-time jobs should run.
	// 设置一次性任务运行的特定时间。
	//
	// Parameters:
	//	workTimes - Variable number of time values / 可变数量的时间值
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AtTimes(workTimes ...time.Time) OnceJobClientInterface

	// Alias sets an alias for the jobs.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) OnceJobClientInterface

	// JobID sets the unique identifier for the jobs.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The jobs ID string / 任务ID字符串
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) OnceJobClientInterface

	// Name sets the name for the jobs.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The jobs name string / 任务名称字符串
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) OnceJobClientInterface

	// Tag adds tags to the jobs for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tag(tags ...string) OnceJobClientInterface

	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) OnceJobClientInterface

	// Watch sets a watcher function to monitor jobs events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for jobs events / 任务事件的监视函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event common.JobWatchInterface)) OnceJobClientInterface

	// DefaultHooks enables default event hooks for the jobs.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() OnceJobClientInterface

	// BeforeJobRuns sets a callback function to be executed before the jobs runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface

	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) OnceJobClientInterface

	// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface

	// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface

	// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) OnceJobClientInterface

	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface

	// WithRetry sets the retry configuration for the jobs.
	// 设置任务的重试配置。
	//
	// Parameters:
	//	maxRetries - Maximum number of retries / 最大重试次数
	//	policy     - Retry policy / 重试策略
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetry(maxRetries int, policy retry.RetryPolicy) OnceJobClientInterface

	// WithRetryConfig sets the complete retry configuration.
	// 设置完整的重试配置。
	//
	// Parameters:
	//	config - Retry configuration / 重试配置
	// Returns:
	//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetryConfig(config *retry.RetryConfig) OnceJobClientInterface

	// Add adds the configured one-time jobs to the scheduler.
	// 将配置的一次性任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added jobs / 添加的任务
	//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)

	// BatchAdd adds multiple one-time jobs to the scheduler.
	// 批量添加多个一次性任务到调度器。
	//
	// Parameters:
	//	onceJobs - Variable number of one-time jobs to add / 可变数量的要添加的一次性任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(onceJobs ...*jobs.OnceJob) ([]gocron.Job, error)

	// Remove removes the one-time jobs from the scheduler.
	// 从调度器中删除一次性任务。
	//
	// Returns:
	//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
	Remove() error

	// Get retrieves the one-time jobs from the scheduler.
	// 从调度器中获取一次性任务。
	//
	// Returns:
	//	gocron.Job - The retrieved jobs / 获取的任务
	//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// OnceJobClient is a client for managing one-time jobs with a scheduler.
// OnceJobClient 是一个用于通过调度器管理一次性任务的客户端。
type OnceJobClient struct {
	// scheduler is the scheduler instance.
	// scheduler 是调度器实例。
	scheduler common.SchedulerInterface
	// job is the one-time job being managed.
	// job 是正在管理的一次性任务。
	job *jobs.OnceJob
}

// var _ OnceJobClientInterface = (*OnceJobClient)(nil) ensures that OnceJobClient implements OnceJobClientInterface.
// var _ OnceJobClientInterface = (*OnceJobClient)(nil) 确保 OnceJobClient 实现了 OnceJobClientInterface。
var _ OnceJobClientInterface = (*OnceJobClient)(nil)

// AtTimes sets the specific times when the one-time jobs should run.
// AtTimes 设置一次性任务运行的特定时间。
//
// Parameters:
//
//	workTimes - Variable number of time values / 可变数量的时间值
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) AtTimes(workTimes ...time.Time) OnceJobClientInterface {
	o.job.AtTimes(workTimes...)
	return o
}

// Alias sets the alias for the one-time jobs.
// Alias 设置一次性任务的别名。
//
// Parameters:
//
//	alias - The alias string for the jobs / 任务的别名字符串
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) Alias(alias string) OnceJobClientInterface {
	o.job.Alias(alias)
	return o
}

// JobID sets the unique identifier for the one-time jobs.
// JobID 设置一次性任务的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) JobID(id string) OnceJobClientInterface {
	o.job.JobID(id)
	return o
}

// Name sets the name for the one-time jobs.
// Name 设置一次性任务的名称。
//
// Parameters:
//
//	name - The name string for the jobs / 任务的名称字符串
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) Name(name string) OnceJobClientInterface {
	o.job.Names(name)
	return o
}

// Tag adds tags to the one-time jobs.
// Tag 添加标签到一次性任务。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) Tag(tags ...string) OnceJobClientInterface {
	o.job.Tags(tags...)
	return o
}

// Task sets the task function and its parameters for the one-time jobs.
// Task 设置一次性任务的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) Task(task any, parameters ...any) OnceJobClientInterface {
	o.job.Task(task, parameters...)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) Watch(watch func(event common.JobWatchInterface)) OnceJobClientInterface {
	o.job.Watch(watch)
	return o
}

// DefaultHooks enables default event hooks for the one-time jobs.
// DefaultHooks 为一次性任务启用默认的事件钩子。
//
// Returns:
//
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) DefaultHooks() OnceJobClientInterface {
	o.job.DefaultHooks()
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface {
	o.job.BeforeJobRuns(eventListenerFunc)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) OnceJobClientInterface {
	o.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface {
	o.job.AfterJobRuns(eventListenerFunc)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface {
	o.job.AfterJobRunsWithError(eventListenerFunc)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) OnceJobClientInterface {
	o.job.AfterJobRunsWithPanic(eventListenerFunc)
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface {
	o.job.AfterLockError(eventListenerFunc)
	return o
}

// Add adds the configured one-time jobs to the scheduler.
// Add 将配置的一次性任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
func (o *OnceJobClient) Add() (gocron.Job, error) {
	if o.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if o.job == nil {
		return nil, common.ErrOnceJobNil
	}
	return o.scheduler.AddOnceJob(o.job)
}

// BatchAdd adds multiple one-time jobs to the scheduler.
// BatchAdd 批量添加多个一次性任务到调度器。
//
// Parameters:
//
//	onceJobs - Variable number of one-time jobs to add / 可变数量的要添加的一次性任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (o *OnceJobClient) BatchAdd(onceJobs ...*jobs.OnceJob) ([]gocron.Job, error) {
	if o.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	jobs := make([]any, len(onceJobs))
	for i, j := range onceJobs {
		jobs[i] = j
	}
	return o.scheduler.AddOnceJobs(jobs...)
}

// Remove removes the one-time jobs from the scheduler.
// Remove 从调度器中删除一次性任务。
//
// Returns:
//
//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
func (o *OnceJobClient) Remove() error {
	if o.scheduler == nil {
		return common.ErrScheduleNil
	}
	if o.job == nil {
		return common.ErrOnceJobNil
	}
	if o.job.ID != "" {
		return o.scheduler.RemoveJob(o.job.ID)
	}
	if o.job.Ali != "" {
		return o.scheduler.RemoveJobByAlias(o.job.Ali)
	}
	if o.job.Name != "" {
		return o.scheduler.RemoveJobByName(o.job.Name)
	}
	return common.ErrJobNotFound
}

// Get retrieves the one-time jobs from the scheduler.
// Get 从调度器中获取一次性任务。
//
// Returns:
//
//	gocron.Job - The retrieved jobs / 获取的任务
//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
func (o *OnceJobClient) Get() (gocron.Job, error) {
	if o.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if o.job == nil {
		return nil, common.ErrOnceJobNil
	}
	if o.job.ID != "" {
		return o.scheduler.GetJobByID(o.job.ID)
	}
	if o.job.Ali != "" {
		return o.scheduler.GetJobByAlias(o.job.Ali)
	}
	if o.job.Name != "" {
		return o.scheduler.GetJobByName(o.job.Name)
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) WithRetry(maxRetries int, policy retry.RetryPolicy) OnceJobClientInterface {
	if o.job == nil {
		return o
	}
	if o.job.JobOptions == nil {
		o.job.JobOptions = common.NewJobOptions()
	}
	o.job.JobOptions.SetRetryConfig(&retry.RetryConfig{
		MaxRetries: maxRetries,
		Policy:     policy,
	})
	return o
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
//	OnceJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (o *OnceJobClient) WithRetryConfig(config *retry.RetryConfig) OnceJobClientInterface {
	if o.job == nil {
		return o
	}
	if o.job.JobOptions == nil {
		o.job.JobOptions = common.NewJobOptions()
	}
	o.job.JobOptions.SetRetryConfig(config)
	return o
}

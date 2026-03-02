package client

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/retry"
)

// 定义用于配置和管理每日任务的接口。
type DailyJobClientInterface interface {
	// AtDayTime sets the jobs to run at a specific time every day.
	// 设置任务在每天的指定时间运行。
	//
	// Parameters:
	//	hour   - The hour (0-23) / 小时（0-23）
	//	minute - The minute (0-59) / 分钟（0-59）
	//	second - The second (0-59) / 秒（0-59）
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AtDayTime(hour, minute, second uint) DailyJobClientInterface
	// Alias sets an alias for the jobs.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) DailyJobClientInterface
	// JobID sets the unique identifier for the jobs.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The jobs ID string / 任务ID字符串
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) DailyJobClientInterface
	// Name sets the name for the jobs.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The jobs name string / 任务名称字符串
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) DailyJobClientInterface
	// Tags adds tags to the jobs for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tags(tags ...string) DailyJobClientInterface
	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) DailyJobClientInterface
	// Watch sets a watcher function to monitor jobs events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for jobs events / 任务事件的监视函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event common.JobWatchInterface)) DailyJobClientInterface
	// DefaultHooks enables default event hooks for the jobs.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() DailyJobClientInterface
	// BeforeJobRuns sets a callback function to be executed before the jobs runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface
	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) DailyJobClientInterface
	// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface
	// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface
	// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) DailyJobClientInterface
	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface

	// WithRetry sets the retry configuration for the jobs.
	// 设置任务的重试配置。
	//
	// Parameters:
	//	maxRetries - Maximum number of retries / 最大重试次数
	//	policy     - Retry policy / 重试策略
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetry(maxRetries int, policy retry.RetryPolicy) DailyJobClientInterface

	// WithRetryConfig sets the complete retry configuration.
	// 设置完整的重试配置。
	//
	// Parameters:
	//	config - Retry configuration / 重试配置
	// Returns:
	//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	WithRetryConfig(config *retry.RetryConfig) DailyJobClientInterface

	// Add adds the configured daily jobs to the scheduler.
	// 将配置的每日任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added jobs / 添加的任务
	//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)
	// BatchAdd adds multiple daily jobs to the scheduler.
	// 批量添加多个每日任务到调度器。
	//
	// Parameters:
	//	dailyJobs - Variable number of daily jobs to add / 可变数量的要添加的每日任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(dailyJobs ...*jobs.DailyJob) ([]gocron.Job, error)
	// Remove removes the daily jobs from the scheduler.
	// 从调度器中删除每日任务。
	//
	// Returns:
	//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
	Remove() error
	// Get retrieves the daily jobs from the scheduler.
	// 从调度器中获取每日任务。
	//
	// Returns:
	//	gocron.Job - The retrieved jobs / 获取的任务
	//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// DailyJobClient provides a fluent interface for configuring and managing daily jobs.
// 提供用于配置和管理每日任务的流畅接口。
//
// It acts as a builder pattern wrapper around DailyJob, delegating configuration
// to the underlying jobs and managing operations through the scheduler.
// 它作为 DailyJob 的构建器模式包装器，将配置委托给底层任务，并通过调度器管理操作。
type DailyJobClient struct {
	// scheduler is the scheduler instance / 调度器实例
	scheduler common.SchedulerInterface
	// job is the daily job being configured / 正在配置的每日任务
	job *jobs.DailyJob
}

// Ensure DailyJobClient implements DailyJobClientInterface at compile time.
// 在编译时确保 DailyJobClient 实现了 DailyJobClientInterface。
var _ DailyJobClientInterface = (*DailyJobClient)(nil)

// AtDayTime sets the jobs to run at a specific time every day.
// 设置任务在每天的指定时间运行。
//
// Parameters:
//
//	hour   - The hour (0-23) / 小时（0-23）
//	minute - The minute (0-59) / 分钟（0-59）
//	second - The second (0-59) / 秒（0-59）
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) AtDayTime(hour, minute, second uint) DailyJobClientInterface {
	c.job.AtDayTime(hour, minute, second)
	return c
}

// Alias sets an alias for the jobs.
// 为任务设置别名。
//
// Parameters:
//
//	alias - The alias string / 别名字符串
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) Alias(alias string) DailyJobClientInterface {
	c.job.Alias(alias)
	return c
}

// JobID sets the unique identifier for the jobs.
// 设置任务的唯一标识符。
//
// Parameters:
//
//	id - The jobs ID string / 任务ID字符串
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) JobID(id string) DailyJobClientInterface {
	c.job.JobID(id)
	return c
}

// Name sets the name for the jobs.
// 为任务设置名称。
//
// Parameters:
//
//	name - The jobs name string / 任务名称字符串
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) Name(name string) DailyJobClientInterface {
	c.job.Names(name)
	return c
}

// Tags adds tags to the jobs for categorization.
// 为任务添加标签以便分类。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) Tags(tags ...string) DailyJobClientInterface {
	c.job.Tag(tags...)
	return c
}

// Task sets the task function and its parameters to be executed.
// 设置要执行的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters to pass to the task / 传递给任务的可变参数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) Task(task any, parameters ...any) DailyJobClientInterface {
	c.job.Task(task, parameters...)
	return c
}

// Watch sets a watcher function to monitor jobs events.
// 设置监视器函数以监控任务事件。
//
// Parameters:
//
//	watch - The watch function for jobs events / 任务事件的监视函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) Watch(watch func(event common.JobWatchInterface)) DailyJobClientInterface {
	c.job.Watch(watch)
	return c
}

// DefaultHooks enables default event hooks for the jobs.
// 为任务启用默认的事件钩子。
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) DefaultHooks() DailyJobClientInterface {
	c.job.DefaultHooks()
	return c
}

// BeforeJobRuns sets a callback function to be executed before the jobs runs.
// 设置在任务运行前执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface {
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
}

// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip jobs execution if it returns an error.
// 设置一个回调函数，如果返回错误则跳过任务执行。
//
// Parameters:
//
//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) DailyJobClientInterface {
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

// AfterJobRuns sets a callback function to be executed after the jobs runs successfully.
// 设置在任务成功运行后执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface {
	c.job.AfterJobRuns(eventListenerFunc)
	return c
}

// AfterJobRunsWithError sets a callback function to be executed when the jobs runs with an error.
// 设置在任务运行出错时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface {
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
}

// AfterJobRunsWithPanic sets a callback function to be executed when the jobs panics.
// 设置在任务发生 panic 时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) DailyJobClientInterface {
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

// AfterLockError sets a callback function to be executed when a lock error occurs.
// 设置在发生锁错误时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface {
	c.job.AfterLockError(eventListenerFunc)
	return c
}

// Add adds the configured daily jobs to the scheduler.
// 将配置的每日任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if the jobs cannot be added / 如果无法添加任务的错误
func (c *DailyJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if c.job == nil {
		return nil, common.ErrDailyJobNil
	}
	return c.scheduler.AddDailyJob(c.job)
}

// BatchAdd adds multiple daily jobs to the scheduler.
// 批量添加多个每日任务到调度器。
//
// Parameters:
//
//	dailyJobs - Variable number of daily jobs to add / 可变数量的要添加的每日任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (c *DailyJobClient) BatchAdd(dailyJobs ...*jobs.DailyJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	jobs := make([]any, len(dailyJobs))
	for i, j := range dailyJobs {
		jobs[i] = j
	}
	return c.scheduler.AddDailyJobs(jobs...)
}

// Remove removes the daily jobs from the scheduler.
// 从调度器中删除每日任务。
//
// The removal strategy follows a priority order:
// 1. Remove by JobID if set
// 2. Remove by Alias if JobID is not set
// 3. Remove by Name if neither JobID nor Alias is set
//
// 删除策略按优先级顺序执行：
// 1. 如果设置了 JobID，则通过 JobID 删除
// 2. 如果未设置 JobID，则通过 Alias 删除
// 3. 如果既未设置 JobID 也未设置 Alias，则通过 Name 删除
//
// Returns:
//
//	error - Error if the jobs cannot be removed / 如果无法删除任务的错误
func (c *DailyJobClient) Remove() error {
	if c.scheduler == nil {
		return common.ErrScheduleNil
	}
	if c.job == nil {
		return common.ErrCronJobNil
	}
	// Priority 1: Try to remove by JobID
	// 优先级 1：尝试通过 JobID 删除
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	// Priority 2: Try to remove by Alias
	// 优先级 2：尝试通过 Alias 删除
	if c.job.Ali != "" {
		return c.scheduler.RemoveJobByAlias(c.job.Ali)
	}
	// Priority 3: Try to remove by Name
	// 优先级 3：尝试通过 Name 删除
	if c.job.Name != "" {
		return c.scheduler.RemoveJobByName(c.job.Name)
	}
	return common.ErrJobNotFound
}

// Get retrieves the daily jobs from the scheduler.
// 从调度器中获取每日任务。
//
// The retrieval strategy follows a priority order:
// 1. Get by JobID if set
// 2. Get by Alias if JobID is not set
// 3. Get by Name if neither JobID nor Alias is set
//
// 获取策略按优先级顺序执行：
// 1. 如果设置了 JobID，则通过 JobID 获取
// 2. 如果未设置 JobID，则通过 Alias 获取
// 3. 如果既未设置 JobID 也未设置 Alias，则通过 Name 获取
//
// Returns:
//
//	gocron.Job - The retrieved jobs / 获取的任务
//	error      - Error if the jobs cannot be found / 如果找不到任务的错误
func (c *DailyJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, common.ErrScheduleNil
	}
	if c.job == nil {
		return nil, common.ErrCronJobNil
	}
	// Priority 1: Try to get by JobID
	// 优先级 1：尝试通过 JobID 获取
	if c.job.ID != "" {
		return c.scheduler.GetJobByID(c.job.ID)
	}
	// Priority 2: Try to get by Alias
	// 优先级 2：尝试通过 Alias 获取
	if c.job.Ali != "" {
		return c.scheduler.GetJobByAlias(c.job.Ali)
	}
	// Priority 3: Try to get by Name
	// 优先级 3：尝试通过 Name 获取
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
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) WithRetry(maxRetries int, policy retry.RetryPolicy) DailyJobClientInterface {
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
//	DailyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *DailyJobClient) WithRetryConfig(config *retry.RetryConfig) DailyJobClientInterface {
	if c.job == nil {
		return c
	}
	if c.job.JobOptions == nil {
		c.job.JobOptions = common.NewJobOptions()
	}
	c.job.JobOptions.SetRetryConfig(config)
	return c
}

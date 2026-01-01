package chrono

import (
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// CronJobClientInterface defines the interface for configuring and managing cron jobs.
// 定义用于配置和管理定时任务的接口。
type CronJobClientInterface interface {
	// CronExpr sets the Linux cron expression for the job (required field).
	// 设置任务的 Linux Cron 表达式（必填字段）。
	//
	// Parameters:
	//	expr - The cron expression string / Cron 表达式字符串
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	CronExpr(expr string) CronJobClientInterface

	// Alias sets an alias for the job.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) CronJobClientInterface

	// JobID sets the unique identifier for the job.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The job ID string / 任务ID字符串
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) CronJobClientInterface

	// Name sets the name for the job.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The job name string / 任务名称字符串
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) CronJobClientInterface

	// Tags adds tags to the job for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tags(tags ...string) CronJobClientInterface

	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) CronJobClientInterface

	// Watch sets a watcher function to monitor job events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for job events / 任务事件的监视函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event JobWatchInterface)) CronJobClientInterface

	// DefaultHooks enables default event hooks for the job.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() CronJobClientInterface

	// BeforeJobRuns sets a callback function to be executed before the job runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface

	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip job execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) CronJobClientInterface

	// AfterJobRuns sets a callback function to be executed after the job runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface

	// AfterJobRunsWithError sets a callback function to be executed when the job runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface

	// AfterJobRunsWithPanic sets a callback function to be executed when the job panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) CronJobClientInterface

	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface

	// Add adds the configured cron job to the scheduler.
	// 将配置的定时任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added job / 添加的任务
	//	error      - Error if the job cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)

	// BatchAdd adds multiple cron jobs to the scheduler.
	// 批量添加多个定时任务到调度器。
	//
	// Parameters:
	//	cronJobs - Variable number of cron jobs to add / 可变数量的要添加的定时任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error)

	// Remove removes the cron job from the scheduler.
	// 从调度器中删除定时任务。
	//
	// Returns:
	//	error - Error if the job cannot be removed / 如果无法删除任务的错误
	Remove() error

	// Get retrieves the cron job from the scheduler.
	// 从调度器中获取定时任务。
	//
	// Returns:
	//	gocron.Job - The retrieved job / 获取的任务
	//	error      - Error if the job cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// CronJobClient provides a fluent interface for configuring and managing cron jobs.
// 提供用于配置和管理定时任务的流畅接口。
//
// It acts as a builder pattern wrapper around CronJob, delegating configuration
// to the underlying job and managing operations through the scheduler.
// 它作为 CronJob 的构建器模式包装器，将配置委托给底层任务，并通过调度器管理操作。
//
// The client uses a mutex to ensure thread-safe access to the scheduler and job fields.
// 客户端使用互斥锁确保对调度器和任务字段的线程安全访问。
type CronJobClient struct {
	mu        sync.RWMutex
	scheduler *Scheduler
	job       *CronJob
}

// Ensure CronJobClient implements CronJobClientInterface at compile time.
// 在编译时确保 CronJobClient 实现了 CronJobClientInterface。
var _ CronJobClientInterface = (*CronJobClient)(nil)

// CronExpr sets the Linux cron expression for the job.
// 设置任务的 Linux Cron 表达式。
//
// Parameters:
//
//	expr - The cron expression string / Cron 表达式字符串
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) CronExpr(expr string) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.CronExpr(expr)
	return c
}

// Alias sets an alias for the job.
// 为任务设置别名。
//
// Parameters:
//
//	alias - The alias string / 别名字符串
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) Alias(alias string) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.Alias(alias)
	return c
}

// JobID sets the unique identifier for the job.
// 设置任务的唯一标识符。
//
// Parameters:
//
//	id - The job ID string / 任务ID字符串
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) JobID(id string) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.JobID(id)
	return c
}

// Name sets the name for the job.
// 为任务设置名称。
//
// Parameters:
//
//	name - The job name string / 任务名称字符串
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) Name(name string) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.Names(name)
	return c
}

// Tags adds tags to the job for categorization.
// 为任务添加标签以便分类。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) Tags(tags ...string) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.Tag(tags...)
	return c
}

// Task sets the task function and its parameters to be executed.
// 设置要执行的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function / 任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) Task(task any, parameters ...any) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.Task(task, parameters...)
	return c
}

// Watch sets a watcher function to monitor job events.
// 设置监视器函数以监控任务事件。
//
// Parameters:
//
//	watch - The watch function for job events / 任务事件的监视函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) Watch(watch func(event JobWatchInterface)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.Watch(watch)
	return c
}

// DefaultHooks enables default event hooks for the job.
// 为任务启用默认的事件钩子。
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) DefaultHooks() CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.DefaultHooks()
	return c
}

// BeforeJobRuns sets a callback function to be executed before the job runs.
// 设置在任务运行前执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
}

// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip job execution if it returns an error.
// 设置一个回调函数，如果返回错误则跳过任务执行。
//
// Parameters:
//
//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

// AfterJobRuns sets a callback function to be executed after the job runs successfully.
// 设置在任务成功运行后执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.AfterJobRuns(eventListenerFunc)
	return c
}

// AfterJobRunsWithError sets a callback function to be executed when the job runs with an error.
// 设置在任务运行出错时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
}

// AfterJobRunsWithPanic sets a callback function to be executed when the job panics.
// 设置在任务发生 panic 时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

// AfterLockError sets a callback function to be executed when a lock error occurs.
// 设置在发生锁错误时执行的回调函数。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	CronJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (c *CronJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.job == nil {
		return c
	}
	c.job.AfterLockError(eventListenerFunc)
	return c
}

// Add adds the configured cron job to the scheduler.
// 将配置的定时任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if the job cannot be added / 如果无法添加任务的错误
func (c *CronJobClient) Add() (gocron.Job, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
	}
	return c.scheduler.AddCronJob(c.job)
}

// BatchAdd adds multiple cron jobs to the scheduler.
// 批量添加多个定时任务到调度器。
//
// Parameters:
//
//	cronJobs - Variable number of cron jobs to add / 可变数量的要添加的定时任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (c *CronJobClient) BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if len(cronJobs) == 0 {
		return nil, nil
	}
	return c.scheduler.AddCronJobs(cronJobs...)
}

// Remove removes the cron job from the scheduler.
// 从调度器中删除定时任务。
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
//	error - Error if the job cannot be removed / 如果无法删除任务的错误
func (c *CronJobClient) Remove() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job == nil {
		return ErrCronJobNil
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
	return ErrJobNotFound
}

// Get retrieves the cron job from the scheduler.
// 从调度器中获取定时任务。
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
//	gocron.Job - The retrieved job / 获取的任务
//	error      - Error if the job cannot be found / 如果找不到任务的错误
func (c *CronJobClient) Get() (gocron.Job, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
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
	return nil, ErrJobNotFound
}

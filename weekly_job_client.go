package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// WeeklyJobClientInterface defines interface for configuring and managing weekly jobs.
// 定义用于配置和管理周任务的接口。
type WeeklyJobClientInterface interface {
	// AtTimes sets specific days and time for the weekly job to run.
	// 设置周任务运行的特定日期和时间。
	//
	// Parameters:
	//	days   - The days of the week when the job should run / 任务应该运行的星期几
	//	hour   - The hour of day (0-23) / 一天中的小时（0-23）
	//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
	//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AtTimes(days []time.Weekday, hour, minute, second uint) WeeklyJobClientInterface

	// Alias sets an alias for the job.
	// 为任务设置别名。
	//
	// Parameters:
	//	alias - The alias string / 别名字符串
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Alias(alias string) WeeklyJobClientInterface

	// JobID sets the unique identifier for the job.
	// 设置任务的唯一标识符。
	//
	// Parameters:
	//	id - The job ID string / 任务ID字符串
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	JobID(id string) WeeklyJobClientInterface

	// Name sets the name for the job.
	// 为任务设置名称。
	//
	// Parameters:
	//	name - The job name string / 任务名称字符串
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Name(name string) WeeklyJobClientInterface

	// Tags adds tags to the job for categorization.
	// 为任务添加标签以便分类。
	//
	// Parameters:
	//	tags - Variable number of tag strings / 可变数量的标签字符串
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Tags(tags ...string) WeeklyJobClientInterface

	// Task sets the task function and its parameters to be executed.
	// 设置要执行的任务函数及其参数。
	//
	// Parameters:
	//	task      - The task function / 任务函数
	//	parameters - Variable parameters for the task / 任务的可变参数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Task(task any, parameters ...any) WeeklyJobClientInterface

	// Watch sets a watcher function to monitor job events.
	// 设置监视器函数以监控任务事件。
	//
	// Parameters:
	//	watch - The watch function for job events / 任务事件的监视函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	Watch(watch func(event JobWatchInterface)) WeeklyJobClientInterface

	// DefaultHooks enables default event hooks for the job.
	// 为任务启用默认的事件钩子。
	//
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	DefaultHooks() WeeklyJobClientInterface

	// BeforeJobRuns sets a callback function to be executed before the job runs.
	// 设置在任务运行前执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface

	// BeforeJobRunsSkipIfBeforeFuncErrors sets a callback that can skip job execution if it returns an error.
	// 设置一个回调函数，如果返回错误则跳过任务执行。
	//
	// Parameters:
	//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) WeeklyJobClientInterface

	// AfterJobRuns sets a callback function to be executed after the job runs successfully.
	// 设置在任务成功运行后执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface

	// AfterJobRunsWithError sets a callback function to be executed when the job runs with an error.
	// 设置在任务运行出错时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface

	// AfterJobRunsWithPanic sets a callback function to be executed when the job panics.
	// 设置在任务发生 panic 时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) WeeklyJobClientInterface

	// AfterLockError sets a callback function to be executed when a lock error occurs.
	// 设置在发生锁错误时执行的回调函数。
	//
	// Parameters:
	//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
	// Returns:
	//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface

	// Add adds the configured weekly job to the scheduler.
	// 将配置的周任务添加到调度器。
	//
	// Returns:
	//	gocron.Job - The added job / 添加的任务
	//	error      - Error if the job cannot be added / 如果无法添加任务的错误
	Add() (gocron.Job, error)

	// BatchAdd adds multiple weekly jobs to the scheduler.
	// 批量添加多个周任务到调度器。
	//
	// Parameters:
	//	weeklyJobs - Variable number of weekly jobs to add / 可变数量的要添加的周任务
	// Returns:
	//	[]gocron.Job - The list of added jobs / 添加的任务列表
	//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
	BatchAdd(weeklyJobs ...*WeeklyJob) ([]gocron.Job, error)

	// Remove removes the weekly job from the scheduler.
	// 从调度器中删除周任务。
	//
	// Returns:
	//	error - Error if the job cannot be removed / 如果无法删除任务的错误
	Remove() error

	// Get retrieves the weekly job from the scheduler.
	// 从调度器中获取周任务。
	//
	// Returns:
	//	gocron.Job - The retrieved job / 获取的任务
	//	error      - Error if the job cannot be found / 如果找不到任务的错误
	Get() (gocron.Job, error)
}

// WeeklyJobClient is the client implementation for configuring and managing weekly jobs.
// 用于配置和管理周任务的客户端实现。
type WeeklyJobClient struct {
	// scheduler is the scheduler instance that manages the weekly jobs.
	// 管理周任务的调度器实例。
	scheduler *Scheduler
	// job is the weekly job configuration.
	// 周任务配置。
	job *WeeklyJob
}

var _ WeeklyJobClientInterface = (*WeeklyJobClient)(nil)

// AtTimes sets specific days and time for the weekly job to run.
// 设置周任务运行的特定日期和时间。
//
// Parameters:
//
//	days   - The days of the week when the job should run / 任务应该运行的星期几
//	hour   - The hour of day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) AtTimes(days []time.Weekday, hour, minute, second uint) WeeklyJobClientInterface {
	w.job.AtTimes(days, hour, minute, second)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) Alias(alias string) WeeklyJobClientInterface {
	w.job.Alias(alias)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) JobID(id string) WeeklyJobClientInterface {
	w.job.JobID(id)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) Name(name string) WeeklyJobClientInterface {
	w.job.Names(name)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) Tags(tags ...string) WeeklyJobClientInterface {
	w.job.Tags(tags...)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) Task(task any, parameters ...any) WeeklyJobClientInterface {
	w.job.Task(task, parameters...)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) Watch(watch func(event JobWatchInterface)) WeeklyJobClientInterface {
	w.job.Watch(watch)
	return w
}

// DefaultHooks enables default event hooks for the job.
// 为任务启用默认的事件钩子。
//
// Returns:
//
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) DefaultHooks() WeeklyJobClientInterface {
	w.job.DefaultHooks()
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface {
	w.job.BeforeJobRuns(eventListenerFunc)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) WeeklyJobClientInterface {
	w.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface {
	w.job.AfterJobRuns(eventListenerFunc)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface {
	w.job.AfterJobRunsWithError(eventListenerFunc)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) WeeklyJobClientInterface {
	w.job.AfterJobRunsWithPanic(eventListenerFunc)
	return w
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
//	WeeklyJobClientInterface - The client interface for method chaining / 客户端接口，支持链式调用
func (w *WeeklyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface {
	w.job.AfterLockError(eventListenerFunc)
	return w
}

// Add adds the configured weekly job to the scheduler.
// 将配置的周任务添加到调度器。
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if the job cannot be added / 如果无法添加任务的错误
func (w *WeeklyJobClient) Add() (gocron.Job, error) {
	// Check if scheduler is nil
	// 检查调度器是否为空
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	// Check if job is nil
	// 检查任务是否为空
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	return w.scheduler.AddWeeklyJob(w.job)
}

// BatchAdd adds multiple weekly jobs to the scheduler.
// 批量添加多个周任务到调度器。
//
// Parameters:
//
//	weeklyJobs - Variable number of weekly jobs to add / 可变数量的要添加的周任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error         - Error if jobs cannot be added / 如果无法添加任务的错误
func (w *WeeklyJobClient) BatchAdd(weeklyJobs ...*WeeklyJob) ([]gocron.Job, error) {
	// Check if scheduler is nil
	// 检查调度器是否为空
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	// Check if job is nil
	// 检查任务是否为空
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	return w.scheduler.AddWeeklyJobs(weeklyJobs...)
}

// Remove removes the weekly job from the scheduler.
// 从调度器中删除周任务。
//
// Returns:
//
//	error - Error if the job cannot be removed / 如果无法删除任务的错误
func (w *WeeklyJobClient) Remove() error {
	// Check if scheduler is nil
	// 检查调度器是否为空
	if w.scheduler == nil {
		return ErrScheduleNil
	}
	// Check if job is nil
	// 检查任务是否为空
	if w.job == nil {
		return ErrWeeklyJobNil
	}
	// Remove by name if ID is empty
	// 如果ID为空，则按名称删除
	if w.job.ID == "" {
		return w.scheduler.RemoveJob(w.job.Name)
	}
	// Remove by alias if alias is set
	// 如果设置了别名，则按别名删除
	if w.job.Ali != "" {
		return w.scheduler.RemoveJobByAlias(w.job.Ali)
	}
	// Remove by name if name is set
	// 如果设置了名称，则按名称删除
	if w.job.Name != "" {
		return w.scheduler.RemoveJobByName(w.job.ID)
	}
	return ErrJobNotFound
}

// Get retrieves the weekly job from the scheduler.
// 从调度器中获取周任务。
//
// Returns:
//
//	gocron.Job - The retrieved job / 获取的任务
//	error      - Error if the job cannot be found / 如果找不到任务的错误
func (w *WeeklyJobClient) Get() (gocron.Job, error) {
	// Check if scheduler is nil
	// 检查调度器是否为空
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	// Check if job is nil
	// 检查任务是否为空
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	// Get by ID if ID is set
	// 如果设置了ID，则按ID获取
	if w.job.ID != "" {
		return w.scheduler.GetJobByID(w.job.ID)
	}
	// Get by alias if alias is set
	// 如果设置了别名，则按别名获取
	if w.job.Ali != "" {
		return w.scheduler.GetJobByAlias(w.job.Ali)
	}
	// Get by name if name is set
	// 如果设置了名称，则按名称获取
	if w.job.Name != "" {
		return w.scheduler.GetJobByName(w.job.Name)
	}
	return nil, ErrJobNotFound
}

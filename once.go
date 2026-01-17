package chrono

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// OnceJob represents a job that runs only once at specified times.
// OnceJob 表示一个在指定时间只运行一次的任务。
type OnceJob struct {
	// ID is the unique identifier for the job.
	// ID 是任务的唯一标识符。
	ID string
	// Type is the job type.
	// Type 是任务的类型。
	Type JobType
	// Ali is the alias for the job.
	// Ali 是任务的别名。
	Ali string
	// Name is the name of the job.
	// Name 是任务的名称。
	Name string
	// WorkTime is the specific times when the job should run.
	// WorkTime 是任务应该运行的特定时间。
	WorkTime []time.Time
	// Tag are the tags for the job.
	// Tag 是任务的标签。
	Tag []string
	// TaskFunc is the function to execute as the job.
	// TaskFunc 是作为任务执行的函数。
	TaskFunc any
	// Parameters are the parameters to pass to the task function.
	// Parameters 是传递给任务函数的参数。
	Parameters []any
	// Hooks are the event hooks for job lifecycle events.
	// Hooks 是任务生命周期事件的钩子。
	Hooks []gocron.EventListener
	// WatchFunc is the function to watch job events.
	// WatchFunc 是监听任务事件的函数。
	WatchFunc func(event JobWatchInterface)
	// jobOptions holds job-level options.
	// jobOptions 保存任务级别的选项。
	jobOptions *jobOptions
	// err is the error state for the job.
	// err 是任务的错误状态。
	err error
}

// NewOnceJob creates a new OnceJob with specified work times.
// NewOnceJob 创建一个具有指定工作时间的 OnceJob。
//
// Parameters:
//
//	workTimes - Variable number of time values when job should run / 任务应该运行的可变数量时间值
//
// Returns:
//
//	*OnceJob - The newly created one-time job / 新创建的一次性任务
func NewOnceJob(workTimes ...time.Time) *OnceJob {
	return &OnceJob{
		WorkTime: workTimes,
		Type:     JobTypeOnce,
	}
}

// AtTimes sets the specific times when the one-time job should run.
// AtTimes 设置一次性任务运行的特定时间。
//
// Parameters:
//
//	workTimes - Variable number of time values / 可变数量的时间值
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) AtTimes(workTimes ...time.Time) *OnceJob {
	c.WorkTime = workTimes
	return nil
}

// Alias sets the alias for the OnceJob.
// Alias 设置 OnceJob 的别名。
//
// Parameters:
//
//	alias - The alias string for the job / 任务的别名字符串
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) Alias(alias string) *OnceJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the OnceJob.
// JobID 设置 OnceJob 的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) JobID(id string) *OnceJob {
	c.ID = id
	return c
}

// Names sets the name for the OnceJob. If name is empty, a UUID is generated.
// Names 设置 OnceJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The name string for the job / 任务的名称字符串
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) Names(name string) *OnceJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tags sets the tags for the OnceJob.
// Tags 设置 OnceJob 的标签。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) Tags(tags ...string) *OnceJob {
	c.Tag = tags
	return c
}

// Task sets the task function and its parameters for the OnceJob.
// It wraps the task with error and timeout handling.
// Task 设置 OnceJob 的任务函数及其参数，并包装错误和超时处理。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) Task(task any, parameters ...any) *OnceJob {
	if task == nil {
		c.err = errors.Join(c.err, ErrTaskFuncNil)
		return c
	}
	c.Parameters = append(c.Parameters, parameters)
	c.TaskFunc = func() error {
		return callJobFunc(task, c.Parameters...)
	}
	return c
}

// Watch sets a watcher function for job events.
// Watch 设置任务事件的监听函数。
//
// Parameters:
//
//	watch - The watch function for job events / 任务事件的监听函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) Watch(watch func(event JobWatchInterface)) *OnceJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the OnceJob.
// addHooks 向 OnceJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners / 可变数量的事件监听器
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) addHooks(hook ...gocron.EventListener) *OnceJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the OnceJob.
// DefaultHooks 向 OnceJob 添加一组默认事件监听器。
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) DefaultHooks() *OnceJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns adds a hook to be called before the job runs.
// BeforeJobRuns 添加一个在任务运行前调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the job runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
//
// Parameters:
//
//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *OnceJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the job runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the job runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the job panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *OnceJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during job execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with job ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	*OnceJob - The one-time job for method chaining / 支持链式调用的一次性任务
func (c *OnceJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

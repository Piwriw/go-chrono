package chrono

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// IntervalJob represents a job that runs at a fixed interval.
// IntervalJob 表示一个按固定时间间隔运行的任务。
type IntervalJob struct {
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
	// Interval is the time duration between job runs.
	// Interval 是任务运行之间的时间间隔。
	Interval time.Duration
	// TaskFunc is the function to execute as the job.
	// TaskFunc 是要作为任务执行的函数。
	TaskFunc any
	// Tags are the tags for the job.
	// Tags 是任务的标签。
	Tags []string
	// Parameters are the parameters to pass to the task function.
	// Parameters 是要传递给任务函数的参数。
	Parameters []any
	// Hooks are the event hooks for job lifecycle events.
	// Hooks 是任务生命周期事件的事件钩子。
	Hooks []gocron.EventListener
	// WatchFunc is the function to watch job events.
	// WatchFunc 是监听任务事件的函数。
	WatchFunc func(event JobWatchInterface)
	// err is the error state for the job.
	// err 是任务的错误状态。
	err error
}

// NewIntervalJob creates a new IntervalJob with the specified interval.
// NewIntervalJob 创建一个具有指定时间间隔的 IntervalJob。
//
// Parameters:
//
//	interval - The time duration between job executions / 任务执行之间的时间间隔
//
// Returns:
//
//	*IntervalJob - The newly created interval job / 新创建的间隔任务
func NewIntervalJob(interval time.Duration) *IntervalJob {
	return &IntervalJob{
		Interval: interval,
		Type:     JobInterval,
	}
}

// IntervalTime sets the time interval for the IntervalJob.
// IntervalTime 设置 IntervalJob 的时间间隔。
//
// Parameters:
//
//	interval - The time duration between job executions / 任务执行之间的时间间隔
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) IntervalTime(interval time.Duration) *IntervalJob {
	c.Interval = interval
	return c
}

// Alias sets the alias for the IntervalJob.
// Alias 设置 IntervalJob 的别名。
//
// Parameters:
//
//	alias - The alias string for the job / 任务的别名字符串
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Alias(alias string) *IntervalJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the IntervalJob.
// JobID 设置 IntervalJob 的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) JobID(id string) *IntervalJob {
	c.ID = id
	return c
}

// Names sets the name for the IntervalJob. If name is empty, a UUID is generated.
// Names 设置 IntervalJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The name string for the job / 任务的名称字符串
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Names(name string) *IntervalJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tag adds tags to the IntervalJob.
// Tag 添加标签到 IntervalJob。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Tag(tags ...string) *IntervalJob {
	c.Tags = tags
	return c
}

// Task sets the task function and its parameters for the IntervalJob.
// It wraps the task with error and timeout handling.
// Task 设置 IntervalJob 的任务函数及其参数，并包装错误和超时处理。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Task(task any, parameters ...any) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Watch(watch func(event JobWatchInterface)) *IntervalJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the IntervalJob.
// addHooks 向 IntervalJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners / 可变数量的事件监听器
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) addHooks(hook ...gocron.EventListener) *IntervalJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the IntervalJob.
// DefaultHooks 向 IntervalJob 添加一组默认事件监听器。
//
// Returns:
//
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) DefaultHooks() *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
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
//	*IntervalJob - The interval job for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

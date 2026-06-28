package jobs

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/pkg/executor"
)

// IntervalJob represents a jobs that runs at a fixed interval.
// IntervalJob 表示一个按固定时间间隔运行的任务。
type IntervalJob struct {
	// ID is the unique identifier for the jobs.
	// ID 是任务的唯一标识符。
	ID string
	// Type is the job type.
	// Type 是任务的类型。
	Type JobType
	// Ali is the alias for the jobs.
	// Ali 是任务的别名。
	Ali string
	// Name is the name of the jobs.
	// Name 是任务的名称。
	Name string
	// Interval is the time duration between jobs runs.
	// Interval 是任务运行之间的时间间隔。
	Interval time.Duration
	// TaskFunc is the function to execute as the jobs.
	// TaskFunc 是要作为任务执行的函数。
	TaskFunc any
	// Tags are the tags for the jobs.
	// Tags 是任务的标签。
	Tags []string
	// Parameters are the parameters to pass to the task function.
	// Parameters 是要传递给任务函数的参数。
	Parameters []any
	// Hooks are the event hooks for jobs lifecycle events.
	// Hooks 是任务生命周期事件的事件钩子。
	Hooks []gocron.EventListener
	// WatchFunc is the function to watch jobs events.
	// WatchFunc 是监听任务事件的函数。
	WatchFunc func(event common.JobWatchInterface)
	// jobOptions holds jobs-level options.
	// jobOptions 保存任务级别的选项。
	JobOptions *common.JobOptions
	// err is the error state for the jobs.
	// err 是任务的错误状态。
	err error
}

// NewIntervalJob creates a new IntervalJob with the specified interval.
// NewIntervalJob 创建一个具有指定时间间隔的 IntervalJob。
//
// Parameters:
//
//	interval - The time duration between jobs executions / 任务执行之间的时间间隔
//
// Returns:
//
//	*IntervalJob - The newly created interval jobs / 新创建的间隔任务
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
//	interval - The time duration between jobs executions / 任务执行之间的时间间隔
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) IntervalTime(interval time.Duration) *IntervalJob {
	c.Interval = interval
	return c
}

// Alias sets the alias for the IntervalJob.
// Alias 设置 IntervalJob 的别名。
//
// Parameters:
//
//	alias - The alias string for the jobs / 任务的别名字符串
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) JobID(id string) *IntervalJob {
	c.ID = id
	return c
}

// Names sets the name for the IntervalJob. If name is empty, a UUID is generated.
// Names 设置 IntervalJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The name string for the jobs / 任务的名称字符串
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Task(task any, parameters ...any) *IntervalJob {
	if task == nil {
		c.err = errors.Join(c.err, common.ErrTaskFuncNil)
		return c
	}
	c.Parameters = append(c.Parameters, parameters)
	c.TaskFunc = func() error {
		return executor.CallJobFunc(task, c.Parameters...)
	}
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) Watch(watch func(event common.JobWatchInterface)) *IntervalJob {
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
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
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) DefaultHooks() *IntervalJob {
	return c.addHooks(
		gocron.BeforeJobRuns(common.DefaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(common.DefaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(common.DefaultAfterJobRuns),
		gocron.AfterJobRunsWithError(common.DefaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(common.DefaultAfterJobRunsWithPanic),
		gocron.AfterLockError(common.DefaultAfterLockError))
}

// BeforeJobRuns adds a hook to be called before the jobs runs.
// BeforeJobRuns 添加一个在任务运行前调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the jobs runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
//
// Parameters:
//
//	eventListenerFunc - The callback function that returns error / 返回错误的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the jobs runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID and name / 包含任务ID和名称的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the jobs runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the jobs panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and recovered data / 包含任务ID、名称和恢复数据的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during jobs execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
//
// Parameters:
//
//	eventListenerFunc - The callback function with jobs ID, name, and error / 包含任务ID、名称和错误的回调函数
//
// Returns:
//
//	*IntervalJob - The interval jobs for method chaining / 支持链式调用的间隔任务
func (c *IntervalJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

// GetError returns the configuration error for the jobs.
// GetError 返回任务的配置错误。
//
// Returns:
//
//	error - The configuration error, or nil if no error / 配置错误，如果没有错误则返回 nil
func (j *IntervalJob) GetError() error {
	return j.err
}

package chrono

import (
	"errors"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// MonthJob represents a job that runs on a monthly schedule.
// MonthJob 表示一个按月调度运行的任务。
type MonthJob struct {
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
	// Interval is the interval in months between job runs.
	// Interval 是任务运行的月数间隔。
	Interval uint
	// DaysOfTheMonth are the days of the month to run the job.
	// DaysOfTheMonth 是任务每月运行的具体日期。
	DaysOfTheMonth gocron.DaysOfTheMonth
	// AtTimes are the specific times of day to run the job.
	// AtTimes 是任务每天运行的具体时间点。
	AtTimes gocron.AtTimes
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
	// err is the error state for the job.
	// err 是任务的错误状态。
	err error
}

// NewMonthJob creates a new MonthJob with the specified interval, days, and time.
// NewMonthJob 创建一个具有指定间隔、日期和时间的 MonthJob。
//
// Parameters:
//
//	interval - The interval in months between job runs / 任务运行之间的月数间隔
//	days    - The days of the month to run the job / 任务运行的月份日期
//	atTime  - The specific times to run the job / 任务运行的具体时间
//
// Returns:
//
//	*MonthJob - The newly created monthly job / 新创建的月度任务
func NewMonthJob(interval uint, days gocron.DaysOfTheMonth, atTime gocron.AtTimes) *MonthJob {
	return &MonthJob{
		Interval:       interval,
		DaysOfTheMonth: days,
		AtTimes:        atTime,
		Type:           JobTypeMonthly,
	}
}

// NewMonthJobAtTime creates a new MonthJob that runs at specific days and times every month.
// NewMonthJobAtTime 创建一个每月在特定日期和时间运行的 MonthJob。
//
// Parameters:
//
//	days   - The days of the month to run the job / 任务运行的月份日期
//	hour   - The hour of the day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	*MonthJob - The newly created monthly job / 新创建的月度任务
func NewMonthJobAtTime(days []int, hour, minute, second int) *MonthJob {
	if len(days) == 0 {
		return &MonthJob{
			err: ErrAtTimeDaysNil,
		}
	}
	return &MonthJob{
		Interval:       1,
		DaysOfTheMonth: gocron.NewDaysOfTheMonth(days[0], days[1:]...),
		AtTimes:        gocron.NewAtTimes(gocron.NewAtTime(uint(hour), uint(minute), uint(second))),
	}
}

// AtTime sets the time of day for the MonthJob to run.
// AtTime 设置 MonthJob 运行的时间。
//
// Parameters:
//
//	days   - The days of the month to run the job / 任务运行的月份日期
//	hour   - The hour of the day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) AtTime(days []int, hour, minute, second int) *MonthJob {
	if len(days) == 0 {
		return &MonthJob{
			err: ErrAtTimeDaysNil,
		}
	}
	return &MonthJob{
		Interval:       1,
		DaysOfTheMonth: gocron.NewDaysOfTheMonth(days[0], days[1:]...),
		AtTimes:        gocron.NewAtTimes(gocron.NewAtTime(uint(hour), uint(minute), uint(second))),
	}
}

// Alias sets the alias for the MonthJob.
// Alias 设置 MonthJob 的别名。
//
// Parameters:
//
//	alias - The alias string for the job / 任务的别名字符串
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) Alias(alias string) *MonthJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the MonthJob.
// JobID 设置 MonthJob 的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) JobID(id string) *MonthJob {
	c.ID = id
	return c
}

// Names sets the name for the MonthJob. If name is empty, a UUID is generated.
// Names 设置 MonthJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The name string for the job / 任务的名称字符串
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) Names(name string) *MonthJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tags sets the tags for the MonthJob.
// Tags 设置 MonthJob 的标签。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) Tags(tags ...string) *MonthJob {
	c.Tag = tags
	return c
}

// Task sets the task function and its parameters for the MonthJob.
// It wraps the task with error and timeout handling.
// Task 设置 MonthJob 的任务函数及其参数，并包装错误和超时处理。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) Task(task any, parameters ...any) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) Watch(watch func(event JobWatchInterface)) *MonthJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the MonthJob.
// addHooks 向 MonthJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners / 可变数量的事件监听器
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) addHooks(hook ...gocron.EventListener) *MonthJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the MonthJob.
// DefaultHooks 向 MonthJob 添加一组默认事件监听器。
//
// Returns:
//
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) DefaultHooks() *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *MonthJob {
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
//	*MonthJob - The monthly job for method chaining / 支持链式调用的月度任务
func (c *MonthJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

package jobs

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/pkg/executor"
)

// WeeklyJob represents a jobs that runs on a weekly schedule.
// WeeklyJob 表示一个按周调度运行的任务。
type WeeklyJob struct {
	// Unique identifier for the jobs
	// 任务的唯一标识符
	ID string
	// 任务的类型
	// JobType of the jobs
	Type JobType
	// Alias for the jobs
	// 任务的别名
	Ali string
	// Name of the jobs
	// 任务名称
	Name string
	// Interval in weeks between jobs runs
	// 任务运行的周数间隔
	Interval uint
	// Days of the week to run the jobs
	// 任务每周运行的具体星期几
	DaysOfTheWeek gocron.Weekdays
	// Specific times of day to run the jobs
	// 任务每天运行的具体时间点
	WorkTimes gocron.AtTimes
	// Tag for the jobs
	// 任务的标签
	Tag []string
	// The function to execute as the jobs
	// 作为任务执行的函数
	TaskFunc any
	// Parameters to pass to the task function
	// 传递给任务函数的参数
	Parameters []any
	// Event hooks for jobs lifecycle events
	// 任务生命周期事件的钩子
	Hooks []gocron.EventListener
	// Function to watch jobs events
	// 监听任务事件的函数
	WatchFunc func(event common.JobWatchInterface)
	// jobOptions holds jobs-level options
	// jobOptions 保存任务级别的选项
	JobOptions *common.JobOptions
	// Error state for the jobs
	// 任务的错误状态
	err error
}

// NewWeeklyJob creates a new WeeklyJob with the specified interval, days, and time.
// NewWeeklyJob 创建一个具有指定间隔、星期和时间的 WeeklyJob。
//
// Parameters:
//
//	interval - The interval in weeks between jobs runs / 任务运行的周数间隔
//	days     - The days of the week when the jobs should run / 任务应该运行的星期几
//	atTime   - The specific times of day to run the jobs / 任务每天运行的具体时间点
//
// Returns:
//
//	*WeeklyJob - The created weekly jobs / 创建的周任务
func NewWeeklyJob(interval uint, days gocron.Weekdays, atTime gocron.AtTimes) *WeeklyJob {
	return &WeeklyJob{
		Interval:      interval,
		DaysOfTheWeek: days,
		WorkTimes:     atTime,
		Type:          JobTypeWeekly,
	}
}

// NewWeeklyJobAtTime creates a new WeeklyJob that runs at specific days and times every week.
// NewWeeklyJobAtTime 创建一个每周在特定星期和时间运行的 WeeklyJob。
//
// Parameters:
//
//	days   - The days of the week when the jobs should run / 任务应该运行的星期几
//	hour   - The hour of day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	*WeeklyJob - The created weekly jobs / 创建的周任务
func NewWeeklyJobAtTime(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	// Check if days slice is empty
	// 检查星期数组是否为空
	if len(days) == 0 {
		return &WeeklyJob{
			err: common.ErrAtTimeDaysNil,
		}
	}
	return &WeeklyJob{
		Interval:      1,
		DaysOfTheWeek: gocron.NewWeekdays(days[0], days[1:]...),
		WorkTimes:     gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
		Type:          JobTypeWeekly,
	}
}

// AtTimes sets specific days and time for the weekly jobs to run.
// AtTimes 设置周任务运行的特定星期和时间。
//
// Parameters:
//
//	days   - The days of the week when the jobs should run / 任务应该运行的星期几
//	hour   - The hour of day (0-23) / 一天中的小时（0-23）
//	minute - The minute of the hour (0-59) / 小时中的分钟（0-59）
//	second - The second of the minute (0-59) / 分钟中的秒（0-59）
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) AtTimes(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	// Check if days slice is empty
	// 检查星期数组是否为空
	if len(days) == 0 {
		return &WeeklyJob{
			err: common.ErrAtTimeDaysNil,
		}
	}
	return &WeeklyJob{
		Interval:      1,
		DaysOfTheWeek: gocron.NewWeekdays(days[0], days[1:]...),
		WorkTimes:     gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
		Type:          JobTypeWeekly,
	}
}

// Alias sets the alias for the WeeklyJob.
// Alias 设置 WeeklyJob 的别名。
//
// Parameters:
//
//	alias - The alias string / 别名字符串
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) Alias(alias string) *WeeklyJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the WeeklyJob.
// JobID 设置 WeeklyJob 的唯一标识符。
//
// Parameters:
//
//	id - The jobs ID string / 任务ID字符串
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) JobID(id string) *WeeklyJob {
	c.ID = id
	return c
}

// Names sets the name for the WeeklyJob. If name is empty, a UUID is generated.
// Names 设置 WeeklyJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The jobs name string / 任务名称字符串
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) Names(name string) *WeeklyJob {
	// Generate UUID if name is empty
	// 如果名称为空，则生成UUID
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tags sets the tags for the WeeklyJob.
// Tags 设置 WeeklyJob 的标签。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) Tags(tags ...string) *WeeklyJob {
	c.Tag = tags
	return c
}

// Task sets the task function and its parameters for the WeeklyJob.
// It wraps the task with error and timeout handling.
// Task 设置 WeeklyJob 的任务函数及其参数，并包装错误和超时处理。
//
// Parameters:
//
//	task      - The task function / 任务函数
//	parameters - Variable parameters for the task / 任务的可变参数
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) Task(task any, parameters ...any) *WeeklyJob {
	// Check if task function is nil
	// 检查任务函数是否为空
	if task == nil {
		c.err = errors.Join(c.err, common.ErrTaskFuncNil)
		return c
	}
	// Append parameters to the existing parameters
	// 将参数追加到现有参数列表
	c.Parameters = append(c.Parameters, parameters)
	// Wrap task function with error handling
	// 包装任务函数以进行错误处理
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
//	watch - The watch function for jobs events / 任务事件的监视函数
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) Watch(watch func(event common.JobWatchInterface)) *WeeklyJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the WeeklyJob.
// addHooks 向 WeeklyJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners to add / 可变数量的要添加的事件监听器
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) addHooks(hook ...gocron.EventListener) *WeeklyJob {
	// Initialize hooks slice if nil
	// 如果钩子切片为空，则初始化
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	// Append hooks to the existing hooks
	// 将钩子追加到现有钩子列表
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the WeeklyJob.
// DefaultHooks 向 WeeklyJob 添加一组默认事件监听器。
//
// Returns:
//
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) DefaultHooks() *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob {
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
//	*WeeklyJob - The weekly jobs for method chaining / 周任务，支持链式调用
func (c *WeeklyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

// GetError returns the configuration error for the jobs.
// GetError 返回任务的配置错误。
//
// Returns:
//
//	error - The configuration error, or nil if no error / 配置错误，如果没有错误则返回 nil
func (j *WeeklyJob) GetError() error {
	return j.err
}

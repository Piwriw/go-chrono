package jobs

import (
	"errors"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/pkg/executor"
)

// DailyJob represents a jobs that runs on a daily schedule.
// DailyJob 表示一个按天调度运行的任务。
type DailyJob struct {
	// ID is the unique identifier for the jobs / 任务的唯一标识符
	ID string
	// Type is the job type / 任务的类型
	Type JobType
	// Ali is the alias for the jobs / 任务的别名
	Ali string
	// Name is the name of the jobs / 任务名称
	Name string
	// Interval is the number of days between jobs runs / 任务运行的天数间隔
	Interval uint
	// AtTimes are the specific times of day to run the jobs / 任务每天运行的具体时间点
	AtTimes gocron.AtTimes
	// Tags are labels for categorization / 标签用于分类
	Tags []string
	// TaskFunc is the function to execute as the jobs / 作为任务执行的函数
	TaskFunc any
	// Parameters are values to pass to the task function / 传递给任务函数的参数
	Parameters []any
	// Hooks are event listeners for jobs lifecycle events / 任务生命周期事件的钩子
	Hooks []gocron.EventListener
	// WatchFunc monitors jobs events / 监听任务事件的函数
	WatchFunc func(event common.JobWatchInterface)
	// jobOptions holds jobs-level options / 保存任务级别的选项
	JobOptions *common.JobOptions
	// err holds any configuration errors / 存储配置错误
	err error
}

// NewDailyJob creates a new DailyJob with specified interval and times.
// NewDailyJob 创建一个具有指定间隔和时间的 DailyJob。
//
// Parameters:
//
//	interval - The number of days between runs / 运行之间的天数
//	atTime   - The specific times to run / 运行的具体时间
//
// Returns:
//
//	*DailyJob - The newly created DailyJob / 新创建的 DailyJob
func NewDailyJob(interval uint, atTime gocron.AtTimes) *DailyJob {
	return &DailyJob{
		Interval: interval,
		AtTimes:  atTime,
		Type:     JobTypeDaily,
	}
}

// NewDailyJobAtTime creates a new DailyJob that runs at a specific time every day.
// NewDailyJobAtTime 创建一个每天在特定时间运行的 DailyJob。
//
// Parameters:
//
//	hour   - The hour (0-23) / 小时（0-23）
//	minute - The minute (0-59) / 分钟（0-59）
//	second - The second (0-59) / 秒（0-59）
//
// Returns:
//
//	*DailyJob - The newly created DailyJob / 新创建的 DailyJob
func NewDailyJobAtTime(hour, minute, second uint) *DailyJob {
	return &DailyJob{
		Interval: 1,
		AtTimes:  gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
		Type:     JobTypeDaily,
	}
}

// AtDayTime sets the specific time of day to run the jobs every day.
// AtDayTime 设置任务每天运行的具体时间。
//
// Parameters:
//
//	hour   - The hour (0-23) / 小时（0-23）
//	minute - The minute (0-59) / 分钟（0-59）
//	second - The second (0-59) / 秒（0-59）
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) AtDayTime(hour, minute, second uint) *DailyJob {
	c.AtTimes = gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second))
	return c
}

// Alias sets the alias for the DailyJob.
// Alias 设置 DailyJob 的别名。
//
// Parameters:
//
//	alias - The alias string / 别名字符串
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) Alias(alias string) *DailyJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the DailyJob.
// JobID 设置 DailyJob 的唯一标识符。
//
// Parameters:
//
//	id - The jobs ID string / 任务ID字符串
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) JobID(id string) *DailyJob {
	c.ID = id
	return c
}

// Names sets the name for the DailyJob. If name is empty, a UUID is generated.
// Names 设置 DailyJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The jobs name string / 任务名称字符串
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) Names(name string) *DailyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tag adds tags to the DailyJob for categorization.
// Tag 为 DailyJob 添加标签以便分类。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) Tag(tags ...string) *DailyJob {
	c.Tags = tags
	return c
}

// Task sets the task function and its parameters for the DailyJob.
// Task 设置 DailyJob 的任务函数及其参数，并包装错误和超时处理。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters to pass to the task / 传递给任务的可变参数
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) Task(task any, parameters ...any) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) Watch(watch func(event common.JobWatchInterface)) *DailyJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the DailyJob.
// addHooks 向 DailyJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners / 可变数量的事件监听器
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) addHooks(hook ...gocron.EventListener) *DailyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the DailyJob.
// DefaultHooks 向 DailyJob 添加一组默认事件监听器。
//
// Returns:
//
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) DefaultHooks() *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *DailyJob {
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
//	*DailyJob - The DailyJob for method chaining / 支持链式调用的 DailyJob
func (c *DailyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

// GetError returns the configuration error for the jobs.
// GetError 返回任务的配置错误。
//
// Returns:
//
//	error - The configuration error, or nil if no error / 配置错误，如果没有错误则返回 nil
func (j *DailyJob) GetError() error {
	return j.err
}

package chrono

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// CronJob represents a job that runs based on a cron expression.
// CronJob 表示一个基于 cron 表达式运行的任务。
type CronJob struct {
	// ID is the unique identifier for the job / 任务的唯一标识符
	ID string
	// Type is the job type / 任务的类型
	Type JobType
	// Ali is the alias for the job / 任务的别名
	Ali string
	// Name is the name of the job / 任务名称
	Name string
	// Expr is the cron expression for scheduling / 用于调度的 cron 表达式
	Expr string
	// Tags are labels for categorization / 标签用于分类
	Tags []string
	// TaskFunc is the function to execute as the job / 作为任务执行的函数
	TaskFunc any
	// Parameters are values to pass to the task function / 传递给任务函数的参数
	Parameters []any
	// Hooks are event listeners for job lifecycle events / 任务生命周期事件的钩子
	Hooks []gocron.EventListener
	// WatchFunc monitors job events / 监听任务事件的函数
	WatchFunc func(event JobWatchInterface)
	// err holds any configuration errors / 存储配置错误
	err error
}

// NewCronJob creates a new CronJob with the specified cron expression.
// NewCronJob 创建一个具有指定 cron 表达式的 CronJob。
//
// Parameters:
//
//	expr - The cron expression string / Cron 表达式字符串
//
// Returns:
//
//	*CronJob - The newly created CronJob / 新创建的 CronJob
func NewCronJob(expr string) *CronJob {
	return &CronJob{
		Expr: expr,
		Type: JobTypeCron,
	}
}

// CronExpr sets the cron expression for the CronJob.
// CronExpr 设置 CronJob 的 cron 表达式。
//
// Parameters:
//
//	expr - The cron expression string / Cron 表达式字符串
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) CronExpr(expr string) *CronJob {
	// Set cron expression which will be validated during scheduling
	// 设置 cron 表达式，在调度中将会校验 cron 表达式的有效性
	c.Expr = expr
	return c
}

// Alias sets the alias for the CronJob.
// Alias 设置 CronJob 的别名。
//
// Parameters:
//
//	alias - The alias string / 别名字符串
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) Alias(alias string) *CronJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the CronJob.
// JobID 设置 CronJob 的唯一标识符。
//
// Parameters:
//
//	id - The unique identifier string / 唯一标识符字符串
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) JobID(id string) *CronJob {
	c.ID = id
	return c
}

// Names sets the name for the CronJob. If name is empty, a UUID is generated.
// Names 设置 CronJob 的名称。如果名称为空，则生成一个 UUID。
//
// Parameters:
//
//	name - The job name string / 任务名称字符串
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) Names(name string) *CronJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tag adds tags to the CronJob for categorization.
// Tag 为 CronJob 添加标签以便分类。
//
// Parameters:
//
//	tags - Variable number of tag strings / 可变数量的标签字符串
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) Tag(tags ...string) *CronJob {
	c.Tags = tags
	return c
}

// Task sets the task function and its parameters for the CronJob.
// Task 设置 CronJob 的任务函数及其参数。
//
// Parameters:
//
//	task      - The task function to execute / 要执行的任务函数
//	parameters - Variable parameters to pass to the task / 传递给任务的可变参数
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) Task(task any, parameters ...any) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) Watch(watch func(event JobWatchInterface)) *CronJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the CronJob.
// addHooks 向 CronJob 添加一个或多个事件监听器（钩子）。
//
// Parameters:
//
//	hook - Variable number of event listeners / 可变数量的事件监听器
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) addHooks(hook ...gocron.EventListener) *CronJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the CronJob.
// DefaultHooks 向 CronJob 添加一组默认事件监听器。
//
// Returns:
//
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) DefaultHooks() *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
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
//	*CronJob - The CronJob for method chaining / 支持链式调用的 CronJob
func (c *CronJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

// TimeType represents the type of time-based scheduling.
// TimeType 表示基于时间的调度类型。

// DayTimeType is the cron format for daily jobs.
// DayTimeType 是每日任务的 cron 格式。
const (
	DayTimeType   string = "%d %d * * *"
	WeekTimeType  string = "%d %d * * %d"
	MonthTimeType string = "%d %d * %d *"
)

// DayTimeToCron converts a time.Time to a daily cron expression.
// DayTimeToCron 将 time.Time 转换为每日的 Cron 表达式。
//
// Parameters:
//
//	t - The time to convert / 要转换的时间
//
// Returns:
//
//	string - The cron expression string / Cron 表达式字符串
func DayTimeToCron(t time.Time) string {
	// Extract time fields
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// Return cron expression
	// 返回 Cron 表达式
	return fmt.Sprintf(DayTimeType, minute, hour)
}

// WeekTimeToCron converts a time.Time and weekday to a weekly cron expression.
// WeekTimeToCron 将 time.Time 和星期几转换为每周的 Cron 表达式。
//
// Parameters:
//
//	t    - The time to convert / 要转换的时间
//	week - The weekday to schedule / 要调度的星期几
//
// Returns:
//
//	string - The cron expression string / Cron 表达式字符串
func WeekTimeToCron(t time.Time, week time.Weekday) string {
	// Extract time fields
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// Return cron expression
	// 返回 Cron 表达式
	return fmt.Sprintf(WeekTimeType, minute, hour, week)
}

// MonthTimeToCron converts a time.Time and month to a monthly cron expression.
// MonthTimeToCron 将 time.Time 和月份转换为每月的 Cron 表达式。
//
// Parameters:
//
//	t     - The time to convert / 要转换的时间
//	month - The month to schedule / 要调度的月份
//
// Returns:
//
//	string - The cron expression string / Cron 表达式字符串
func MonthTimeToCron(t time.Time, month time.Month) string {
	// Extract time fields
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// Return cron expression
	// 返回 Cron 表达式
	return fmt.Sprintf(MonthTimeType, minute, hour, month)
}

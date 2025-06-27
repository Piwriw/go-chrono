package chrono

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// CronJob represents a job that runs based on a cron expression.
// CronJob 表示一个基于 cron 表达式运行的任务。
type CronJob struct {
	// 任务的唯一标识符
	// Unique identifier for the job
	ID string
	// 任务的类型
	// JobType of the job
	Type JobType
	// Alias for the job
	// 任务的别名
	Ali string
	// Name of the job
	// 任务名称
	Name string
	// Cron expression for scheduling
	// 用于调度的 cron 表达式
	Expr string
	// The function to execute as the job
	// 作为任务执行的函数
	Tags     []string
	TaskFunc any
	// Parameters to pass to the task function
	// 传递给任务函数的参数
	Parameters []any
	// Event hooks for job lifecycle events
	// 任务生命周期事件的钩子
	Hooks []gocron.EventListener
	// Function to watch job events
	// 监听任务事件的函数
	WatchFunc func(event JobWatchInterface)
	// Timeout for the job execution
	// 任务执行的超时时间
	timeout time.Duration
	// Timeout for the job execution
	// 任务执行的超时时间
	err error
}

// NewCronJob creates a new CronJob with the specified cron expression.
// NewCronJob 创建一个具有指定 cron 表达式的 CronJob。
func NewCronJob(expr string) *CronJob {
	return &CronJob{
		Expr: expr,
		Type: JobTypeCron,
	}
}

// CronExpr sets the cron expression for the CronJob.
// CronExpr 设置 cron 表达式。
func (c *CronJob) CronExpr(expr string) *CronJob {
	c.Expr = expr
	return c
}

// Alias sets the alias for the CronJob.
// Alias 设置 CronJob 的别名。
func (c *CronJob) Alias(alias string) *CronJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the CronJob.
// JobID 设置 CronJob 的唯一标识符。
func (c *CronJob) JobID(id string) *CronJob {
	c.ID = id
	return c
}

// Names sets the name for the CronJob. If name is empty, a UUID is generated.
// Names 设置 CronJob 的名称。如果名称为空，则生成一个 UUID。
func (c *CronJob) Names(name string) *CronJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *CronJob) Tag(tags ...string) *CronJob {
	c.Tags = tags
	return c
}

// Task sets the task function and its parameters for the CronJob.
// Task 设置 CronJob 的任务函数及其参数。
func (c *CronJob) Task(task any, parameters ...any) *CronJob {
	if task == nil {
		c.err = errors.Join(c.err, ErrTaskFuncNil)
		return c
	}
	c.TaskFunc = func() error {
		var ctx context.Context
		var cancel context.CancelFunc
		// 如果设置了超时时间，则使用 context.WithTimeout
		if c.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
			defer cancel()
		} else {
			// 如果没有设置超时时间，则直接使用背景上下文
			ctx = context.Background()
		}

		done := make(chan error, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					done <- fmt.Errorf("chrono:task panicked: %v", r)
				}
			}()
			done <- callJobFunc(task, parameters...)
		}()

		select {
		case err := <-done:
			if err != nil {
				slog.Error("chrono:task exec failed", "err", err)
				return ErrTaskFailed
			}
		case <-ctx.Done():
			return ErrTaskTimeout
		}
		return nil
	}
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

// Timeout sets the timeout duration for the CronJob execution.
// Timeout 设置 CronJob 执行的超时时间。
func (c *CronJob) Timeout(timeout time.Duration) *CronJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

// Watch sets a watcher function for job events.
// Watch 设置任务事件的监听函数。
func (c *CronJob) Watch(watch func(event JobWatchInterface)) *CronJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the CronJob.
// addHooks 向 CronJob 添加一个或多个事件监听器（钩子）。
func (c *CronJob) addHooks(hook ...gocron.EventListener) *CronJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the CronJob.
// DefaultHooks 向 CronJob 添加一组默认事件监听器。
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
func (c *CronJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the job runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
func (c *CronJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the job runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
func (c *CronJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the job runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
func (c *CronJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the job panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
func (c *CronJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during job execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
func (c *CronJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

type TimeType string

// DayTimeType is the cron format for daily jobs.
// DayTimeType 是每日任务的 cron 格式。
const (
	DayTimeType   string = "%d %d * * * "
	WeekTimeType  string = "%d %d * * %d "
	MonthTimeType string = "%d %d * %d * "
)

// DayTimeToCron converts a time.Time to a daily cron expression.
// DayTimeToCron 将 time.Time 转换为每日的 Cron 表达式。
func DayTimeToCron(t time.Time) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(DayTimeType, minute, hour)
}

// WeekTimeToCron converts a time.Time and weekday to a weekly cron expression.
// WeekTimeToCron 将 time.Time 和星期几转换为每周的 Cron 表达式。
func WeekTimeToCron(t time.Time, week time.Weekday) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(WeekTimeType, minute, hour, week)
}

// MonthTimeToCron converts a time.Time and month to a monthly cron expression.
// MonthTimeToCron 将 time.Time 和月份转换为每月的 Cron 表达式。
func MonthTimeToCron(t time.Time, month time.Month) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(MonthTimeType, minute, hour, month)
}

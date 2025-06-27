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

// DailyJob represents a job that runs on a daily schedule.
// DailyJob 表示一个按天调度运行的任务。
type DailyJob struct {
	// Unique identifier for the job
	// 任务的唯一标识符
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
	// Interval in days between job runs
	// 任务运行的天数间隔
	Interval uint
	// Specific times of day to run the job
	// 任务每天运行的具体时间点
	AtTimes gocron.AtTimes
	// Tags for the job
	// 任务的标签
	Tags []string
	// The function to execute as the job
	// 作为任务执行的函数
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
	// Error state for the job
	// 任务的错误状态
	err error
}

func NewDailyJob(interval uint, atTime gocron.AtTimes) *DailyJob {
	return &DailyJob{
		Interval: interval,
		AtTimes:  atTime,
		Type:     JobTypeDaily,
	}
}

// NewDailyJobAtTime creates a new DailyJob that runs at a specific time every day.
// NewDailyJobAtTime 创建一个每天在特定时间运行的 DailyJob。
func NewDailyJobAtTime(hour, minute, second uint) *DailyJob {
	return &DailyJob{
		Interval: 1,
		AtTimes:  gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
	}
}

// AtDayTime sets the specific time of day to run the job every day.
func (c *DailyJob) AtDayTime(hour, minute, second uint) *DailyJob {
	c.AtTimes = gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second))
	return c
}

// Alias sets the alias for the DailyJob.
// Alias 设置 DailyJob 的别名。
func (c *DailyJob) Alias(alias string) *DailyJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the DailyJob.
// JobID 设置 DailyJob 的唯一标识符。
func (c *DailyJob) JobID(id string) *DailyJob {
	c.ID = id
	return c
}

// Names sets the name for the DailyJob. If name is empty, a UUID is generated.
// Names 设置 DailyJob 的名称。如果名称为空，则生成一个 UUID。
func (c *DailyJob) Names(name string) *DailyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *DailyJob) Tag(tags ...string) *DailyJob {
	c.Tags = tags
	return c
}

// Task sets the task function and its parameters for the DailyJob.
// It wraps the task with error and timeout handling.
// Task 设置 DailyJob 的任务函数及其参数，并包装错误和超时处理。
func (c *DailyJob) Task(task any, parameters ...any) *DailyJob {
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

// Timeout sets the timeout duration for the DailyJob execution.
// Timeout 设置 DailyJob 执行的超时时间。
func (c *DailyJob) Timeout(timeout time.Duration) *DailyJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

// Watch sets a watcher function for job events.
// Watch 设置任务事件的监听函数。
func (c *DailyJob) Watch(watch func(event JobWatchInterface)) *DailyJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the DailyJob.
// addHooks 向 DailyJob 添加一个或多个事件监听器（钩子）。
func (c *DailyJob) addHooks(hook ...gocron.EventListener) *DailyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the DailyJob.
// DefaultHooks 向 DailyJob 添加一组默认事件监听器。
func (c *DailyJob) DefaultHooks() *DailyJob {
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
func (c *DailyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the job runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
func (c *DailyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *DailyJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the job runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
func (c *DailyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the job runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
func (c *DailyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the job panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
func (c *DailyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *DailyJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during job execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
func (c *DailyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

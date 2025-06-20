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

// IntervalJob represents a job that runs at a fixed interval.
// IntervalJob 表示一个按固定时间间隔运行的任务。
type IntervalJob struct {
	// 任务的唯一标识符
	ID string // Unique identifier for the job
	// 任务的别名
	Ali string // Alias for the job
	// 任务的名称
	Name string // Name of the job
	// 任务的时间间隔
	Interval time.Duration // Interval duration between job runs
	// 任务的执行函数
	TaskFunc any // The function to execute as the job
	// 任务的参数
	Parameters []any // Parameters to pass to the task function
	// 任务的钩子函数
	Hooks []gocron.EventListener // Event hooks for job lifecycle events
	//  监听任务事件的函数
	WatchFunc func(event JobWatchInterface) // Function to watch job events
	// 任务执行的超时时间
	timeout time.Duration // Timeout for the job execution
	// 任务的错误状态
	err error // Error state for the job
}

// NewIntervalJob creates a new IntervalJob with the specified interval.
// NewIntervalJob 创建一个具有指定时间间隔的 IntervalJob。
func NewIntervalJob(interval time.Duration) *IntervalJob {
	return &IntervalJob{
		Interval: interval,
	}
}

// Error returns the error message associated with the IntervalJob, if any.
// Error 返回与 IntervalJob 相关的错误信息（如果有）。
func (c *IntervalJob) Error() string {
	return c.err.Error()
}

func (c *IntervalJob) IntervalTime(interval time.Duration) *IntervalJob {
	c.Interval = interval
	return c
}

// Alias sets the alias for the IntervalJob.
// Alias 设置 IntervalJob 的别名。
func (c *IntervalJob) Alias(alias string) *IntervalJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the IntervalJob.
// JobID 设置 IntervalJob 的唯一标识符。
func (c *IntervalJob) JobID(id string) *IntervalJob {
	c.ID = id
	return c
}

// Names sets the name for the IntervalJob. If name is empty, a UUID is generated.
// Names 设置 IntervalJob 的名称。如果名称为空，则生成一个 UUID。
func (c *IntervalJob) Names(name string) *IntervalJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Timeout sets the timeout duration for the IntervalJob execution.
// Timeout 设置 IntervalJob 执行的超时时间。
func (c *IntervalJob) Timeout(timeout time.Duration) *IntervalJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

// Task sets the task function and its parameters for the IntervalJob.
// It wraps the task with error and timeout handling.
// Task 设置 IntervalJob 的任务函数及其参数，并包装错误和超时处理。
func (c *IntervalJob) Task(task any, parameters ...any) *IntervalJob {
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

// Watch sets a watcher function for job events.
// Watch 设置任务事件的监听函数。
func (c *IntervalJob) Watch(watch func(event JobWatchInterface)) *IntervalJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the IntervalJob.
// addHooks 向 IntervalJob 添加一个或多个事件监听器（钩子）。
func (c *IntervalJob) addHooks(hook ...gocron.EventListener) *IntervalJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the IntervalJob.
// DefaultHooks 向 IntervalJob 添加一组默认事件监听器。
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
func (c *IntervalJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the job runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
func (c *IntervalJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the job runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
func (c *IntervalJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the job runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
func (c *IntervalJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the job panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
func (c *IntervalJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during job execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
func (c *IntervalJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

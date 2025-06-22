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

// WeeklyJob represents a job that runs on a weekly schedule.
// WeeklyJob 表示一个按周调度运行的任务。
type WeeklyJob struct {
	// Unique identifier for the job
	// 任务的唯一标识符
	ID string
	// Alias for the job
	// 任务的别名
	Ali string
	// Name of the job
	// 任务名称
	Name string
	// Interval in weeks between job runs
	// 任务运行的周数间隔
	Interval uint
	// Days of the week to run the job
	// 任务每周运行的具体星期几
	DaysOfTheWeek gocron.Weekdays
	// Specific times of day to run the job
	// 任务每天运行的具体时间点
	WorkTimes gocron.AtTimes
	// Tag for the job
	// 任务的标签
	Tag []string
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

func NewWeeklyJob(interval uint, days gocron.Weekdays, atTime gocron.AtTimes) *WeeklyJob {
	return &WeeklyJob{
		Interval:      interval,
		DaysOfTheWeek: days,
		WorkTimes:     atTime,
	}
}

// NewWeeklyJobAtTime creates a new WeeklyJob that runs at specific days and times every week.
// NewWeeklyJobAtTime 创建一个每周在特定星期和时间运行的 WeeklyJob。
func NewWeeklyJobAtTime(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	if len(days) == 0 {
		return &WeeklyJob{
			err: ErrAtTimeDaysNil,
		}
	}
	return &WeeklyJob{
		Interval:      1,
		DaysOfTheWeek: gocron.NewWeekdays(days[0], days[1:]...),
		WorkTimes:     gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
	}
}

func (c *WeeklyJob) AtTimes(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	if len(days) == 0 {
		return &WeeklyJob{
			err: ErrAtTimeDaysNil,
		}
	}
	return &WeeklyJob{
		Interval:      1,
		DaysOfTheWeek: gocron.NewWeekdays(days[0], days[1:]...),
		WorkTimes:     gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
	}
}

// Alias sets the alias for the WeeklyJob.
// Alias 设置 WeeklyJob 的别名。
func (c *WeeklyJob) Alias(alias string) *WeeklyJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the WeeklyJob.
// JobID 设置 WeeklyJob 的唯一标识符。
func (c *WeeklyJob) JobID(id string) *WeeklyJob {
	c.ID = id
	return c
}

// Names sets the name for the WeeklyJob. If name is empty, a UUID is generated.
// Names 设置 WeeklyJob 的名称。如果名称为空，则生成一个 UUID。
func (c *WeeklyJob) Names(name string) *WeeklyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tags sets the tags for the WeeklyJob.
// Tags 设置 WeeklyJob 的标签
func (c *WeeklyJob) Tags(tags ...string) *WeeklyJob {
	c.Tag = tags
	return c
}

// Task sets the task function and its parameters for the WeeklyJob.
// It wraps the task with error and timeout handling.
// Task 设置 WeeklyJob 的任务函数及其参数，并包装错误和超时处理。
func (c *WeeklyJob) Task(task any, parameters ...any) *WeeklyJob {
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

// Timeout sets the timeout duration for the WeeklyJob execution.
// Timeout 设置 WeeklyJob 执行的超时时间。
func (c *WeeklyJob) Timeout(timeout time.Duration) *WeeklyJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

// Watch sets a watcher function for job events.
// Watch 设置任务事件的监听函数。
func (c *WeeklyJob) Watch(watch func(event JobWatchInterface)) *WeeklyJob {
	c.WatchFunc = watch
	return c
}

// addHooks adds one or more event listeners (hooks) to the WeeklyJob.
// addHooks 向 WeeklyJob 添加一个或多个事件监听器（钩子）。
func (c *WeeklyJob) addHooks(hook ...gocron.EventListener) *WeeklyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

// DefaultHooks adds a set of default event listeners to the WeeklyJob.
// DefaultHooks 向 WeeklyJob 添加一组默认事件监听器。
func (c *WeeklyJob) DefaultHooks() *WeeklyJob {
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
func (c *WeeklyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors adds a hook to be called before the job runs, skipping if the hook returns an error.
// BeforeJobRunsSkipIfBeforeFuncErrors 添加一个在任务运行前调用的钩子，如果钩子返回错误则跳过。
func (c *WeeklyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns adds a hook to be called after the job runs.
// AfterJobRuns 添加一个在任务运行后调用的钩子。
func (c *WeeklyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError adds a hook to be called after the job runs with an error.
// AfterJobRunsWithError 添加一个在任务运行出错后调用的钩子。
func (c *WeeklyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic adds a hook to be called after the job panics.
// AfterJobRunsWithPanic 添加一个在任务发生 panic 后调用的钩子。
func (c *WeeklyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError adds a hook to be called when a lock error occurs during job execution.
// AfterLockError 添加一个在任务加锁出错时调用的钩子。
func (c *WeeklyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

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

// MonthJob represents a job that runs on a monthly schedule.
// MonthJob 表示一个按月调度运行的任务。
type MonthJob struct {
	// Unique identifier for the job
	// 任务的唯一标识符
	ID string
	// Alias for the job
	// 任务的别名
	Ali string
	// Name of the job
	// 任务名称
	Name string
	// Interval in months between job runs
	// 任务运行的月数间隔
	Interval uint
	// Days of the month to run the job
	// 任务每月运行的具体日期
	DaysOfTheMonth gocron.DaysOfTheMonth
	// Specific times of day to run the job
	// 任务每天运行的具体时间点
	AtTimes gocron.AtTimes
	// Tags for the job
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

func NewMonthJob(interval uint, days gocron.DaysOfTheMonth, atTime gocron.AtTimes) *MonthJob {
	return &MonthJob{
		Interval:       interval,
		DaysOfTheMonth: days,
		AtTimes:        atTime,
	}
}

// NewMonthJobAtTime creates a new MonthJob that runs at specific days and times every month.
// NewMonthJobAtTime 创建一个每月在特定日期和时间运行的 MonthJob。
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
func (c *MonthJob) Alias(alias string) *MonthJob {
	c.Ali = alias
	return c
}

// JobID sets the unique identifier for the MonthJob.
// JobID 设置 MonthJob 的唯一标识符。
func (c *MonthJob) JobID(id string) *MonthJob {
	c.ID = id
	return c
}

// Names sets the name for the MonthJob. If name is empty, a UUID is generated.
// Names 设置 MonthJob 的名称。如果名称为空，则生成一个 UUID。
func (c *MonthJob) Names(name string) *MonthJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

// Tags sets the tags for the MonthJob.
// Tags 设置 MonthJob 的标签
func (c *MonthJob) Tags(tags ...string) *MonthJob {
	c.Tag = tags
	return c
}

// Task sets the task function and its parameters for the MonthJob.
// It wraps the task with error and timeout handling.
// Task 设置 MonthJob 的任务函数及其参数，并包装错误和超时处理。
func (c *MonthJob) Task(task any, parameters ...any) *MonthJob {
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

// Timeout sets the timeout duration for the MonthJob execution.
// Timeout 设置 MonthJob 执行的超时时间。
func (c *MonthJob) Timeout(timeout time.Duration) *MonthJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

func (c *MonthJob) Watch(watch func(event JobWatchInterface)) *MonthJob {
	c.WatchFunc = watch
	return c
}

func (c *MonthJob) addHooks(hook ...gocron.EventListener) *MonthJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *MonthJob) DefaultHooks() *MonthJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *MonthJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *MonthJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *MonthJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *MonthJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *MonthJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *MonthJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *MonthJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *MonthJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

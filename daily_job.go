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

type DailyJob struct {
	ID         string
	Ali        string
	Name       string
	Interval   uint
	AtTimes    gocron.AtTimes
	TaskFunc   any
	Parameters []any
	Hooks      []gocron.EventListener
	WatchFunc  func(event JobWatchInterface)
	timeout    time.Duration
	err        error
}

func NewDailyJob(interval uint, atTime gocron.AtTimes) *DailyJob {
	return &DailyJob{
		Interval: interval,
		AtTimes:  atTime,
	}
}

func NewDailyJobAtTime(hour, minute, second uint) *DailyJob {
	return &DailyJob{
		Interval: 1,
		AtTimes:  gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
	}
}

func (c *DailyJob) Error() string {
	return c.err.Error()
}

func (c *DailyJob) Alias(alias string) *DailyJob {
	c.Ali = alias
	return c
}

func (c *DailyJob) JobID(id string) *DailyJob {
	c.ID = id
	return c
}

func (c *DailyJob) Names(name string) *DailyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

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

func (c *DailyJob) Timeout(timeout time.Duration) *DailyJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

func (c *DailyJob) Watch(watch func(event JobWatchInterface)) *DailyJob {
	c.WatchFunc = watch
	return c
}

func (c *DailyJob) addHooks(hook ...gocron.EventListener) *DailyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *DailyJob) DefaultHooks() *DailyJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *DailyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *DailyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *DailyJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *DailyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *DailyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *DailyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *DailyJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *DailyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

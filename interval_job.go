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

type IntervalJob struct {
	ID         string
	Ali        string
	Name       string
	Interval   time.Duration
	TaskFunc   any
	Parameters []any
	Hooks      []gocron.EventListener
	WatchFunc  func(event JobWatchInterface)
	timeout    time.Duration
	err        error
}

func NewIntervalJob(interval time.Duration) *IntervalJob {
	return &IntervalJob{
		Interval: interval,
	}
}

func (c *IntervalJob) Error() string {
	return c.err.Error()
}

func (c *IntervalJob) Alias(alias string) *IntervalJob {
	c.Ali = alias
	return c
}

func (c *IntervalJob) JobID(id string) *IntervalJob {
	c.ID = id
	return c
}

func (c *IntervalJob) Names(name string) *IntervalJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *IntervalJob) Timeout(timeout time.Duration) *IntervalJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

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

func (c *IntervalJob) Watch(watch func(event JobWatchInterface)) *IntervalJob {
	c.WatchFunc = watch
	return c
}

func (c *IntervalJob) addHooks(hook ...gocron.EventListener) *IntervalJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *IntervalJob) DefaultHooks() *IntervalJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *IntervalJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *IntervalJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *IntervalJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *IntervalJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *IntervalJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *IntervalJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

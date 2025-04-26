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

type MonthJob struct {
	ID             string
	Ali            string
	Name           string
	Interval       uint
	DaysOfTheMonth gocron.DaysOfTheMonth
	AtTimes        gocron.AtTimes
	TaskFunc       any
	Parameters     []any
	Hooks          []gocron.EventListener
	WatchFunc      func(event JobWatchInterface)
	timeout        time.Duration
	err            error
}

func NewMonthJob(interval uint, days gocron.DaysOfTheMonth, atTime gocron.AtTimes) *MonthJob {
	return &MonthJob{
		Interval:       interval,
		DaysOfTheMonth: days,
		AtTimes:        atTime,
	}
}

func NewMonthJobAtTime(days []int, hour, minute, second int) *MonthJob {
	if len(days) == 0 {
		err := errors.New("chrono:at time must have at least one day")
		return &MonthJob{
			err: err,
		}
	}
	return &MonthJob{
		Interval:       1,
		DaysOfTheMonth: gocron.NewDaysOfTheMonth(days[0], days[1:]...),
		AtTimes:        gocron.NewAtTimes(gocron.NewAtTime(uint(hour), uint(minute), uint(second))),
	}
}

func (c *MonthJob) Error() string {
	return c.err.Error()
}

func (c *MonthJob) Alias(alias string) *MonthJob {
	c.Ali = alias
	return c
}

func (c *MonthJob) JobID(id string) *MonthJob {
	c.ID = id
	return c
}

func (c *MonthJob) Names(name string) *MonthJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

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

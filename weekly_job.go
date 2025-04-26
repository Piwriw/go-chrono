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

type WeeklyJob struct {
	ID            string
	Ali           string
	Name          string
	Interval      uint
	DaysOfTheWeek gocron.Weekdays
	AtTimes       gocron.AtTimes
	TaskFunc      any
	Parameters    []any
	Hooks         []gocron.EventListener
	WatchFunc     func(event JobWatchInterface)
	timeout       time.Duration
	err           error
}

func NewWeeklyJob(interval uint, days gocron.Weekdays, atTime gocron.AtTimes) *WeeklyJob {
	return &WeeklyJob{
		Interval:      interval,
		DaysOfTheWeek: days,
		AtTimes:       atTime,
	}
}

func NewWeeklyJobAtTime(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	if len(days) == 0 {
		err := errors.New("chrono:at time must have at least one day")
		return &WeeklyJob{
			err: err,
		}
	}
	return &WeeklyJob{
		Interval:      1,
		DaysOfTheWeek: gocron.NewWeekdays(days[0], days[1:]...),
		AtTimes:       gocron.NewAtTimes(gocron.NewAtTime(hour, minute, second)),
	}
}

func (c *WeeklyJob) Error() string {
	return c.err.Error()
}

func (c *WeeklyJob) Alias(alias string) *WeeklyJob {
	c.Ali = alias
	return c
}

func (c *WeeklyJob) JobID(id string) *WeeklyJob {
	c.ID = id
	return c
}

func (c *WeeklyJob) Names(name string) *WeeklyJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

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

func (c *WeeklyJob) Timeout(timeout time.Duration) *WeeklyJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

func (c *WeeklyJob) Watch(watch func(event JobWatchInterface)) *WeeklyJob {
	c.WatchFunc = watch
	return c
}

func (c *WeeklyJob) addHooks(hook ...gocron.EventListener) *WeeklyJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *WeeklyJob) DefaultHooks() *WeeklyJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *WeeklyJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *WeeklyJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *WeeklyJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *WeeklyJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *WeeklyJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *WeeklyJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

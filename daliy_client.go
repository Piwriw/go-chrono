package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type DailyJobClientInterface interface {
	// AtDayTime 设置任务在每天的指定时间运行
	AtDayTime(hour, minute, second uint) DailyJobClientInterface
	Alias(alias string) DailyJobClientInterface
	JobID(id string) DailyJobClientInterface
	Name(name string) DailyJobClientInterface
	Tags(tags ...string) DailyJobClientInterface
	Task(task any, parameters ...any) DailyJobClientInterface
	Timeout(timeout time.Duration) DailyJobClientInterface
	Watch(watch func(event JobWatchInterface)) DailyJobClientInterface
	DefaultHooks() DailyJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) DailyJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) DailyJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface
	Add() (gocron.Job, error)
	BatchAdd(dailyJobs ...*DailyJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type DailyJobClient struct {
	scheduler *Scheduler
	job       *DailyJob
}

var _ DailyJobClientInterface = (*DailyJobClient)(nil)

func (c *DailyJobClient) AtDayTime(hour, minute, second uint) DailyJobClientInterface {
	c.job.AtDayTime(hour, minute, second)
	return c
}

func (c *DailyJobClient) Alias(alias string) DailyJobClientInterface {
	c.job.Alias(alias)
	return c
}

func (c *DailyJobClient) JobID(id string) DailyJobClientInterface {
	c.job.JobID(id)
	return c
}

func (c *DailyJobClient) Name(name string) DailyJobClientInterface {
	c.job.Names(name)
	return c
}

func (c *DailyJobClient) Tags(tags ...string) DailyJobClientInterface {
	c.job.Tag(tags...)
	return c
}

func (c *DailyJobClient) Task(task any, parameters ...any) DailyJobClientInterface {
	c.job.Task(task, parameters...)
	return c
}

func (c *DailyJobClient) Timeout(timeout time.Duration) DailyJobClientInterface {
	c.job.Timeout(timeout)
	return c
}

func (c *DailyJobClient) Watch(watch func(event JobWatchInterface)) DailyJobClientInterface {
	c.job.Watch(watch)
	return c
}

func (c *DailyJobClient) DefaultHooks() DailyJobClientInterface {
	c.job.DefaultHooks()
	return c
}

func (c *DailyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface {
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
}

func (c *DailyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) DailyJobClientInterface {
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

func (c *DailyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) DailyJobClientInterface {
	c.job.AfterJobRuns(eventListenerFunc)
	return c
}

func (c *DailyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface {
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
}

func (c *DailyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) DailyJobClientInterface {
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

func (c *DailyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) DailyJobClientInterface {
	c.job.AfterLockError(eventListenerFunc)
	return c
}

func (c *DailyJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrDailyJobNil
	}
	return c.scheduler.AddDailyJob(c.job)
}

func (c *DailyJobClient) BatchAdd(dailyJobs ...*DailyJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	return c.scheduler.AddDailyJobs(dailyJobs...)
}

func (c *DailyJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job == nil {
		return ErrCronJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.RemoveJobByAlias(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.RemoveJobByName(c.job.Name)
	}
	return ErrJobNotFound
}

// Get 获取任务
func (c *DailyJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.GetJobByID(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.GetJobByAlias(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.GetJobByName(c.job.Name)
	}
	return nil, ErrJobNotFound
}

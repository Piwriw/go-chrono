package chrono

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
)

type DailyJobClientInterface interface {
	// AtDayTime 设置任务在每天的指定时间运行
	AtDayTime(hour, minute, second uint) *DailyJob
	Alias(alias string) *DailyJob
	JobID(id string) *DailyJob
	Name(name string) *DailyJob
	Task(task any, parameters ...any) *DailyJob
	Timeout(timeout time.Duration) *DailyJob
	Watch(watch func(event JobWatchInterface)) *DailyJob
	DefaultHooks() *DailyJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *DailyJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *DailyJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob
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

func (c *DailyJobClient) AtDayTime(hour, minute, second uint) *DailyJob {
	return c.job.AtDayTime(hour, minute, second)
}

func (c *DailyJobClient) Alias(alias string) *DailyJob {
	return c.job.Alias(alias)
}

func (c *DailyJobClient) JobID(id string) *DailyJob {
	return c.job.JobID(id)
}

func (c *DailyJobClient) Name(name string) *DailyJob {
	return c.job.Names(name)
}

func (c *DailyJobClient) Task(task any, parameters ...any) *DailyJob {
	return c.job.Task(task, parameters...)
}

func (c *DailyJobClient) Timeout(timeout time.Duration) *DailyJob {
	return c.job.Timeout(timeout)
}

func (c *DailyJobClient) Watch(watch func(event JobWatchInterface)) *DailyJob {
	return c.job.Watch(watch)
}

func (c *DailyJobClient) DefaultHooks() *DailyJob {
	return c.job.DefaultHooks()
}

func (c *DailyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.job.BeforeJobRuns(eventListenerFunc)
}

func (c *DailyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *DailyJob {
	return c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (c *DailyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *DailyJob {
	return c.job.AfterJobRuns(eventListenerFunc)
}

func (c *DailyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.job.AfterJobRunsWithError(eventListenerFunc)
}

func (c *DailyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *DailyJob {
	return c.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (c *DailyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *DailyJob {
	return c.job.AfterLockError(eventListenerFunc)
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

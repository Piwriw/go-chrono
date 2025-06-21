package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type IntervalJobClientInterface interface {
	// Interval 设置间隔时间
	Interval(interval time.Duration) *IntervalJob
	Alias(alias string) *IntervalJob
	JobID(id string) *IntervalJob
	Name(name string) *IntervalJob
	Task(task any, parameters ...any) *IntervalJob
	Timeout(timeout time.Duration) *IntervalJob
	Watch(watch func(event JobWatchInterface)) *IntervalJob
	DefaultHooks() *IntervalJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob
	Add() (gocron.Job, error)
	BatchAdd(intervalJobs ...*IntervalJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type IntervalJobClient struct {
	scheduler *Scheduler
	job       *IntervalJob
}

var _ IntervalJobClientInterface = (*IntervalJobClient)(nil)

func (c *IntervalJobClient) Interval(interval time.Duration) *IntervalJob {
	return c.job.IntervalTime(interval)
}

func (c *IntervalJobClient) Alias(alias string) *IntervalJob {
	return c.job.Alias(alias)
}

func (c *IntervalJobClient) JobID(id string) *IntervalJob {
	return c.job.JobID(id)
}

func (c *IntervalJobClient) Name(name string) *IntervalJob {
	return c.job.Names(name)
}

func (c *IntervalJobClient) Task(task any, parameters ...any) *IntervalJob {
	return c.job.Task(task, parameters...)
}

func (c *IntervalJobClient) Timeout(timeout time.Duration) *IntervalJob {
	return c.job.Timeout(timeout)
}

func (c *IntervalJobClient) Watch(watch func(event JobWatchInterface)) *IntervalJob {
	return c.job.Watch(watch)
}

func (c *IntervalJobClient) DefaultHooks() *IntervalJob {
	return c.job.DefaultHooks()
}

func (c *IntervalJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.job.BeforeJobRuns(eventListenerFunc)
}

func (c *IntervalJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *IntervalJob {
	return c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (c *IntervalJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *IntervalJob {
	return c.job.AfterJobRuns(eventListenerFunc)
}

func (c *IntervalJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.job.AfterJobRunsWithError(eventListenerFunc)
}

func (c *IntervalJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *IntervalJob {
	return c.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (c *IntervalJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *IntervalJob {
	return c.job.AfterLockError(eventListenerFunc)
}

func (c *IntervalJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job != nil {
		return nil, ErrIntervalJobNil
	}
	return c.scheduler.AddIntervalJob(c.job)
}

func (c *IntervalJobClient) BatchAdd(intervalJobs ...*IntervalJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job != nil {
		return nil, ErrIntervalJobNil
	}
	return c.scheduler.AddIntervalJobs(intervalJobs...)
}

func (c *IntervalJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job != nil {
		return ErrIntervalJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.RemoveJob(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.RemoveJob(c.job.Name)
	}
	return ErrJobNotFound
}

func (c *IntervalJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job != nil {
		return nil, ErrIntervalJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.GetJobByID(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.GetJobByID(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.GetJobByName(c.job.Name)
	}
	return nil, ErrJobNotFound
}

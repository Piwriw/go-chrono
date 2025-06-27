package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type IntervalJobClientInterface interface {
	// Interval 设置间隔时间
	Interval(interval time.Duration) IntervalJobClientInterface
	Alias(alias string) IntervalJobClientInterface
	JobID(id string) IntervalJobClientInterface
	Name(name string) IntervalJobClientInterface
	Tag(tags ...string) IntervalJobClientInterface
	Task(task any, parameters ...any) IntervalJobClientInterface
	Timeout(timeout time.Duration) IntervalJobClientInterface
	Watch(watch func(event JobWatchInterface)) IntervalJobClientInterface
	DefaultHooks() IntervalJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) IntervalJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) IntervalJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface
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

func (c *IntervalJobClient) Interval(interval time.Duration) IntervalJobClientInterface {
	c.job.IntervalTime(interval)
	return c
}

func (c *IntervalJobClient) Alias(alias string) IntervalJobClientInterface {
	c.job.Alias(alias)
	return c
}

func (c *IntervalJobClient) JobID(id string) IntervalJobClientInterface {
	c.job.JobID(id)
	return c
}

func (c *IntervalJobClient) Name(name string) IntervalJobClientInterface {
	c.job.Names(name)
	return c
}

func (c *IntervalJobClient) Tag(tags ...string) IntervalJobClientInterface {
	c.job.Tag(tags...)
	return c
}

func (c *IntervalJobClient) Task(task any, parameters ...any) IntervalJobClientInterface {
	c.job.Task(task, parameters...)
	return c
}

func (c *IntervalJobClient) Timeout(timeout time.Duration) IntervalJobClientInterface {
	c.job.Timeout(timeout)
	return c
}

func (c *IntervalJobClient) Watch(watch func(event JobWatchInterface)) IntervalJobClientInterface {
	c.job.Watch(watch)
	return c
}

func (c *IntervalJobClient) DefaultHooks() IntervalJobClientInterface {
	c.job.DefaultHooks()
	return c
}

func (c *IntervalJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface {
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) IntervalJobClientInterface {
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) IntervalJobClientInterface {
	c.job.AfterJobRuns(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface {
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) IntervalJobClientInterface {
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) IntervalJobClientInterface {
	c.job.AfterLockError(eventListenerFunc)
	return c
}

func (c *IntervalJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrIntervalJobNil
	}
	return c.scheduler.AddIntervalJob(c.job)
}

func (c *IntervalJobClient) BatchAdd(intervalJobs ...*IntervalJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrIntervalJobNil
	}
	return c.scheduler.AddIntervalJobs(intervalJobs...)
}

func (c *IntervalJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job == nil {
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
	if c.job == nil {
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

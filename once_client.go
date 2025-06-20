package chrono

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
)

type OnceJobClientInterface interface {
	// AtTimes 设置运行时间
	AtTimes(workTimes ...time.Time) *OnceJob
	Alias(alias string) *OnceJob
	JobID(id string) *OnceJob
	Name(name string) *OnceJob
	Task(task any, parameters ...any) *OnceJob
	Timeout(timeout time.Duration) *OnceJob
	Watch(watch func(event JobWatchInterface)) *OnceJob
	DefaultHooks() *OnceJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *OnceJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *OnceJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob
	Add() (gocron.Job, error)
	BatchAdd(cronJobs ...*OnceJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type OnceJobClient struct {
	scheduler *Scheduler
	job       *OnceJob
}

var _ OnceJobClientInterface = (*OnceJobClient)(nil)

func (o *OnceJobClient) AtTimes(workTimes ...time.Time) *OnceJob {
	return o.job.AtTimes(workTimes...)
}

func (o *OnceJobClient) Alias(alias string) *OnceJob {
	return o.job.Alias(alias)
}

func (o *OnceJobClient) JobID(id string) *OnceJob {
	return o.job.JobID(id)
}

func (o *OnceJobClient) Name(name string) *OnceJob {
	return o.job.Names(name)
}

func (o *OnceJobClient) Task(task any, parameters ...any) *OnceJob {
	return o.job.Task(task, parameters...)
}

func (o *OnceJobClient) Timeout(timeout time.Duration) *OnceJob {
	return o.job.Timeout(timeout)
}

func (o *OnceJobClient) Watch(watch func(event JobWatchInterface)) *OnceJob {
	return o.job.Watch(watch)
}

func (o *OnceJobClient) DefaultHooks() *OnceJob {
	return o.job.DefaultHooks()
}

func (o *OnceJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return o.job.BeforeJobRuns(eventListenerFunc)
}

func (o *OnceJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *OnceJob {
	return o.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (o *OnceJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *OnceJob {
	return o.job.AfterJobRuns(eventListenerFunc)
}

func (o *OnceJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return o.job.AfterJobRunsWithError(eventListenerFunc)
}

func (o *OnceJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *OnceJob {
	return o.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (o *OnceJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *OnceJob {
	return o.job.AfterLockError(eventListenerFunc)
}

func (o *OnceJobClient) Add() (gocron.Job, error) {
	if o.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if o.job == nil {
		return nil, ErrOnceJobNil
	}
	return o.scheduler.AddOnceJob(o.job)
}

func (o *OnceJobClient) BatchAdd(onceJobs ...*OnceJob) ([]gocron.Job, error) {
	if o.scheduler == nil {
		return nil, ErrScheduleNil
	}
	return o.scheduler.AddOnceJobs(onceJobs...)
}

func (o *OnceJobClient) Remove() error {
	if o.scheduler == nil {
		return ErrScheduleNil
	}
	if o.job == nil {
		return ErrOnceJobNil
	}
	if o.job.ID != "" {
		return o.scheduler.RemoveJob(o.job.ID)
	}
	if o.job.Ali != "" {
		return o.scheduler.RemoveJobByAlias(o.job.Ali)
	}
	if o.job.Name != "" {
		return o.scheduler.RemoveJobByName(o.job.Name)
	}
	return ErrJobNotFound
}

func (o *OnceJobClient) Get() (gocron.Job, error) {
	if o.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if o.job == nil {
		return nil, ErrOnceJobNil
	}
	if o.job.ID != "" {
		return o.scheduler.GetJobByID(o.job.ID)
	}
	if o.job.Ali != "" {
		return o.scheduler.GetJobByAlias(o.job.Ali)
	}
	if o.job.Name != "" {
		return o.scheduler.GetJobByName(o.job.Name)
	}
	return nil, ErrJobNotFound
}

package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type OnceJobClientInterface interface {
	// AtTimes 设置运行时间
	AtTimes(workTimes ...time.Time) OnceJobClientInterface
	Alias(alias string) OnceJobClientInterface
	JobID(id string) OnceJobClientInterface
	Name(name string) OnceJobClientInterface
	Tag(tags ...string) OnceJobClientInterface
	Task(task any, parameters ...any) OnceJobClientInterface
	Timeout(timeout time.Duration) OnceJobClientInterface
	Watch(watch func(event JobWatchInterface)) OnceJobClientInterface
	DefaultHooks() OnceJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) OnceJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) OnceJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface
	Add() (gocron.Job, error)
	BatchAdd(onceJobs ...*OnceJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type OnceJobClient struct {
	scheduler *Scheduler
	job       *OnceJob
}

var _ OnceJobClientInterface = (*OnceJobClient)(nil)

func (o *OnceJobClient) AtTimes(workTimes ...time.Time) OnceJobClientInterface {
	o.job.AtTimes(workTimes...)
	return o
}

func (o *OnceJobClient) Alias(alias string) OnceJobClientInterface {
	o.job.Alias(alias)
	return o
}

func (o *OnceJobClient) JobID(id string) OnceJobClientInterface {
	o.job.JobID(id)
	return o
}

func (o *OnceJobClient) Name(name string) OnceJobClientInterface {
	o.job.Names(name)
	return o
}

func (o *OnceJobClient) Tag(tags ...string) OnceJobClientInterface {
	o.job.Tags(tags...)
	return o
}

func (o *OnceJobClient) Task(task any, parameters ...any) OnceJobClientInterface {
	o.job.Task(task, parameters...)
	return o
}

func (o *OnceJobClient) Timeout(timeout time.Duration) OnceJobClientInterface {
	o.job.Timeout(timeout)
	return o
}

func (o *OnceJobClient) Watch(watch func(event JobWatchInterface)) OnceJobClientInterface {
	o.job.Watch(watch)
	return o
}

func (o *OnceJobClient) DefaultHooks() OnceJobClientInterface {
	o.job.DefaultHooks()
	return o
}

func (o *OnceJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface {
	o.job.BeforeJobRuns(eventListenerFunc)
	return o
}

func (o *OnceJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) OnceJobClientInterface {
	o.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return o
}

func (o *OnceJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) OnceJobClientInterface {
	o.job.AfterJobRuns(eventListenerFunc)
	return o
}

func (o *OnceJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface {
	o.job.AfterJobRunsWithError(eventListenerFunc)
	return o
}

func (o *OnceJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) OnceJobClientInterface {
	o.job.AfterJobRunsWithPanic(eventListenerFunc)
	return o
}

func (o *OnceJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) OnceJobClientInterface {
	o.job.AfterLockError(eventListenerFunc)
	return o
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

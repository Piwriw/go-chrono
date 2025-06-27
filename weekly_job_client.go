package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type WeeklyJobClientInterface interface {
	// AtTimes 设置运行时间
	AtTimes(days []time.Weekday, hour, minute, second uint) WeeklyJobClientInterface
	Alias(alias string) WeeklyJobClientInterface
	JobID(id string) WeeklyJobClientInterface
	Name(name string) WeeklyJobClientInterface
	Tags(tags ...string) WeeklyJobClientInterface
	Task(task any, parameters ...any) WeeklyJobClientInterface
	Timeout(timeout time.Duration) WeeklyJobClientInterface
	Watch(watch func(event JobWatchInterface)) WeeklyJobClientInterface
	DefaultHooks() WeeklyJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) WeeklyJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) WeeklyJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface
	Add() (gocron.Job, error)
	BatchAdd(weeklyJobs ...*WeeklyJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type WeeklyJobClient struct {
	scheduler *Scheduler
	job       *WeeklyJob
}

var _ WeeklyJobClientInterface = (*WeeklyJobClient)(nil)

func (w *WeeklyJobClient) AtTimes(days []time.Weekday, hour, minute, second uint) WeeklyJobClientInterface {
	w.job.AtTimes(days, hour, minute, second)
	return w
}

func (w *WeeklyJobClient) Alias(alias string) WeeklyJobClientInterface {
	w.job.Alias(alias)
	return w
}

func (w *WeeklyJobClient) JobID(id string) WeeklyJobClientInterface {
	w.job.JobID(id)
	return w
}

func (w *WeeklyJobClient) Name(name string) WeeklyJobClientInterface {
	w.job.Names(name)
	return w
}

func (w *WeeklyJobClient) Tags(tags ...string) WeeklyJobClientInterface {
	w.job.Tags(tags...)
	return w
}

func (w *WeeklyJobClient) Task(task any, parameters ...any) WeeklyJobClientInterface {
	w.job.Task(task, parameters...)
	return w
}

func (w *WeeklyJobClient) Timeout(timeout time.Duration) WeeklyJobClientInterface {
	w.job.Timeout(timeout)
	return w
}

func (w *WeeklyJobClient) Watch(watch func(event JobWatchInterface)) WeeklyJobClientInterface {
	w.job.Watch(watch)
	return w
}

func (w *WeeklyJobClient) DefaultHooks() WeeklyJobClientInterface {
	w.job.DefaultHooks()
	return w
}

func (w *WeeklyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface {
	w.job.BeforeJobRuns(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) WeeklyJobClientInterface {
	w.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) WeeklyJobClientInterface {
	w.job.AfterJobRuns(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface {
	w.job.AfterJobRunsWithError(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) WeeklyJobClientInterface {
	w.job.AfterJobRunsWithPanic(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) WeeklyJobClientInterface {
	w.job.AfterLockError(eventListenerFunc)
	return w
}

func (w *WeeklyJobClient) Add() (gocron.Job, error) {
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	return w.scheduler.AddWeeklyJob(w.job)
}

func (w *WeeklyJobClient) BatchAdd(weeklyJobs ...*WeeklyJob) ([]gocron.Job, error) {
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	return w.scheduler.AddWeeklyJobs(weeklyJobs...)
}

func (w *WeeklyJobClient) Remove() error {
	if w.scheduler == nil {
		return ErrScheduleNil
	}
	if w.job == nil {
		return ErrWeeklyJobNil
	}
	if w.job.ID == "" {
		return w.scheduler.RemoveJob(w.job.Name)
	}
	if w.job.Ali != "" {
		return w.scheduler.RemoveJobByAlias(w.job.Ali)
	}
	if w.job.Name != "" {
		return w.scheduler.RemoveJobByName(w.job.ID)
	}
	return ErrJobNotFound
}

func (w *WeeklyJobClient) Get() (gocron.Job, error) {
	if w.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if w.job == nil {
		return nil, ErrWeeklyJobNil
	}
	if w.job.ID != "" {
		return w.scheduler.GetJobByID(w.job.ID)
	}
	if w.job.Ali != "" {
		return w.scheduler.GetJobByAlias(w.job.Ali)
	}
	if w.job.Name != "" {
		return w.scheduler.GetJobByName(w.job.Name)
	}
	return nil, ErrJobNotFound
}

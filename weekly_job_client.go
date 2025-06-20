package chrono

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
)

type WeeklyJobClientInterface interface {
	// AtTimes 设置运行时间
	AtTimes(days []time.Weekday, hour, minute, second uint) *WeeklyJob
	Alias(alias string) *WeeklyJob
	JobID(id string) *WeeklyJob
	Name(name string) *WeeklyJob
	Task(task any, parameters ...any) *WeeklyJob
	Timeout(timeout time.Duration) *WeeklyJob
	Watch(watch func(event JobWatchInterface)) *WeeklyJob
	DefaultHooks() *WeeklyJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob
	Add() (gocron.Job, error)
	BatchAdd(cronJobs ...*WeeklyJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type WeeklyJobClient struct {
	scheduler *Scheduler
	job       *WeeklyJob
}

var _ WeeklyJobClientInterface = (*WeeklyJobClient)(nil)

func (w *WeeklyJobClient) AtTimes(days []time.Weekday, hour, minute, second uint) *WeeklyJob {
	return w.job.AtTimes(days, hour, minute, second)
}

func (w *WeeklyJobClient) Alias(alias string) *WeeklyJob {
	return w.job.Alias(alias)
}

func (w *WeeklyJobClient) JobID(id string) *WeeklyJob {
	return w.job.JobID(id)
}

func (w *WeeklyJobClient) Name(name string) *WeeklyJob {
	return w.job.Names(name)
}

func (w *WeeklyJobClient) Task(task any, parameters ...any) *WeeklyJob {
	return w.job.Task(task, parameters...)
}

func (w *WeeklyJobClient) Timeout(timeout time.Duration) *WeeklyJob {
	return w.job.Timeout(timeout)
}

func (w *WeeklyJobClient) Watch(watch func(event JobWatchInterface)) *WeeklyJob {
	return w.job.Watch(watch)
}

func (w *WeeklyJobClient) DefaultHooks() *WeeklyJob {
	return w.job.DefaultHooks()
}

func (w *WeeklyJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return w.job.BeforeJobRuns(eventListenerFunc)
}

func (w *WeeklyJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *WeeklyJob {
	return w.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (w *WeeklyJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *WeeklyJob {
	return w.job.AfterJobRuns(eventListenerFunc)
}

func (w *WeeklyJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return w.job.AfterJobRunsWithError(eventListenerFunc)
}

func (w *WeeklyJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *WeeklyJob {
	return w.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (w *WeeklyJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *WeeklyJob {
	return w.job.AfterLockError(eventListenerFunc)
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

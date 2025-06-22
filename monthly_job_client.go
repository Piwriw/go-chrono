package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type MonthJobClientInterface interface {
	// AtTime 设置运行时间
	AtTime(days []int, hour, minute, second int) *MonthJob
	Alias(alias string) *MonthJob
	JobID(id string) *MonthJob
	Name(name string) *MonthJob
	Tags(tags ...string) *MonthJob
	Task(task any, parameters ...any) *MonthJob
	Timeout(timeout time.Duration) *MonthJob
	Watch(watch func(event JobWatchInterface)) *MonthJob
	DefaultHooks() *MonthJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *MonthJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *MonthJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob
	Add() (gocron.Job, error)
	BatchAdd(monthJobs ...*MonthJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type MonthJobClient struct {
	scheduler *Scheduler
	job       *MonthJob
}

var _ MonthJobClientInterface = (*MonthJobClient)(nil)

func (m *MonthJobClient) AtTime(days []int, hour, minute, second int) *MonthJob {
	return m.job.AtTime(days, hour, minute, second)
}

func (m *MonthJobClient) Alias(alias string) *MonthJob { // 设置别名
	return m.job.Alias(alias)
}

func (m *MonthJobClient) JobID(id string) *MonthJob {
	return m.job.JobID(id)
}

func (m *MonthJobClient) Name(name string) *MonthJob {
	return m.job.Names(name)
}

func (m *MonthJobClient) Tags(tags ...string) *MonthJob {
	return m.job.Tags(tags...)
}

func (m *MonthJobClient) Task(task any, parameters ...any) *MonthJob {
	return m.job.Task(task, parameters...)
}

func (m *MonthJobClient) Timeout(timeout time.Duration) *MonthJob {
	return m.job.Timeout(timeout)
}

func (m *MonthJobClient) Watch(watch func(event JobWatchInterface)) *MonthJob {
	return m.job.Watch(watch)
}

func (m *MonthJobClient) DefaultHooks() *MonthJob {
	return m.job.DefaultHooks()
}

func (m *MonthJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return m.job.BeforeJobRuns(eventListenerFunc)
}

func (m *MonthJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *MonthJob {
	return m.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (m *MonthJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *MonthJob {
	return m.job.AfterJobRuns(eventListenerFunc)
}

func (m *MonthJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return m.job.AfterJobRunsWithError(eventListenerFunc)
}

func (m *MonthJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *MonthJob {
	return m.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (m *MonthJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *MonthJob {
	return m.job.AfterLockError(eventListenerFunc)
}

func (m *MonthJobClient) Add() (gocron.Job, error) {
	if m.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if m.job == nil {
		return nil, ErrMonthJobNil
	}
	return m.scheduler.AddMonthlyJob(m.job)
}

func (m *MonthJobClient) BatchAdd(monthlyJobs ...*MonthJob) ([]gocron.Job, error) {
	if m.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if m.job == nil {
		return nil, ErrMonthJobNil
	}
	return m.scheduler.AddMonthlyJobs(monthlyJobs...)
}

func (m *MonthJobClient) Remove() error {
	if m.scheduler == nil {
		return ErrScheduleNil
	}
	if m.job == nil {
		return ErrMonthJobNil
	}
	if m.job.ID != "" {
		return m.scheduler.RemoveJob(m.job.ID)
	}
	if m.job.Ali != "" {
		return m.scheduler.RemoveJobByAlias(m.job.Ali)
	}
	if m.job.Name != "" {
		return m.scheduler.RemoveJobByName(m.job.Name)
	}
	return ErrJobNotFound
}

func (m *MonthJobClient) Get() (gocron.Job, error) {
	if m.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if m.job == nil {
		return nil, ErrMonthJobNil
	}
	if m.job.ID != "" {
		return m.scheduler.GetJobByID(m.job.ID)
	}
	if m.job.Ali != "" {
		return m.scheduler.GetJobByAlias(m.job.Ali)
	}
	if m.job.Name != "" {
		return m.scheduler.GetJobByName(m.job.Name)
	}
	return nil, ErrJobNotFound
}

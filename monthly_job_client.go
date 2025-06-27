package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type MonthJobClientInterface interface {
	// AtTime 设置运行时间
	AtTime(days []int, hour, minute, second int) MonthJobClientInterface
	Alias(alias string) MonthJobClientInterface
	JobID(id string) MonthJobClientInterface
	Name(name string) MonthJobClientInterface
	Tags(tags ...string) MonthJobClientInterface
	Task(task any, parameters ...any) MonthJobClientInterface
	Timeout(timeout time.Duration) MonthJobClientInterface
	Watch(watch func(event JobWatchInterface)) MonthJobClientInterface
	DefaultHooks() MonthJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) MonthJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) MonthJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface
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

func (m *MonthJobClient) AtTime(days []int, hour, minute, second int) MonthJobClientInterface {
	m.job.AtTime(days, hour, minute, second)
	return m
}

func (m *MonthJobClient) Alias(alias string) MonthJobClientInterface { // 设置别名
	m.job.Alias(alias)
	return m
}

func (m *MonthJobClient) JobID(id string) MonthJobClientInterface {
	m.job.JobID(id)
	return m
}

func (m *MonthJobClient) Name(name string) MonthJobClientInterface {
	m.job.Names(name)
	return m
}

func (m *MonthJobClient) Tags(tags ...string) MonthJobClientInterface {
	m.job.Tags(tags...)
	return m
}

func (m *MonthJobClient) Task(task any, parameters ...any) MonthJobClientInterface {
	m.job.Task(task, parameters...)
	return m
}

func (m *MonthJobClient) Timeout(timeout time.Duration) MonthJobClientInterface {
	m.job.Timeout(timeout)
	return m
}

func (m *MonthJobClient) Watch(watch func(event JobWatchInterface)) MonthJobClientInterface {
	m.job.Watch(watch)
	return m
}

func (m *MonthJobClient) DefaultHooks() MonthJobClientInterface {
	m.job.DefaultHooks()
	return m
}

func (m *MonthJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface {
	m.job.BeforeJobRuns(eventListenerFunc)
	return m
}

func (m *MonthJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) MonthJobClientInterface {
	m.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return m
}

func (m *MonthJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) MonthJobClientInterface {
	m.job.AfterJobRuns(eventListenerFunc)
	return m
}

func (m *MonthJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface {
	m.job.AfterJobRunsWithError(eventListenerFunc)
	return m
}

func (m *MonthJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) MonthJobClientInterface {
	m.job.AfterJobRunsWithPanic(eventListenerFunc)
	return m
}

func (m *MonthJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) MonthJobClientInterface {
	m.job.AfterLockError(eventListenerFunc)
	return m
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

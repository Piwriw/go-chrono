package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type CronJobClientInterface interface {
	// CronExpr Linux Cron 表达式必须填写
	CronExpr(expr string) *CronJob
	Alias(alias string) *CronJob
	JobID(id string) *CronJob
	Name(name string) *CronJob
	Task(task any, parameters ...any) *CronJob
	Timeout(timeout time.Duration) *CronJob
	Watch(watch func(event JobWatchInterface)) *CronJob
	DefaultHooks() *CronJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob
	Add() (gocron.Job, error)
	BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type CronJobClient struct {
	scheduler *Scheduler
	job       *CronJob
}

var _ CronJobClientInterface = (*CronJobClient)(nil)

func (c *CronJobClient) CronExpr(expr string) *CronJob {
	return c.job.CronExpr(expr)
}

func (c *CronJobClient) Alias(alias string) *CronJob {
	return c.job.Alias(alias)
}

func (c *CronJobClient) JobID(id string) *CronJob {
	return c.job.JobID(id)
}

func (c *CronJobClient) Name(name string) *CronJob {
	return c.job.Names(name)
}

func (c *CronJobClient) Task(task any, parameters ...any) *CronJob {
	return c.job.Task(task, parameters...)
}

func (c *CronJobClient) Timeout(timeout time.Duration) *CronJob {
	return c.job.Timeout(timeout)
}

func (c *CronJobClient) Watch(watch func(event JobWatchInterface)) *CronJob {
	return c.job.Watch(watch)
}

func (c *CronJobClient) DefaultHooks() *CronJob {
	return c.job.DefaultHooks()
}

func (c *CronJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.job.BeforeJobRuns(eventListenerFunc)
}

func (c *CronJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
	return c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.job.AfterJobRuns(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.job.AfterJobRunsWithError(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
	return c.job.AfterJobRunsWithPanic(eventListenerFunc)
}

func (c *CronJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.job.AfterLockError(eventListenerFunc)
}

// Add 添加任务
func (c *CronJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
	}
	return c.scheduler.AddCronJob(c.job)
}

func (c *CronJobClient) BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	return c.scheduler.AddCronJobs(cronJobs...)
}

// Remove 删除任务
func (c *CronJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job == nil {
		return ErrCronJobNil
	}

	// 1. 优先 JobID
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	// 2. 其次 Alias
	if c.job.Ali != "" {
		return c.scheduler.RemoveJobByAlias(c.job.Ali)
	}
	// 3. Name 或多条件才遍历
	if c.job.Name != "" {
		return c.scheduler.RemoveJobByName(c.job.Name)
	}
	return ErrJobNotFound
}

// Get 获取任务
func (c *CronJobClient) Get() (gocron.Job, error) {
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
